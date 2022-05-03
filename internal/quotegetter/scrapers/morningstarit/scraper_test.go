package morningstarit

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers/testingscraper"
)

func getTestScraper() scrapers.Scraper {
	return &scraper{"morningstarit", nil}
}

func TestSource(t *testing.T) {
	const name = "dummy"
	scr := &scraper{name, nil}
	if nameFound := scr.Source(); nameFound != name {
		t.Errorf("Source: found %q, expected %q", nameFound, name)
	}
}

func TestGetSearch(t *testing.T) {
	scr := getTestScraper()
	testingscraper.TestGetSearch(t, "", scr)
}

func TestParseSearch(t *testing.T) {
	tests := []struct {
		name    string
		isin    string
		html    string
		wantURL string
		wantErr error
	}{
		{
			name: "ok",
			isin: "IE00B4TG9K96",
			html: `<table id="ctl00_MainContent_fundTable" cellspacing="0" cellpadding="0" border="0" style="border-collapse:collapse;">
			<tr class="searchGridHeader">
				<th>Nome</th><th>ISIN</th>
			</tr><tr class="gridItem">
				<td class="msDataText searchLink"><a href="/it/funds/snapshot/snapshot.aspx?id=F000005GUM">PIMCO GIS Divers Inc E EURH Inc</a></td><td class="msDataText searchIsin"><span>IE00B4TG9K96</span></td>
			</tr>
		</table>`,
			wantURL: "/it/funds/snapshot/snapshot.aspx?id=F000005GUM",
			wantErr: nil,
		},
		{
			name:    "ko",
			isin:    "IE00B4TG9KAA",
			html:    ``,
			wantURL: "",
			wantErr: scrapers.ErrNoResultFound,
		},
		// {"IE00B4TG9KAA", "morningstar.it/search|IE00B4TG9KAA|ko.html", "", scrapers.ErrNoResultFound},
	}

	scr := getTestScraper()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if testingscraper.CheckError(t, "goquery", err, nil) {
				return
			}
			url, err := scr.ParseSearch(doc, tt.isin)
			if url != tt.wantURL {
				t.Errorf("ParseSearch: URL found %q, expected %q", url, tt.wantURL)
			}
			if err != tt.wantErr {
				t.Errorf("ParseSearch: ERR found %q, expected %q", err, tt.wantErr)
			}
		})
	}
}

func TestGetInfo(t *testing.T) {
	scr := getTestScraper()
	testingscraper.TestGetInfo(t, "", scr)
}

func TestParseInfo(t *testing.T) {

	testCases := []struct {
		title    string
		html     string
		priceStr string
		dateStr  string
		isinStr  string
	}{
		{
			title: "ok",
			html: `<div id="overviewQuickstatsDiv" xmlns:funs="funs" xmlns:rtxo="urn:RTExtensionObj" xmlns:erto="urn:ETRObj">
			<table class="snapshotTextColor snapshotTextFontStyle snapshotTable overviewKeyStatsTable" border="0"><tbody>
			<tr><td class="titleBarHeading" colspan="3">Sintesi</td></tr>
			<tr><td class="line heading">NAV<span class="heading"><br>28/08/2020</span></td><td class="line">&nbsp;</td><td class="line text">EUR&nbsp;126,370</td></tr>
			<tr><td class="line heading">Var.Ultima Quotazione</td><td class="line">&nbsp;</td><td class="line text">0,24%</td></tr>
			<tr><td class="line heading">Categoria Morningstar™</td><td class="line">&nbsp;</td><td class="line value text"><a href="https://www.morningstar.it/it/fundquickrank/default.aspx?category=EUCA000640" style="width:100%!important;">Azionari Italia</a></td></tr>
			<tr><td class="line heading">Categoria Assogestioni</td><td class="line">&nbsp;</td><td class="line text">Azionari Italia</td></tr>
			<tr><td class="line heading">Isin</td><td class="line">&nbsp;</td><td class="line text">IT0005247157</td></tr>
			<tr><td class="line heading">Fund Size (Mil)<span class="heading"><br>26/06/2020</span></td><td class="line">&nbsp;</td><td class="line text">EUR&nbsp;17,32</td></tr>
			<tr><td class="line heading">Share Class Size (Mil)<span class="heading"><br>26/06/2020</span></td><td class="line">&nbsp;</td><td class="line text">EUR&nbsp;5,06</td></tr>
			<tr><td class="line heading">Entrata (max)</td><td class="line">&nbsp;</td><td class="line text">-</td></tr>
			<tr><td class="line heading"><a href="http://www.morningstar.it/it/glossary/121049/spese-correnti-ongoing-charge.aspx" style="text-decoration:underline;">Spese correnti</a><span class="heading"><br>26/02/2019</span></td><td class="line">&nbsp; </td><td class="line text">1,51%</td></tr>
			</tbody></table>
			</div>`,
			priceStr: "126,370",
			dateStr:  "28/08/2020",
			isinStr:  "IT0005247157",
		},
	}

	scr := getTestScraper()

	for _, tc := range testCases {
		prefix := fmt.Sprintf("ParseSearch[%s]", tc.title)

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(tc.html))
		if testingscraper.CheckError(t, prefix+": Goquery", err, nil) {
			continue
		}

		res, err := scr.ParseInfo(doc, "")
		if err != nil {
			if tc.priceStr != "" {
				t.Errorf("[%s] Unexpected error %q", tc.title, err)
			}
			if err != scrapers.ErrNoResultFound {
				t.Errorf("[%s] Unexpected error %q, expected %q", tc.title, err, scrapers.ErrNoResultFound)
			}

			continue
		}

		t.Logf("[%s] -> %q", tc.title, res)

		if res.PriceStr != tc.priceStr {
			t.Errorf("[%s] PriceStr: expected %q, found %q", tc.title, tc.priceStr, res.PriceStr)
		}
		if res.DateStr != tc.dateStr {
			t.Errorf("[%s] DateStr: expected %q, found %q", tc.title, tc.dateStr, res.DateStr)
		}
		if res.IsinStr != tc.isinStr {
			t.Errorf("[%s] IsinStr: expected %q, found %q", tc.title, tc.isinStr, res.IsinStr)
		}
	}
}
