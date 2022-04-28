package quotes

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mmbros/taskengine"
)

func TestNewStats(t *testing.T) {

	T0 := time.Now().Unix()
	T := func(seconds int64) time.Time { return time.Unix(T0+seconds, 0) }

	tests := []struct {
		name    string
		results []*Result
		stats   *Stats
	}{
		{
			name:    "nil",
			results: nil,
			stats:   nil,
		},
		{
			name:    "zero",
			results: []*Result{},
			stats: &Stats{
				Task:   map[string]*Stat{},
				Worker: map[string]*Stat{},
			},
		},
		{
			name: "one",
			results: []*Result{
				{
					TimeStart: T(0),
					TimeEnd:   T(1),
					Status:    taskengine.EventError,
					Isin:      "isin1",
					Source:    "source1",
				},
			},
			stats: &Stats{
				TimeStart: T(0),
				TimeEnd:   T(1),
				Task: map[string]*Stat{
					"isin1": {0, 1, 0},
				},
				Worker: map[string]*Stat{
					"source1": {0, 1, 0},
				},
			},
		},
		{
			name: "multi",
			results: []*Result{
				{
					TimeStart: T(0),
					TimeEnd:   T(20),
					Status:    taskengine.EventError,
					Isin:      "isin1",
					Source:    "source1",
				},
				{
					TimeStart: T(10),
					TimeEnd:   T(30),
					Status:    taskengine.EventSuccess,
					Isin:      "isin1",
					Source:    "source2",
				},
				{
					TimeStart: T(15),
					TimeEnd:   T(30),
					Status:    taskengine.EventCanceled,
					Isin:      "isin1",
					Source:    "source3",
				},
				{
					TimeStart: T(0),
					TimeEnd:   T(9),
					Status:    taskengine.EventSuccess,
					Isin:      "isin2",
					Source:    "source2",
				},
			},

			stats: &Stats{
				TimeStart: T(0),
				TimeEnd:   T(30),
				Task: map[string]*Stat{
					"isin1": {1, 1, 1},
					"isin2": {1, 0, 0},
				},
				Worker: map[string]*Stat{
					"source1": {0, 1, 0},
					"source2": {2, 0, 0},
					"source3": {0, 0, 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := NewStats(tt.results)
			if diff := cmp.Diff(tt.stats, stats, nil); diff != "" {
				t.Errorf("%s: mismatch (-want +got):\n%s", tt.name, diff)
			}
		})
	}
}

func TestNewStats_Info(t *testing.T) {

	T0 := time.Now().Unix()
	T := func(seconds int64) time.Time { return time.Unix(T0+seconds, 0) }

	tests := []struct {
		name          string
		results       []*Result
		tasks         int
		workers       int
		elapsed       time.Duration
		taskSummary   *Stat
		workerSummary *Stat
	}{
		{
			name:          "zero",
			results:       []*Result{},
			taskSummary:   &Stat{},
			workerSummary: &Stat{},
		},
		{
			name: "one",
			results: []*Result{
				{
					TimeStart: T(0),
					TimeEnd:   T(1),
					Status:    taskengine.EventError,
					Isin:      "isin1",
					Source:    "source1",
				},
			},
			tasks:         1,
			workers:       1,
			elapsed:       time.Duration(1 * time.Second),
			taskSummary:   &Stat{0, 1, 0},
			workerSummary: &Stat{0, 1, 0},
		},
		{
			name: "multi",
			results: []*Result{
				{
					TimeStart: T(0),
					TimeEnd:   T(20),
					Status:    taskengine.EventError,
					Isin:      "isin1",
					Source:    "source1",
				},
				{
					TimeStart: T(10),
					TimeEnd:   T(30),
					Status:    taskengine.EventSuccess,
					Isin:      "isin1",
					Source:    "source2",
				},
				{
					TimeStart: T(15),
					TimeEnd:   T(30),
					Status:    taskengine.EventCanceled,
					Isin:      "isin1",
					Source:    "source3",
				},
				{
					TimeStart: T(0),
					TimeEnd:   T(9),
					Status:    taskengine.EventSuccess,
					Isin:      "isin2",
					Source:    "source2",
				},
			},
			tasks:         2,
			workers:       3,
			elapsed:       time.Duration(30 * time.Second),
			taskSummary:   &Stat{2, 0, 0},
			workerSummary: &Stat{1, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := NewStats(tt.results)

			if got := stats.Workers(); got != tt.workers {
				t.Errorf("%s: workers: want %d, got %d", tt.name, tt.workers, got)
			}
			if got := stats.Tasks(); got != tt.tasks {
				t.Errorf("%s: tasks: want %d, got %d", tt.name, tt.tasks, got)
			}
			if got := stats.Elapsed(); got != tt.elapsed {
				t.Errorf("%s: tasks: want %v, got %v", tt.name, tt.elapsed, got)
			}

			if diff := cmp.Diff(tt.taskSummary, stats.TaskSummary(), nil); diff != "" {
				t.Errorf("%s: mismatch (-want +got):\n%s", tt.name, diff)
			}
			if diff := cmp.Diff(tt.workerSummary, stats.WorkerSummary(), nil); diff != "" {
				t.Errorf("%s: mismatch (-want +got):\n%s", tt.name, diff)
			}
		})
	}
}

func TestNewStats_Fprintln(t *testing.T) {

	T0 := time.Now().Unix()
	T := func(seconds int64) time.Time { return time.Unix(T0+seconds, 0) }

	stats := &Stats{
		TimeStart: T(0),
		TimeEnd:   T(30),
		Task: map[string]*Stat{
			"isin1": {0, 1, 1},
			"isin2": {1, 0, 0},
		},
		Worker: map[string]*Stat{
			"source1": {},
			"source2": {},
			"source3": {},
		},
	}

	sb := &strings.Builder{}

	re := regexp.MustCompile(`2 quotes (.*1 success.*, .*1 error.*) from 3 sources in 30s`)
	want := "2 quotes (1 success, 1 error) from 3 sources in 30s"

	stats.Fprintln(sb)
	got := sb.String()
	if !re.MatchString(got) {
		if diff := cmp.Diff(want, got, nil); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}

}
