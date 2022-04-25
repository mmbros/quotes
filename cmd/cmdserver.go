package cmd

import (
	"flag"
	"fmt"

	"github.com/mmbros/flagx"
	"github.com/mmbros/quotes/internal/server"
)

const usageServer = `Usage:
    %[1]s [options]

Start an http server to view a page with graphs based upon
the json files created with the get command.

Options:
    -a, --address    string   server address and port (default %[4]q)
    -f, --folder     path     folder containing the json files (default %[2]q)
    -r, --recursive  bool     search recursively all the json files of the 
                              sub-folders (default %[3]v)
Examples:
    # view the json files of the ./demo folder 
    %[1]s --folder ./demo
`

// Names of the command line arguments (flagx names)
const (
	namesFolder    = "folder,f"
	namesRecursive = "recursive,r"
	namesAddress   = "address,a"
)

// Default args value
const (
	defaultFolder    = "."
	defaultAddress   = ":6789"
	defaultRecursive = false
)

func parseExecServer(fullname string, arguments []string) error {
	var folder string
	var address string
	var recursive bool

	fs := flag.NewFlagSet(fullname, flag.ContinueOnError)

	// use the same output as flag.CommandLine
	fs.SetOutput(flag.CommandLine.Output())

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), usageServer, fullname, defaultFolder, defaultRecursive, defaultAddress)
	}

	flagx.AliasedStringVar(fs, &folder, namesFolder, defaultFolder, "")
	flagx.AliasedStringVar(fs, &address, namesAddress, defaultAddress, "")
	flagx.AliasedBoolVar(fs, &recursive, namesRecursive, defaultRecursive, "")

	// parse the arguments
	err := fs.Parse(arguments)

	// handle help
	if err == flag.ErrHelp {
		fs.Usage()
		return nil
	}
	return execServer(address, folder, recursive)
}

func execServer(serverAddressPort, jsonDataFolder string, recursive bool) error {
	return server.Run(serverAddressPort, jsonDataFolder, recursive)
}
