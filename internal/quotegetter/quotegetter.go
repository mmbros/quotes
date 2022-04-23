package quotegetter

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// QuoteGetter interface
type QuoteGetter interface {
	Source() string
	Client() *http.Client
	GetQuote(ctx context.Context, isin, url string) *Result
}

// Result represents the info returned by the GetQuote function
type Result struct {
	URL      string
	Price    float32
	Currency string
	Date     time.Time
	Err      error
}

func (res *Result) Error() error { return res.Err }

func (res *Result) String() string {
	if res.Err != nil {
		return "n/a"
	}
	return fmt.Sprintf("%.2f %s", res.Price, res.Currency)
}

// NormalizeCurrency return the standard ISO4217 representation
// of the known currency
func NormalizeCurrency(currency string) string {
	if strings.EqualFold(currency, "euro") {
		return "EUR"
	}
	return currency
}
