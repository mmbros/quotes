package googlecrypto

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmbros/quotes/internal/quotegetter"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers"
)

// scraper gets stock/fund prices from www.google.com/finance/quote/
type scraper struct {
	name     string
	client   *http.Client
	currency string
}

// // NewQuoteGetter creates a new QuoteGetter
// // that gets stock/fund prices from www.google.com/finance/quote/
// func NewQuoteGetter(name string, client *http.Client, currency string) quotegetter.QuoteGetter {
// 	return scrapers.NewQuoteGetter(&scraper{name, client, currency})
// }

// NewQuoteGetter creates a new QuoteGetter
// that gets stock/fund prices from www.google.com/finance/quote/
// currency must be in (EUR, USD)
func NewQuoteGetterFactory(currency string) quotegetter.NewQuoteGetterFunc {
	return func(name string, client *http.Client) quotegetter.QuoteGetter {
		return scrapers.NewQuoteGetter(&scraper{name, client, currency})
		// return NewQuoteGetter(name, client, currency)
	}
}

// Name returns the name of the scraper
func (s *scraper) Source() string {
	return s.name
}

// Client returns the http.Client of the scraper
func (s *scraper) Client() *http.Client {
	return s.client
}

// GetSearch creates the http.Request to get the search page for the specified `isin`.
// It returns the http.Response or nil if the scraper can build the url of the info page
// directly from the `isin`.
// The response document will be parsed by ParseSearch to extract the info url.
func (s *scraper) GetSearch(ctx context.Context, isin string) (*http.Request, error) {
	return nil, nil
}

// ParseSearch parse the html of the search page to find the URL of the info page.
// `doc` can be nil if the url of the info page can be build directly from the `isin`.
// It returns the url of the info page.
func (s *scraper) ParseSearch(doc *goquery.Document, isin string) (string, error) {
	url := fmt.Sprintf("www.google.com/finance/quote/%s-%s", isin, s.currency)
	return url, nil
}

// GetInfo executes the http GET of the `url` of info page for the specified `isin`.
// `url` and `isin` must be defined.
// The response document will be parsed by ParseInfo to extract the info url.
func (s *scraper) GetInfo(ctx context.Context, isin, url string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ParseInfo is ...
func (s *scraper) ParseInfo(doc *goquery.Document, isin string) (*scrapers.ParseInfoResult, error) {
	/*
		<div jscontroller="NdbN0c" jsaction="oFr1Ad:uxt3if;" jsname="AS5Pxb" data-mid="/g/11bvvzqspv"
		     data-entity-type="3" data-is-crypto="true"
			 data-source="BTC" data-target="EUR" data-last-price="40401.127700000005"
			 data-last-normal-market-timestamp="1648247340" data-tz-offset="0">
			...
			<div class="ygUjEc" jsname="Vebqub">Mar 25, 10:31:00 PM UTC ·
				<a href="https://www.google.com/intl/en-US_IT/googlefinance/disclaimer/">
					<span class="koPoYd">Disclaimer</span>
				</a>
			</div>
		</div>


		data-last-normal-market-timestamp="1648247340"
		data-tz-offset="0"

		1648247340  -> Mar 25, 10:31:00 PM UTC

	*/

	r := new(scrapers.ParseInfoResult)
	r.DateLayout = scrapers.LayoutUnixTimestamp

	div := doc.Find("div[data-is-crypto]")
	r.DateStr, _ = div.Attr("data-last-normal-market-timestamp")
	r.IsinStr, _ = div.Attr("data-source")
	r.CurrencyStr, _ = div.Attr("data-target")
	r.PriceStr, _ = div.Attr("data-last-price")

	if r.DateStr == "" && r.PriceStr == "" {
		return r, scrapers.ErrNoResultFound
	}
	return r, nil
}
