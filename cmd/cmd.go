package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/mmbros/flagx"
)

const usageApp = `Usage:
    %s <command> [options]

Available Commands:
    get (g)        Get the quotes of the specified isins
    server (se)    Start an http server to show json files
    sources (so)   Show available sources
    tor (t)        Check if Tor network will be used
    version (v)    Version information

Common options:
    -h, --help   Help informations
`

func initApp() *flagx.Command {

	app := &flagx.Command{
		ParseExec: parseExecApp,

		SubCmd: map[string]*flagx.Command{
			"get,g": {
				ParseExec: parseExecGet,
			},
			"server,se": {
				ParseExec: parseExecServer,
			},
			"sources,so": {
				ParseExec: parseExecSources,
			},
			"tor,t": {
				ParseExec: parseExecTor,
			},
			"version,v": {
				ParseExec: parseExecVersion,
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
func Execute(stderr io.Writer) int {
	app := initApp()
	if err := flagx.Run(app); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}
