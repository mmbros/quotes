package quote

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmbros/quotes/internal/quotegetter"
)

// URL to check for TOR network
const torCheckURL = "https://check.torproject.org"

// TorProxyFromEnvironment returns the proxy from environment
// used to get the url used to check the Tor network.
func TorProxyFromEnvironment() string {
	if r, _ := http.NewRequest("GET", torCheckURL, nil); r != nil {
		if u, _ := http.ProxyFromEnvironment(r); u != nil {
			return u.String()
		}
	}
	return ""
}

// TorCheck checks if a Tor connection is used,
// retrieving the "https://check.torproject.org" page.
// It returns:
//  - bool:   true if Tor is used, false otherwise
//  - string: the message contained in the html page
//  - error:  if the message cannot be determined
func TorCheck(proxy string) (bool, string, error) {

	client, err := quotegetter.DefaultClient(proxy)
	if err != nil {
		return false, "", err
	}
	// Make request
	resp, err := client.Get(torCheckURL)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false, "", err
	}

	msg := strings.TrimSpace(doc.Find("h1").Text())
	if msg == "" {
		return false, "", fmt.Errorf("can't determine if you are using Tor")
	}

	return msg == "Congratulations. This browser is configured to use Tor.", msg, nil
}
