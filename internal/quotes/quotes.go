package quotes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mmbros/quotes/internal/progress"
	"github.com/mmbros/quotes/internal/quotegetter"
	"github.com/mmbros/taskengine"
)

// SourceIsins struct represents the isins to get from a specific source
type SourceIsins struct {
	Source  string   `json:"source,omitempty"`
	Workers int      `json:"workers,omitempty"`
	Proxy   string   `json:"proxy,omitempty"`
	Isins   []string `json:"isins,omitempty"`
}

// Result contains the result informations of the retrieved quote.
//
// Result.Date field is a pointer in order to omit zero dates.
// see https://stackoverflow.com/questions/32643815/json-omitempty-with-time-time-field
type Result struct {
	Isin      string               `json:"isin,omitempty"`
	Source    string               `json:"source,omitempty"`
	Instance  int                  `json:"instance"`
	URL       string               `json:"url,omitempty"`
	Price     float32              `json:"price,omitempty"`
	Currency  string               `json:"currency,omitempty"`
	Date      *time.Time           `json:"date,omitempty"` // need a pointer to omit zero date
	TimeStart time.Time            `json:"time_start"`
	TimeEnd   time.Time            `json:"time_end"`
	ErrMsg    string               `json:"error,omitempty"`
	Err       error                `json:"-"`
	Status    taskengine.EventType `json:"status"`
}

// workerTask struct contains the info for retrieve the quote by a source.
// It implements the taskengine.Task interface
type workerTask struct {
	isin string
	url  string
}

// workerResult type is returned by the work functions of the
// taskengine.Worker objects.
// It collects both result and error returned by the
// quotegetter.GetQuote methods to implements the taskengine.Result interface.
type workerResult struct {
	*quotegetter.Result
	Err error
}

// TaskID method of the taskengine.Task interface
func (t *workerTask) TaskID() taskengine.TaskID {
	return taskengine.TaskID(t.isin)
}

// String representation of the task.
// Method of the taskengine.Result interface
func (r *workerResult) String() string {
	if r.Err != nil {
		return "n/a"
	}
	return fmt.Sprintf("%.2f %s", r.Price, r.Currency)
}

// The error returned by the Work function.
// Method of the taskengine.Result interface
func (r *workerResult) Error() error {
	return r.Err
}

// Get retrieves the quotes specified by the SourceIsins object.
// The mode parameters specified the taskengine mode of execution.
func Get(availableSources quotegetter.Sources, items []*SourceIsins, mode taskengine.Mode, wProgress io.Writer) ([]*Result, error) {

	// saveResult return true if the event is a result that have to be saved
	// according to the taskengine.Mode argument.
	var saveResult func(*taskengine.Event) bool
	switch mode {
	case taskengine.FirstSuccessOrLastResult:
		saveResult = func(e *taskengine.Event) bool { return e.IsFirstSuccessOrLastResult() }
	case taskengine.ResultsUntilFirstSuccess:
		saveResult = func(e *taskengine.Event) bool { return e.IsResultUntilFirstSuccess() }
	case taskengine.SuccessOrErrorResults:
		saveResult = func(e *taskengine.Event) bool { return e.IsSuccessOrError() }
	default:
		saveResult = func(e *taskengine.Event) bool { return e.IsResult() }
	}

	// Init the chan that will receive the events (containing the results).
	eventc, err := getEventsChan(availableSources, items)
	if err != nil {
		return nil, err
	}

	// Show progress if wProgress writer is defined.
	// Note that no error is raised if progr is nil.
	var progr *progress.Progress
	if wProgress != nil {
		progr = progress.New(wProgress, len(items))
	}
	go progr.Render()

	results := []*Result{}
	for event := range eventc {

		taskid := string(event.Task.TaskID())
		etype := event.Type()

		progr.InitTrackerIfNew(taskid)

		// handle progress update
		if event.IsFirstSuccessOrLastResult() {
			if etype == taskengine.EventSuccess {
				wres := event.Result.(*workerResult)
				progr.SetSuccess(taskid, string(event.WorkerID), wres.Price, wres.Currency)
			} else {
				progr.SetError(taskid)
			}
		}

		// handle results update
		if saveResult(event) {

			result := &Result{
				Isin:      taskid,
				Source:    string(event.WorkerID),
				Instance:  event.WorkerInst,
				TimeStart: event.TimeStart,
				TimeEnd:   event.TimeEnd,
				Status:    etype,
			}

			if etype == taskengine.EventSuccess {
				wres := event.Result.(*workerResult)
				result.Price = wres.Price
				result.Currency = wres.Currency
				result.URL = wres.URL
				result.Date = &wres.Date
			} else {
				result.Err = event.Result.Error()
				if etype == taskengine.EventCanceled {
					result.ErrMsg = context.Canceled.Error()
				} else {
					result.ErrMsg = result.Err.Error()
				}
			}

			results = append(results, result)
		}
	} // end event loop

	if progr != nil {
		progr.Render()
		time.Sleep(time.Millisecond * 200)
		progr.Stop()
	}

	return results, nil
}

