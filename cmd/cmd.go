package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/mmbros/flagx"
)

const usageApp = `Usage:
    %s <command> [options]

Available Commands:
    get (g)      Get the quotes of the specified isins
    sources (s)  Show available sources
    tor (t)      Checks if Tor network will be used
`

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
				ParseExec: parseExecSources,
			},
		},
	}

	return app
}

func parseExecApp(fullname string, arguments []string) error {
	// parse the arguments
	flags := NewFlags(fullname, fgApp)
	flags.SetUsage(usageApp, fullname)

	err := flags.Parse(arguments)

	// handle help
	if err == nil || err == flag.ErrHelp {
		flags.Usage()
		return nil
	}

	return err
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
