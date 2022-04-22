package quote

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mmbros/quotes/internal/quotegetter"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers"
	"github.com/mmbros/quotes/internal/quotegetterdb"
	"github.com/mmbros/taskengine"
)

// SourceIsins struct represents the isins to get from a specific source
type SourceIsins struct {
	Source  string   `json:"source,omitempty"`
	Workers int      `json:"workers,omitempty"`
	Proxy   string   `json:"proxy,omitempty"`
	Isins   []string `json:"isins,omitempty"`
}

// taskGetQuote struct contains the info for retrieve the quote by a source.
// It implements the taskengine.Task interface
type taskGetQuote struct {
	isin string
	url  string
}

// TaskID method of the taskengine.Task interface
func (t *taskGetQuote) TaskID() taskengine.TaskID {
	return taskengine.TaskID(t.isin)
}

// resultGetQuote contains the result informations of the retrieved quote.
// It implements the taskengine.Result interface.
//
// resultGetQuote.Date field is a pointer in order to omit zero dates.
// see https://stackoverflow.com/questions/32643815/json-omitempty-with-time-time-field
type resultGetQuote struct {
	Isin      string     `json:"isin,omitempty"`
	Source    string     `json:"source,omitempty"`
	Instance  int        `json:"instance"`
	URL       string     `json:"url,omitempty"`
	Price     float32    `json:"price,omitempty"`
	Currency  string     `json:"currency,omitempty"`
	Date      *time.Time `json:"date,omitempty"` // need a pointer to omit zero date
	TimeStart time.Time  `json:"time_start"`
	TimeEnd   time.Time  `json:"time_end"`
	ErrMsg    string     `json:"error,omitempty"`
	Err       error      `json:"-"`
}

// String representation of the task.
// Method of the taskengine.Result interface
func (r *resultGetQuote) String() string {
	if r.Err != nil {
		return "n/a"
	}
	return fmt.Sprintf("%.2f %s", r.Price, r.Currency)
}

// The error returned by the Work function.
// Method of the taskengine.Result interface
func (r *resultGetQuote) Error() error {
	return r.Err
}

func (r *resultGetQuote) dbInsert(db *quotegetterdb.QuoteDatabase) error {
	var qr *quotegetterdb.QuoteRecord

	// assert := func(b bool, label string) {
	// 	if !b {
	// 		panic("failed assert: " + label)
	// 	}
	// }

	// assert(r != nil, "r != nil")
	// assert(db != nil, "db != nil")

	// skip context.Canceled errors
	if r.Err != nil {
		if err, ok := r.Err.(*scrapers.Error); ok {
			if !errors.Is(err, context.Canceled) {
				return nil
			}
		}
	}
	qr = &quotegetterdb.QuoteRecord{
		Isin:     r.Isin,
		Source:   r.Source,
		Price:    r.Price,
		Currency: r.Currency,
		URL:      r.URL,
		ErrMsg:   r.ErrMsg,
	}
	if r.Date != nil {
		qr.Date = *r.Date
	}
	// isin and source are mandatory
	// assert(len(qr.Isin) > 0, "len(qr.Isin) > 0")
	// assert(len(qr.Source) > 0, "len(qr.Source) > 0")

	// save to database
	return db.InsertQuotes(qr)
}

func dbInsert(dbpath string, results []*resultGetQuote) error {
	if len(dbpath) == 0 {
		return nil
	}

	// save to database
	db, err := quotegetterdb.Open(dbpath)
	if db != nil {
		defer db.Close()

		for _, r := range results {
			err = r.dbInsert(db)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// checkListOfSourceIsins checks the validity of the given SourceIsins items
func checkListOfSourceIsins(items []*SourceIsins) error {
	used := map[string]struct{}{}

	for _, item := range items {

		if _, ok := used[item.Source]; ok {
			return fmt.Errorf("duplicate source %q", item.Source)
		}
		used[item.Source] = struct{}{}

		if _, ok := availableSources[item.Source]; !ok {
			return fmt.Errorf("source %q not available", item.Source)
		}
		if item.Workers <= 0 {
			return fmt.Errorf("source %q with invalid workers %d", item.Source, item.Workers)
		}
	}
	return nil
}

// Get retrieves the quotes specified by the SourceIsins object.
// The mode parameters specified the taskengine mode of execution.
// The results quotes are printed in json format.
// The quotes are also saved to the database, if the dbpath is given.
func Get(items []*SourceIsins, dbpath string, mode taskengine.Mode) error {

	results, err := getResults(items, mode)
	if err != nil {
		return err
	}

	// save to database, if not empty
	err = dbInsert(dbpath, results)
	if err != nil {
		fmt.Println(err)
	}

	json, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return err
	}

	fmt.Println(string(json))

	return nil
}

// getResults executes the tasks in order to retrieve the quotes.
func getResults(items []*SourceIsins, mode taskengine.Mode) ([]*resultGetQuote, error) {

	// check input
	if err := checkListOfSourceIsins(items); err != nil {
		return nil, err
	}

	// Workers
	ws := make([]*taskengine.Worker, 0, len(items))

	// WorkerTasks
	wts := make(taskengine.WorkerTasks)

	quoteGetter, err := initQuoteGetters(items)
	if err != nil {
		return nil, err
	}

	for _, item := range items {

		qg := quoteGetter[item.Source]

		// work function of the source
		wfn := func(ctx context.Context, worker *taskengine.Worker, inst int, task taskengine.Task) taskengine.Result {
			//  from taskengine.Task to taskGetQuote
			t := task.(*taskGetQuote)

			time1 := time.Now()
			res, err := qg.GetQuote(ctx, t.isin, t.url)
			time2 := time.Now()

			r := &resultGetQuote{
				Instance:  inst,
				TimeStart: time1,
				TimeEnd:   time2,
				Err:       err,
			}
			if res != nil {
				r.Isin = res.Isin
				r.Source = res.Source
				r.Price = res.Price
				r.Currency = res.Currency
				r.URL = res.URL
				if !res.Date.IsZero() {
					r.Date = &res.Date
				}
			}
			if err != nil {
				r.ErrMsg = err.Error()
				if e, ok := err.(quotegetter.Error); ok {
					r.Isin = e.Isin()
					r.Source = e.Source()
					r.URL = e.URL()
				}
			}
			return r
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
			ts = append(ts, &taskGetQuote{
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
	resChan, err := eng.Execute(context.Background(), mode)
	if err != nil {
		return nil, err
	}

	results := []*resultGetQuote{}
	for r := range resChan {
		res := r.(*resultGetQuote)
		results = append(results, res)
	}

	return results, nil
}
