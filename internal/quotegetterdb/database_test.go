package quotegetterdb

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/mmbros/quotes/internal/quotes"
)

const (
	isin1 = "ISIN00001111"
	isin2 = "ISIN99992222"

	source1 = "quotesource1.com"
	source2 = "quotesource2.it"
)

func testURL(source, isin string) string {
	return fmt.Sprintf("http://127.0.0.1/%s/%s/", source, isin)
}

var loc = time.Local

var records = []*QuoteRecord{
	{
		Isin:      isin1,
		Source:    source1,
		Price:     10.1,
		Currency:  "USD",
		Date:      time.Date(2020, 01, 01, 0, 0, 0, 0, loc),
		Timestamp: time.Date(2020, 01, 01, 10, 11, 0, 0, loc),
		URL:       testURL(source1, isin1),
	},
	{
		Isin:      isin1,
		Source:    source1,
		Timestamp: time.Date(2020, 01, 02, 10, 22, 0, 0, loc),
		ErrMsg:    "Isin not found",
	},
	{
		Isin:      isin1,
		Source:    source1,
		Price:     10.3,
		Currency:  "USD",
		Date:      time.Date(2020, 01, 03, 0, 0, 0, 0, loc),
		Timestamp: time.Date(2020, 01, 03, 10, 33, 0, 0, loc),
		URL:       testURL(source1, isin1),
	},
	{
		Isin:      isin1,
		Source:    source2,
		Price:     10.22,
		Currency:  "EUR",
		Date:      time.Date(2020, 02, 01, 0, 0, 0, 0, loc),
		Timestamp: time.Date(2020, 02, 01, 0, 0, 0, 0, loc),
		URL:       testURL(source2, isin1),
	},
	{
		Isin:      isin1,
		Source:    source2,
		Timestamp: time.Date(2020, 02, 02, 0, 0, 0, 0, loc),
		ErrMsg:    "Isin not found",
	},
	{
		Isin:      isin1,
		Source:    source2,
		Timestamp: time.Date(2020, 02, 03, 0, 0, 0, 0, loc),
		ErrMsg:    "Isin not found",
	},
	{
		Isin:      isin1,
		Source:    source2,
		Timestamp: time.Date(2020, 02, 04, 0, 0, 0, 0, loc),
		ErrMsg:    "Isin not found",
	},
	{
		Isin:      isin2,
		Source:    source1,
		Timestamp: time.Date(2020, 02, 03, 0, 0, 0, 0, loc),
		ErrMsg:    "Isin not found",
	}}

const dbpath = ":memory:"

//const dbpath = "/tmp/quote.sqlite3"

func mustOpenDB() *QuoteDatabase {

	qdb, err := Open(dbpath)
	if err != nil {
		panic(err)
	}
	return qdb
}

func TestOpen(t *testing.T) {

	// create the database
	qdb, err := Open(dbpath)
	if err != nil {
		t.Errorf("open database: unexpected error: %v", err)
	}
	qdb.Close()

	qdb, err = Open("xyz://error")
	if err == nil {
		t.Errorf("open database: expecting error, found no error")
		qdb.Close()
	}

}

func TestInsertQuotes(t *testing.T) {
	qdb := mustOpenDB()
	defer qdb.Close()

	err := qdb.InsertQuotesRecords(records...)
	if err != nil {
		t.Fatal(err)
	}
}

/*
func TestSelectLastQuotes(t *testing.T) {
	qdb := mustOpenDB()
	defer qdb.Close()

	err := qdb.InsertQuotes(records...)
	if err != nil {
		t.Fatal(err)
	}

	res, err := qdb.SelectLastQuotes()
	if err != nil {
		t.Fatal(err)
	}
	for j, r := range res {

		t.Logf("[%d] %v\n", j, r)
	}
	t.Fail()
}
*/

/*
func TestExtractPath(t *testing.T) {
	testCases := []struct {
		dns      string
		expected string
	}{
		{"/tmp/finanze/db.sqlite3", "/tmp/finanze"},
		{"file:/tmp/finanze/db.sqlite3", "/tmp/finanze"},
		{"file:/tmp/finanze/db.sqlite3?cache=shared", "/tmp/finanze"},
		{":memory:", ""},
		{"file::memory:", ""},
		{"file::memory:?cache=shared", ""},
		{"file:/tmp/finanze/db.sqlite3?cache=shared&mode=memory", ""},
		{"file:/tmp/finanze/db.sqlite3?mode=memory&cache=shared", ""},
	}

	for _, tc := range testCases {
		res := extractDir(tc.dns)
		if res != tc.expected {
			t.Errorf("extractPath(%q): got %q, expected %q", tc.dns, res, tc.expected)
		}
	}
}
*/

func TestDBInsert(t *testing.T) {

	time1 := time.Date(2020, 01, 01, 0, 0, 0, 0, loc)
	time2 := time.Date(2020, 02, 02, 0, 0, 0, 0, loc)

	res1 := &quotes.Result{
		Isin:     isin1,
		Source:   source1,
		Price:    10.1,
		Currency: "USD",
		Date:     &time1,
		URL:      testURL(source1, isin1),
	}

	res2 := &quotes.Result{
		Isin:     isin2,
		Source:   source2,
		Price:    20.2,
		Currency: "EUR",
		Date:     &time2,
		URL:      testURL(source2, isin2),
		Err:      errors.New("isin not found"),
	}

	err := DBInsert(dbpath, []*quotes.Result{res1, res2})
	if err != nil {
		t.Error(err)
	}
}
