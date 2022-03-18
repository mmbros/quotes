package cmd

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/mmbros/flagx"
)

type flagsGroup int

// Type of flags grou

const (
	fgApp flagsGroup = iota
	fgAppGet
	fgAppTor
	fgAppSources
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

type Flags struct {
	config     string
	configType string
	database   string
	dryrun     bool
	isins      []string
	proxy      string
	sources    []string
	workers    int
	mode       string

	flagSet  *flag.FlagSet
	fullname string
}

func NewFlags(fullname string, flagsgroup flagsGroup) *Flags {

	fs := flag.NewFlagSet(fullname, flag.ContinueOnError)

	flags := &Flags{}
	flags.flagSet = fs
	flags.fullname = fullname

	// use the same output as flag.CommandLine
	fs.SetOutput(flag.CommandLine.Output())

	// flagx.AliasedStringVar(fs, &args.Config, app.NamesConfig, "", "config file")
	// flagx.AliasedStringVar(fs, &args.ConfigType, app.NamesConfigType, "", "used if config file does not have the extension in the name; accepted values are: YAML, TOML and JSON")
	// flagx.AliasedBoolVar(fs, &args.Dryrun, app.NamesDryrun, false, "perform a trial run with no request/updates made")
	// flagx.AliasedStringVar(fs, &args.Proxy, app.NamesProxy, "", "default proxy")
	// flagx.AliasedIntVar(fs, &args.Workers, app.NamesWorkers, app.DefaultWorkers, "number of workers")
	// flagx.AliasedStringVar(fs, &args.Database, app.NamesDatabase, "", "sqlite3 database used to save the quotes")
	// flagx.AliasedStringVar(fs, &args.Mode, app.NamesMode, app.DefaultMode, `result mode: "1" first success or last error (default), "U" all errors until first success, "A" all`)
	// flagx.AliasedStringsVar(fs, &args.Isins, app.NamesIsins, "list of isins to get the quotes")
	// flagx.AliasedStringsVar(fs, &args.Sources, app.NamesSources, "list of sources to get the quotes from")

	// flags common to all operation
	flagx.AliasedStringVar(fs, &flags.config, namesConfig, "", "")
	flagx.AliasedStringVar(fs, &flags.configType, namesConfigType, "", "")

	// flags for Get or Tor operation
	if flagsgroup == fgAppGet || flagsgroup == fgAppTor {
		flagx.AliasedStringVar(fs, &flags.proxy, namesProxy, "", "")
	}

	// flags only for Get operation
	if flagsgroup == fgAppGet {

		flagx.AliasedBoolVar(fs, &flags.dryrun, namesDryrun, false, "")
		flagx.AliasedIntVar(fs, &flags.workers, namesWorkers, defaultWorkers, "")
		flagx.AliasedStringVar(fs, &flags.database, namesDatabase, "", "")
		flagx.AliasedStringVar(fs, &flags.mode, namesMode, defaultMode, "")
		flagx.AliasedStringsVar(fs, &flags.isins, namesIsins, "")
		flagx.AliasedStringsVar(fs, &flags.sources, namesSources, "")
	}
	// fs.Usage = func() {
	// 	fmt.Fprintf(fs.Output(), usage, fullname)
	// }

	return flags
}

// SetUsage set the usage function of the inner FlagSet
func (flags *Flags) SetUsage(format string, a ...interface{}) {
	fs := flags.flagSet
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), format, a...)
	}
}

// IsPassed checks if the flag was passed in the command-line arguments.
// names is a string that contains the comma separated aliases of the flag.
func (flags *Flags) IsPassed(names string) bool {
	return flagx.IsPassed(flags.flagSet, names)
}

// Appname returns the app name from the fullname of the command
//
// Example:
//   fullname = "app cmd sub-cmd sub-sub-cmd"
//   output   = "app"
func (flags *Flags) Appname() string {
	astr := strings.Split(flags.fullname, " ")
	if len(astr) == 0 {
		return ""
	}
	return astr[0]
}

func (flags *Flags) Parse(arguments []string) error { return flags.flagSet.Parse(arguments) }

func (flags *Flags) Usage() { flags.flagSet.Usage() }

func (flags *Flags) Output() io.Writer { return flags.flagSet.Output() }
