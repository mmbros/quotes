package quotes

import (
	"context"
	"errors"

	"github.com/mmbros/quotes/internal/quotegetterdb"
)

func (r *Result) insert(db *quotegetterdb.QuoteDatabase) error {
	var qr *quotegetterdb.QuoteRecord

	// assert := func(b bool, label string) {
	// 	if !b {
	// 		panic("failed assert: " + label)
	// 	}
	// }

	// assert(r != nil, "r != nil")
	// assert(db != nil, "db != nil")

	// skip context.Canceled errors
	// if r.Err != nil  {
	// 	if err, ok := r.Err.(*scrapers.Error); ok {
	// 		if !errors.Is(err, context.Canceled) {
	// 			return nil
	// 		}
	// 	}
	// }

	if errors.Is(r.Err, context.Canceled) {
		return nil
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

func dbInsert(dbpath string, results []*Result) error {
	if len(dbpath) == 0 {
		return nil
	}

	// save to database
	db, err := quotegetterdb.Open(dbpath)
	if db != nil {
		defer db.Close()

		for _, r := range results {
			err = r.insert(db)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