func initQuoteGetters(availableSources quotegetter.Sources, src []*SourceIsins) (map[string]quotegetter.QuoteGetter, error) {
	quoteGetter := make(map[string]quotegetter.QuoteGetter)

	proxyClient := map[string]*http.Client{}

	for _, s := range src {
		name := s.Source

		client, ok := proxyClient[s.Proxy]
		if !ok {
			client, err := quotegetter.DefaultClient(s.Proxy)
			if err != nil {
				return nil, err
			}
			proxyClient[s.Proxy] = client
		}

		fn := availableSources[name]
		if fn == nil {
			panic("invalid source: " + name)
		}
		quoteGetter[name] = fn(name, client)
	}

	return quoteGetter, nil
}

// func (scnr *Scenario) ExecuteEvents() (chan *taskengine.Event, error) {

// 	if scnr.ws == nil {
// 		return nil, errors.New("must run RandomWorkersAndTasks before")
// 	}

// 	ctx := context.Background()
// 	eng, err := taskengine.NewEngine(scnr.ws, scnr.wts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return eng.ExecuteEvents(ctx)
// }

// getEventsChan ...
func getEventsChan(availableSources quotegetter.Sources, items []*SourceIsins) (chan *taskengine.Event, error) {

	// check input
	if err := checkListOfSourceIsins(availableSources, items); err != nil {
		return nil, err
	}

	// Workers
	ws := make([]*taskengine.Worker, 0, len(items))

	// WorkerTasks
	wts := make(taskengine.WorkerTasks)

	quoteGetter, err := initQuoteGetters(availableSources, items)
	if err != nil {
		return nil, err
	}

	for _, item := range items {

		qg := quoteGetter[item.Source]

		// work function of the source
		wfn := func(ctx context.Context, worker *taskengine.Worker, inst int, task taskengine.Task) taskengine.Result {
			//  from taskengine.Task to taskGetQuote
			t := task.(*workerTask)
			r, err := qg.GetQuote(ctx, t.isin, t.url)
			return &workerResult{r, err}
		}

		// worker
		w := &taskengine.Worker{
			WorkerID:  taskengine.WorkerID(item.Source),
			Instances: item.Workers,
			Work:      wfn,
		}
		ws = append(ws, w)

		// Tasks
		ts := make(taskengine.Tasks, 0, len(item.Isins))
		for _, isin := range item.Isins {
			ts = append(ts, &workerTask{
				isin: isin,
				url:  "",
			})
		}
		wts[w.WorkerID] = ts

	}

	eng, err := taskengine.NewEngine(ws, wts)
	if err != nil {
		return nil, err
	}
	return eng.ExecuteEvents(context.Background())
}

// checkListOfSourceIsins checks the validity of the given SourceIsins items
func checkListOfSourceIsins(availableSources quotegetter.Sources, items []*SourceIsins) error {
	used := map[string]struct{}{}

	for _, item := range items {

		if _, ok := used[item.Source]; ok {
			return fmt.Errorf("duplicate source %q", item.Source)
		}
		used[item.Source] = struct{}{}

		if !availableSources.Exists(item.Source) {
			return fmt.Errorf("source %q not available", item.Source)
		}
		if item.Workers <= 0 {
			return fmt.Errorf("source %q with invalid workers %d", item.Source, item.Workers)
		}
	}
	return nil
}
