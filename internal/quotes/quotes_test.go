package quotes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/mmbros/quotes/internal/quotegetter"
	"github.com/mmbros/taskengine"
	"github.com/stretchr/testify/assert"
)

// String returns a json string representation of the object.
func jsonString(obj interface{}) string {
	// print config
	json, _ := json.MarshalIndent(obj, "", "  ")
	return string(json)
}

type dummyQuoteGetter struct {
	source string
	client *http.Client
}

func newDummyQuoteGetter(source string, client *http.Client) quotegetter.QuoteGetter {
	return &dummyQuoteGetter{source, client}
}

func (qg *dummyQuoteGetter) Source() string       { return qg.source }
func (qg *dummyQuoteGetter) Client() *http.Client { return qg.client }

func (qg *dummyQuoteGetter) GetQuote(ctx context.Context, isin, url string) (*quotegetter.Result, error) {

	cases := map[string]*struct {
		err  bool
		wait int
	}{
		"source1-isin1": {
			err:  false,
			wait: 10,
		},
		"source2-isin1": {
			err:  true,
			wait: 20,
		},
	}
	key := qg.source + "-" + isin
	c := cases[key]

	if c == nil {
		return nil, fmt.Errorf("dummyQuoteGetter: %s, %s: not implemented", qg.source, isin)
	}

	if c.wait > 0 {
		time.Sleep(time.Duration(c.wait) * time.Millisecond)
	}
	if c.err {
		return nil, fmt.Errorf("dummyQuoteGetter: %s, %s: generic error", qg.source, isin)
	}
	res := &quotegetter.Result{
		Date:     time.Now(),
		Currency: "EUR",
		Price:    12.35,
	}
	return res, nil
}

func TestCheckListOfSourceIsins(t *testing.T) {
	availableSources := quotegetter.Sources{
		"source1": newDummyQuoteGetter,
		"source2": newDummyQuoteGetter,
	}

	cases := []struct {
		input  []*SourceIsins
		errmsg string
	}{
		{
			input: []*SourceIsins{
				{
					Source:  "source1",
					Workers: 1,
					Isins:   []string{"isin1"},
				},
				{
					Source:  "source2",
					Workers: 2,
					Isins:   []string{"isin1", "isin2"},
				},
			},
			errmsg: "",
		},
		{
			input: []*SourceIsins{
				{
					Source:  "sourceY",
					Workers: 1,
					Isins:   []string{"isin1"},
				},
			},
			errmsg: "not available",
		},
		{
			input: []*SourceIsins{
				{
					Source:  "source1",
					Workers: -1,
					Isins:   []string{"isin1"},
				},
			},
			errmsg: "invalid workers",
		},
		{
			input: []*SourceIsins{
				{
					Source:  "source1",
					Workers: 1,
					Isins:   []string{"isin1"},
				},
				{
					Source:  "source2",
					Workers: 2,
					Isins:   []string{"isin1", "isin2"},
				},
				{
					Source:  "source1",
					Workers: 1,
					Isins:   []string{"isin2"},
				},
			},
			errmsg: "duplicate source",
		},
	}

	for _, c := range cases {
		err := checkListOfSourceIsins(availableSources, c.input)
		if c.errmsg == "" {
			assert.NoError(t, err)
		} else {
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), c.errmsg)
			}
		}
	}
}

func TestGetResults(t *testing.T) {
	availableSources := quotegetter.Sources{
		"source1": newDummyQuoteGetter,
		"source2": newDummyQuoteGetter,
	}

	sis := []*SourceIsins{
		{
			Source:  "source1",
			Workers: 1,
			Isins:   []string{"isin1"},
		},
		{
			Source:  "source2",
			Workers: 2,
			Isins:   []string{"isin1", "isin2"},
		},
	}
	res, err := Get(availableSources, sis, taskengine.AllResults, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(res))
		// t.Fatalf("res %v", jsonString(res))
	}

}
