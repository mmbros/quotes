package progress

import (
	"fmt"
	"io"
	"time"

	prog "github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Progress struct {
	pwriter  prog.Writer
	trackers map[string]*prog.Tracker
}

var unitsNA = prog.Units{
	Notation:         "",
	NotationPosition: prog.UnitsNotationPositionAfter,
	Formatter:        func(int64) string { return "n/a" },
}

var unitsCurrency = prog.Units{
	Notation:         "€ ",
	NotationPosition: prog.UnitsNotationPositionBefore,
	Formatter:        formatCurrency,
}

func formatCurrency(value int64) string {
	return fmt.Sprintf("%.2f", float32(value)/100)
}

func trackerUnits(currency string) prog.Units {
	u := unitsCurrency
	// var symbol string
	// switch currency {
	// case "USD":
	// 	symbol = "$"
	// case "GBP":
	// 	symbol = "£"
	// case "EUR":
	// 	symbol = "€"
	// default:
	// 	symbol = currency
	// }
	// u.Notation = symbol + " "

	u.Notation = currency + " "
	return u
}

func trackerValue(price float32) int64 {
	return int64(price * 100)
}

func New(w io.Writer, trackers int) *Progress {

	// instantiate a Progress Writer and set up the options
	pw := prog.NewWriter()

	pw.SetNumTrackersExpected(trackers)
	pw.SetOutputWriter(w)

	pw.SetAutoStop(false)
	pw.SetTrackerLength(7)
	pw.ShowETA(false)
	pw.ShowOverallTracker(true)
	pw.ShowTime(true)
	pw.ShowTracker(false)
	pw.ShowValue(true)
	pw.SetMessageWidth(40)

	pw.SetSortBy(prog.SortByMessage)
	pw.SetStyle(prog.StyleDefault)
	pw.SetTrackerPosition(prog.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.Style().Colors = prog.StyleColorsExample
	// pw.Style().Options.PercentFormat = "%4.1f%%"

	p := Progress{
		pwriter:  pw,
		trackers: map[string]*prog.Tracker{},
	}

	return &p
}

func (p *Progress) Render() {
	if p == nil {
		return
	}
	p.pwriter.Render()
}

func (p *Progress) Stop() {
	if p == nil {
		return
	}
	p.pwriter.Stop()
}

// InitTrackerIfNew adds a tracker, if not already exists.
func (p *Progress) InitTrackerIfNew(name string) {
	if p == nil {
		return
	}
	tracker, ok := p.trackers[name]
	if !ok {
		tracker = &prog.Tracker{
			Message: name,
			Total:   0,
			Units:   unitsNA,
		}
		p.pwriter.AppendTracker(tracker)
		p.trackers[name] = tracker
	}
}

func (p *Progress) SetSuccess(name string, message string, price float32, currency string) {
	if p == nil {
		return
	}
	tracker := p.trackers[name]

	tracker.Units = trackerUnits(currency)
	tracker.SetValue(trackerValue(price))
	tracker.UpdateMessage(name + " " + text.FgHiBlack.Sprint(message))
	tracker.MarkAsDone()
}

func (p *Progress) SetError(name string) {
	if p == nil {
		return
	}
	tracker := p.trackers[name]
	tracker.UpdateMessage(name + " " + text.FgRed.Sprint("error"))
	tracker.MarkAsErrored()
}
