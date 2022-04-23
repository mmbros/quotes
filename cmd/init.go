package cmd

import (
	"github.com/mmbros/quotes/internal/quotegetter"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers/fondidocit"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers/fundsquarenet"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers/googlecrypto"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers/morningstarit"
)

// MODULE VARIABLE
var mAvailableSources quotegetter.Sources

func init() {
	mAvailableSources = initSources()
}

func initSources() quotegetter.Sources {
	return quotegetter.Sources{
		"fondidocit":       fondidocit.NewQuoteGetter,
		"morningstarit":    morningstarit.NewQuoteGetter,
		"fundsquarenet":    fundsquarenet.NewQuoteGetter,
		"googlecrypto-EUR": googlecrypto.NewQuoteGetterFactory("EUR"),
	}
}
