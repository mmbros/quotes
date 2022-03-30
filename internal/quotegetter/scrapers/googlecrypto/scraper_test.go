package googlecrypto

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers/testingscraper"

	"github.com/google/go-cmp/cmp"
	// "github.com/google/go-cmp/cmp/cmpopts"
)

func getTestScraper() scrapers.Scraper {
	return &scraper{"googlecrypto", nil, "EUR"}
}

func TestNewQuoteGetter(t *testing.T) {
	testingscraper.TestNewQuoteGetter(t, NewQuoteGetterFactory("EUR"))
}

func TestGetSearch(t *testing.T) {
	scr := getTestScraper()
	testingscraper.TestGetSearch(t, "", scr)
}

func TestParseSearch(t *testing.T) {

	tests := []struct {
		name string
		isin string
		url  string
		err  error
	}{
		{"ok",
			"BTC",
			"www.google.com/finance/quote/BTC-EUR",
			nil,
		},
	}

	scr := getTestScraper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix := fmt.Sprintf("ParseSearch[%s]", tt.name)

			url, err := scr.ParseSearch(nil, tt.isin)
			if testingscraper.CheckError(t, prefix, err, tt.err) {
				return
			}

			if url != tt.url {
				t.Errorf("%s: URL expected %q, found %q", prefix, tt.url, url)
			}
		})
	}

}

func TestGetInfo(t *testing.T) {
	scr := getTestScraper()
	testingscraper.TestGetInfo(t, "", scr)
}

