package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mmbros/flagx"
	"github.com/mmbros/quote/internal/quote"
)

const (
	usageApp = `Usage:
    %s <command> [options]

Available Commands:
    get      Get the quotes of the specified isins
    sources  Show available sources
    tor      Checks if Tor network will be used
`
	usageGet = `Usage:
    %s [options]

Options:
    -c, --config      path     config file (default is $HOME/.quote.yaml)
        --config-type string   used if config file does not have the extension in the name;
                               accepted values are: YAML, TOML and JSON 
    -i, --isins       strings  list of isins to get the quotes
    -n, --dry-run              perform a trial run with no request/updates made
    -p, --proxy       url      default proxy
    -s, --sources     strings  list of sources to get the quotes from
    -w, --workers     int      number of workers (default 1)
    -d, --database    dns      sqlite3 database used to save the quotes
    -m, --mode        char     result mode: "1" first success or last error (default)
                                            "U" all errors until first success 
                                            "A" all 
`

	usageTor = `Usage:
     %s [options]

Checks if Tor network will be used to get the quote.

To use the Tor network the proxy must be defined through:
	1. proxy argument parameter
	2. proxy config file parameter
	3. HTTP_PROXY, HTTPS_PROXY and NOPROXY enviroment variables.

Options:
    -c, --config      path    config file (default is $HOME/.quote.yaml)
	    --config-type string  used if config file does not have the extension in the name;
	                          accepted values are: YAML, TOML and JSON 
    -p, --proxy       url     proxy to test the Tor network
`

	usageSources = `Usage:
	%s

Prints list of available sources.
`
)

func initApp() *flagx.Command {

	app := &flagx.Command{
		ParseExec: parseExecApp,

		SubCmd: map[string]*flagx.Command{
			"get,g": {
				ParseExec: parseExecGet,
			},
			"tor,t": {
				ParseExec: parseExecTor,
			},
			"sources,s": {
				ParseExec: runQuotesSources,
			},
		},
	}

	return app
}

func parseExecApp(fullname string, arguments []string) error {

	fs := flag.NewFlagSet(fullname, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageApp, fullname)
	}

	err := fs.Parse(arguments)
	if err == nil || err == flag.ErrHelp {
		fs.Usage()
		err = nil
	}

	return err
}

func runQuotesSources(fullname string, arguments []string) error {
	// it is used a module level declaration for test porpouses.
	// normally do: argsQuotes := &appArgs{}
	// argsQuotes := &appArgs{}

	fs := flag.NewFlagSet(fullname, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageSources, fullname)
	}

	err := fs.Parse(arguments)
	if err != nil {
		return err
	}

	sources := quote.Sources()
	fmt.Printf("Available sources: \"%s\"\n", strings.Join(sources, "\", \""))
	return nil
}

// Execute is the main function
func Execute() {
	app := initApp()

	err := flagx.Run(app)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
