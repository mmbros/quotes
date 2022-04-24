package quotegetter

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// Result represents the info returned by the GetQuote function in case of success
type Result struct {
	URL      string
	Price    float32
	Currency string
	Date     time.Time
}

// QuoteGetter interface
type QuoteGetter interface {
	Source() string
	Client() *http.Client
	GetQuote(ctx context.Context, isin, url string) (*Result, error)
}

// NormalizeCurrency return the standard ISO4217 representation
// of the known currency
func NormalizeCurrency(currency string) string {
	if strings.EqualFold(currency, "euro") {
		return "EUR"
	}
	return currency
}