func TestParseInfo(t *testing.T) {

	tests := []struct {
		name     string
		html     string
		expected *scrapers.ParseInfoResult
		err      error
	}{
		{
			name: "ok",
			html: `<div jscontroller="NdbN0c" jsaction="oFr1Ad:uxt3if;" jsname="AS5Pxb" data-mid="/g/11bvvzqspv" data-entity-type="3" data-is-crypto="true" data-source="BTC" data-target="EUR" data-last-price="40474.87415" data-last-normal-market-timestamp="1648292399" data-tz-offset="0"><div class="rPF6Lc" jsname="OYCkv"><div class="ln0Gqe"><div jsname="LXPcOd" class=""><div class="AHmHk"><span class=""><div jsname="ip75Cb" class="kf1m0"><div class="YMlKec fxKbKc">40.474,87</div></div></span></div></div><div jsname="CGyduf" class=""><div class="enJeMd"><span jsname="Fe7oBc" class="NydbP nZQ6l tnNmPe" data-disable-percent-toggle="true" data-multiplier-for-price-change="1" aria-label="Aumento: 0,27%"><div jsname="m6NnIb" class="zWwE1"><div class="JwB6zf" style="font-size: 16px;"><span class="V53LMb" aria-hidden="true"><svg focusable="false" width="16" height="16" viewBox="0 0 24 24" class=" NMm5M"><path d="M4 12l1.41 1.41L11 7.83V20h2V7.83l5.58 5.59L20 12l-8-8-8 8z"></path></svg></span>0,27%</div></div></span><span class="P2Luy Ez2Ioe ZYVHBb">+107,89 Oggi</span></div></div></div></div><div class="ygUjEc" jsname="Vebqub">26 mar, 10:59:59 UTC · <a href="https://www.google.com/intl/it_IT/googlefinance/disclaimer/"><span class="koPoYd">Disclaimer</span></a></div></div>`,
			expected: &scrapers.ParseInfoResult{
				IsinStr:     "BTC",
				CurrencyStr: "EUR",
				PriceStr:    "40474.87415",
				DateStr:     "1648292399",
				DateLayout:  "unix",
			},
			err: nil,
		},
		{
			name: "ko-price-date",
			html: `<div jscontroller="NdbN0c" jsaction="oFr1Ad:uxt3if;" jsname="AS5Pxb" data-mid="/g/11bvvzqspv" data-entity-type="3" data-is-crypto="true" data-source="BTC" data-target="EUR"   data-tz-offset="0"><div class="rPF6Lc" jsname="OYCkv"><div class="ln0Gqe"><div jsname="LXPcOd" class=""><div class="AHmHk"><span class=""><div jsname="ip75Cb" class="kf1m0"><div class="YMlKec fxKbKc">40.474,87</div></div></span></div></div><div jsname="CGyduf" class=""><div class="enJeMd"><span jsname="Fe7oBc" class="NydbP nZQ6l tnNmPe" data-disable-percent-toggle="true" data-multiplier-for-price-change="1" aria-label="Aumento: 0,27%"><div jsname="m6NnIb" class="zWwE1"><div class="JwB6zf" style="font-size: 16px;"><span class="V53LMb" aria-hidden="true"><svg focusable="false" width="16" height="16" viewBox="0 0 24 24" class=" NMm5M"><path d="M4 12l1.41 1.41L11 7.83V20h2V7.83l5.58 5.59L20 12l-8-8-8 8z"></path></svg></span>0,27%</div></div></span><span class="P2Luy Ez2Ioe ZYVHBb">+107,89 Oggi</span></div></div></div></div><div class="ygUjEc" jsname="Vebqub">26 mar, 10:59:59 UTC · <a href="https://www.google.com/intl/it_IT/googlefinance/disclaimer/"><span class="koPoYd">Disclaimer</span></a></div></div>`,
			err:  scrapers.ErrNoResultFound,
		},
		{
			name: "ko-price",
			html: `<div jscontroller="NdbN0c" jsaction="oFr1Ad:uxt3if;" jsname="AS5Pxb" data-mid="/g/11bvvzqspv" data-entity-type="3" data-is-crypto="true" data-source="BTC" data-target="EUR" data-last-normal-market-timestamp="1648292399" data-tz-offset="0"><div class="rPF6Lc" jsname="OYCkv"><div class="ln0Gqe"><div jsname="LXPcOd" class=""><div class="AHmHk"><span class=""><div jsname="ip75Cb" class="kf1m0"><div class="YMlKec fxKbKc">40.474,87</div></div></span></div></div><div jsname="CGyduf" class=""><div class="enJeMd"><span jsname="Fe7oBc" class="NydbP nZQ6l tnNmPe" data-disable-percent-toggle="true" data-multiplier-for-price-change="1" aria-label="Aumento: 0,27%"><div jsname="m6NnIb" class="zWwE1"><div class="JwB6zf" style="font-size: 16px;"><span class="V53LMb" aria-hidden="true"><svg focusable="false" width="16" height="16" viewBox="0 0 24 24" class=" NMm5M"><path d="M4 12l1.41 1.41L11 7.83V20h2V7.83l5.58 5.59L20 12l-8-8-8 8z"></path></svg></span>0,27%</div></div></span><span class="P2Luy Ez2Ioe ZYVHBb">+107,89 Oggi</span></div></div></div></div><div class="ygUjEc" jsname="Vebqub">26 mar, 10:59:59 UTC · <a href="https://www.google.com/intl/it_IT/googlefinance/disclaimer/"><span class="koPoYd">Disclaimer</span></a></div></div>`,
			expected: &scrapers.ParseInfoResult{
				IsinStr:     "BTC",
				CurrencyStr: "EUR",
				PriceStr:    "",
				DateStr:     "1648292399",
				DateLayout:  "unix",
			},
		},
		{
			name: "ko-date",
			html: `<div jscontroller="NdbN0c" jsaction="oFr1Ad:uxt3if;" jsname="AS5Pxb" data-mid="/g/11bvvzqspv" data-entity-type="3" data-is-crypto="true" data-source="BTC" data-target="EUR" data-last-price="40474.87415" data-tz-offset="0"><div class="rPF6Lc" jsname="OYCkv"><div class="ln0Gqe"><div jsname="LXPcOd" class=""><div class="AHmHk"><span class=""><div jsname="ip75Cb" class="kf1m0"><div class="YMlKec fxKbKc">40.474,87</div></div></span></div></div><div jsname="CGyduf" class=""><div class="enJeMd"><span jsname="Fe7oBc" class="NydbP nZQ6l tnNmPe" data-disable-percent-toggle="true" data-multiplier-for-price-change="1" aria-label="Aumento: 0,27%"><div jsname="m6NnIb" class="zWwE1"><div class="JwB6zf" style="font-size: 16px;"><span class="V53LMb" aria-hidden="true"><svg focusable="false" width="16" height="16" viewBox="0 0 24 24" class=" NMm5M"><path d="M4 12l1.41 1.41L11 7.83V20h2V7.83l5.58 5.59L20 12l-8-8-8 8z"></path></svg></span>0,27%</div></div></span><span class="P2Luy Ez2Ioe ZYVHBb">+107,89 Oggi</span></div></div></div></div><div class="ygUjEc" jsname="Vebqub">26 mar, 10:59:59 UTC · <a href="https://www.google.com/intl/it_IT/googlefinance/disclaimer/"><span class="koPoYd">Disclaimer</span></a></div></div>`,
			expected: &scrapers.ParseInfoResult{
				IsinStr:     "BTC",
				CurrencyStr: "EUR",
				PriceStr:    "40474.87415",
				DateStr:     "",
				DateLayout:  "unix",
			},
		},
	}

	scr := getTestScraper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if testingscraper.CheckError(t, "Goquery", err, nil) {
				return
			}

			result, err := scr.ParseInfo(doc, "")
			if testingscraper.CheckError(t, "ParseInfo", err, tt.err) {
				return
			}

			if diff := cmp.Diff(tt.expected, result, nil); diff != "" {
				t.Errorf("%s: mismatch (-want +got):\n%s", tt.name, diff)
			}
		})
	}

}
