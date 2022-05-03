package fundsquarenet

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-cmp/cmp"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers/testingscraper"
)

func getTestScraper() scrapers.Scraper {
	return &scraper{"fundsquarenet", nil}
}

func TestNewQuoteGetter(t *testing.T) {
	testingscraper.TestNewQuoteGetter(t, NewQuoteGetter)
}

func TestGetSearch(t *testing.T) {
	scr := getTestScraper()
	testingscraper.TestGetSearch(t, "", scr)
}

func TestParseSearch(t *testing.T) {
	const isin = "<ISIN>"
	scr := getTestScraper()
	url, err := scr.ParseSearch(nil, isin)
	if err != nil {
		t.Errorf("ParseSearch: %v", err)
	}
	if !strings.Contains(url, isin) {
		t.Errorf("ParseSearch: returned url %q does not contain isin %q", url, isin)
	}
}

func TestGetInfo(t *testing.T) {
	scr := getTestScraper()
	req, _ := testingscraper.TestGetInfo(t, "", scr)
	if req != nil {
		referer := req.Header.Get("Referer")
		isin := testingscraper.TestIsin
		if !strings.Contains(referer, isin) {
			t.Errorf("GetInfo: referer %q not contains isin %q", referer, isin)
		}
	}
}

func TestParseInfo(t *testing.T) {

	tests := []struct {
		name    string
		html    string
		want    *scrapers.ParseInfoResult
		wantErr error
	}{
		{
			name: "ok",
			html: `<div id="content"><table style="width: 100%"><tr><td><span style="font-weight: bold;">IE00B4TG9K96</span>&nbsp;&nbsp;PIMCO GIS Diversified Income Fund E Hgd EUR Dis&nbsp;&nbsp;</td><td></td></tr></table><table width="85%"><tr><td width="30%">Last NAV</td><td width="15%">11/09/2020</td><td width="55%"><span class="surligneorange">11.49&nbsp;EUR</span>&nbsp;<span style="color:#000000;text-align:left;padding:4px 0;"> 0.00 &nbsp;%&nbsp;<img src="/images/share/variationNulle.gif" style="vertical-align:middle;"/></span></td></tr></table>`,
			want: &scrapers.ParseInfoResult{
				IsinStr:     "IE00B4TG9K96",
				PriceStr:    "11.49",
				CurrencyStr: "EUR",
				DateStr:     "11/09/2020",
				DateLayout:  "02/01/2006",
			},
			wantErr: nil,
		},
		{
			name:    "unavailable",
			html:    `<html><head></head><body><div id="content"><table style="width: 100%"><tbody><tr><td><span style="font-weight: bold;">IE00B4TG9K96</span>  PIMCO GIS Diversified Income Fund E Hgd EUR Dis  </td><td></td></tr></tbody></table>`,
			want:    &scrapers.ParseInfoResult{},
			wantErr: scrapers.ErrPriceAndCurrencyString,
		},
		{
			name:    "not-found",
			html:    `<table width="100%"><tbody><tr><td valign="middle" align="center"><img src="/images/share/research.gif" alt="info"/></td><td valign="middle" align="center"><span class="title"><div class="contenu"><p class="zero"></p><p><span class="surligneorange"> No result</span> produced by your request.</p><p>Please<span class="surligneorange"> modify your search criteria</span>.</p><input type="submit" id="valida" name="valider" onclick="window.location=&#39;/search?fastSearch=O&amp;isISIN=O&amp;search=IT0005247157&#39;" class="back_search_w btn_w" value="Back to search"/><br class="clear_r"/><p></p></div><p class="ps">Number of results : <span>0</span></p></span></td></tr></tbody></table>`,
			want:    &scrapers.ParseInfoResult{},
			wantErr: scrapers.ErrNoResultFound,
		},
	}

	scr := getTestScraper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if testingscraper.CheckError(t, "goquery", err, nil) {
				return
			}

			res, err := scr.ParseInfo(doc, "")
			if testingscraper.CheckError(t, "ParseInfo", err, tt.wantErr) {
				return
			}

			if diff := cmp.Diff(tt.want, res, nil); diff != "" {
				t.Errorf("%s: mismatch (-want +got):\n%s", tt.name, diff)
			}
		})
	}
}
