package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmbros/quotes/internal/quotegetter/scrapers/fondidocit"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers/fundsquarenet"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers/googlecrypto"
	"github.com/mmbros/quotes/internal/quotegetter/scrapers/morningstarit"
	"github.com/mmbros/quotes/internal/sources"
)

const usageSources = `Usage:
    %s

Prints list of available sources.
`

// MODULE VARIABLE
var availableQuoteGetters *sources.QuoteGetterSources

func init() {
	availableQuoteGetters = initSources()
}

func initSources() *sources.QuoteGetterSources {
	srcs := sources.NewQuoteGetterSources()
	srcs.Add("fondidocit", fondidocit.NewQuoteGetter)
	srcs.Add("morningstarit", morningstarit.NewQuoteGetter)
	srcs.Add("fundsquarenet", fundsquarenet.NewQuoteGetter)
	srcs.Add("googlecrypto-EUR", googlecrypto.NewQuoteGetterFactory("EUR"))
	return srcs
}

func parseExecSources(fullname string, arguments []string) error {

	// parse the arguments
	flags := NewFlags(fullname, fgAppSources)
	flags.SetUsage(usageSources, fullname)

	err := flags.Parse(arguments)

	// handle help
	if err == flag.ErrHelp {
		flags.Usage()
		return nil
	}
	if err != nil {
		return err
	}

	execSources(os.Stdout)
	return nil
}

func execSources(w io.Writer) {
	fmt.Fprintln(w, availableQuoteGetters)
}
