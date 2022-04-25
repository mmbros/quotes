package quotes

import (
	"fmt"
	"io"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mmbros/taskengine"
)

type Stat struct {
	NumSuccess  int
	NumError    int
	NumCanceled int
}

type Stats struct {
	TimeStart time.Time
	TimeEnd   time.Time
	Task      map[string]*Stat
	Worker    map[string]*Stat
}

func (s *Stat) Update(status taskengine.EventType) {
	switch status {
	case taskengine.EventSuccess:
		s.NumSuccess++
	case taskengine.EventError:
		s.NumError++
	case taskengine.EventCanceled:
		s.NumCanceled++
	}
}

func NewStats(results []*Result) *Stats {

	if results == nil {
		return nil
	}

	stats := &Stats{
		Task:   map[string]*Stat{},
		Worker: map[string]*Stat{},
	}

	for iter, result := range results {
		// update times
		if (iter == 0) || stats.TimeStart.After(result.TimeStart) {
			stats.TimeStart = result.TimeStart
		}
		if stats.TimeEnd.Before(result.TimeEnd) {
			stats.TimeEnd = result.TimeEnd
		}

		// update task
		taskstat := stats.Task[result.Isin]
		if taskstat == nil {
			taskstat = &Stat{}
		}
		taskstat.Update(result.Status)
		stats.Task[result.Isin] = taskstat

		// update task
		workerstat := stats.Worker[result.Source]
		if workerstat == nil {
			workerstat = &Stat{}
		}
		workerstat.Update(result.Status)
		stats.Worker[result.Source] = workerstat
	}

	return stats
}

func statsSummary(astat map[string]*Stat) *Stat {
	var summary Stat

	for _, s := range astat {
		if s.NumSuccess > 0 {
			summary.NumSuccess++
		} else if s.NumError > 0 {
			summary.NumError++
		} else if s.NumCanceled > 0 {
			summary.NumCanceled++
		}
	}
	return &summary
}

func (stats *Stats) TaskSummary() *Stat {
	return statsSummary(stats.Task)
}

func (stats *Stats) WorkerSummary() *Stat {
	return statsSummary(stats.Worker)
}

func (stats *Stats) Elapsed() time.Duration {
	return stats.TimeEnd.Sub(stats.TimeStart)
}

func (stats *Stats) Tasks() int {
	return len(stats.Task)
}
func (stats *Stats) Workers() int {
	return len(stats.Worker)
}

// func (stats *Stats) Fprintln(w io.Writer) {
// 	var sum *Stat
// 	fmt.Fprintf(w, "Elapsed: %v\n", stats.Elapsed())

// 	sum = stats.TaskSummary()
// 	fmt.Fprintf(w, "%d quotes (success:%d, error:%d, canceled:%d)\n", stats.Tasks(), sum.NumSuccess, sum.NumError, sum.NumCanceled)

// 	sum = stats.WorkerSummary()
// 	fmt.Fprintf(w, "%d sources (success:%d, error:%d, canceled:%d)\n", stats.Workers(), sum.NumSuccess, sum.NumError, sum.NumCanceled)
// }

func (stats *Stats) Fprintln(w io.Writer) {
	sum := stats.TaskSummary()
	fmt.Fprintf(w, "%d quotes (%s", stats.Tasks(), text.FgGreen.Sprint(sum.NumSuccess, " success"))
	if sum.NumError > 0 {
		fmt.Fprintf(w, ", %s", text.FgRed.Sprint(sum.NumError, " error"))
	}
	if sum.NumCanceled > 0 {
		fmt.Fprintf(w, ", %s", text.FgBlack.Sprint(sum.NumCanceled, " canceled"))
	}

	fmt.Fprintf(w, ") from %d sources in %v\n", stats.Workers(), stats.Elapsed().Round(time.Millisecond))
}
