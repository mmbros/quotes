package cmd

import (
	"flag"
	"strings"

	"github.com/mmbros/flagx"
)

// Names of the command line arguments (flagx names)
const (
	namesConfig     = "config,c"
	namesConfigType = "config-type"
	namesDatabase   = "database,d"
	namesDryrun     = "dry-run,n"
	namesIsins      = "isins,i"
	namesMode       = "mode,m"
	namesProxy      = "proxy,p"
	namesSources    = "sources,s"
	namesWorkers    = "workers,w"
)

// Default args value
const (
	defaultMode    = "1"
	defaultWorkers = 1
)

type Args struct {
	config     string
	configType string
	database   string
	dryrun     bool
	isins      []string
	proxy      string
	sources    []string
	workers    int
	mode       string

	flagSet *flag.FlagSet
}

// IsPassed checks if the flag was passed in the command-line arguments.
// names is a string that contains the comma separated aliases of the flag.
func (args *Args) IsPassed(names string) bool {
	if args == nil || args.flagSet == nil {
		panic("appargs.FlagSet = nil")
	}
	return flagx.IsPassed(args.flagSet, names)
}

func getAppname(fullname string) string {
	astr := strings.Split(fullname, " ")
	if len(astr) == 0 {
		return ""
	}
	return astr[0]
}

func NewArgs(fullname string) *Args {

	fs := flag.NewFlagSet(fullname, flag.ContinueOnError)

	args := &Args{}
	args.flagSet = fs

	// use the same output as flag.CommandLine
	fs.SetOutput(flag.CommandLine.Output())

	flagx.AliasedStringVar(fs, &args.config, namesConfig, "", "")
	flagx.AliasedStringVar(fs, &args.configType, namesConfigType, "", "")
	flagx.AliasedBoolVar(fs, &args.dryrun, namesDryrun, false, "")
	flagx.AliasedStringVar(fs, &args.proxy, namesProxy, "", "")
	flagx.AliasedIntVar(fs, &args.workers, namesWorkers, defaultWorkers, "")
	flagx.AliasedStringVar(fs, &args.database, namesDatabase, "", "")
	flagx.AliasedStringVar(fs, &args.mode, namesMode, defaultMode, "")
	flagx.AliasedStringsVar(fs, &args.isins, namesIsins, "")
	flagx.AliasedStringsVar(fs, &args.sources, namesSources, "")

	// fs.Usage = func() {
	// 	fmt.Fprintf(fs.Output(), usage, fullname)
	// }

	return args
}
