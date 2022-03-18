package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/mmbros/flagx"
	"github.com/mmbros/quote/internal/quote"
)

func parseExecGet(fullname string, arguments []string) error {
	var cfg *Config

	// parse the arguments
	args, err := parseGet(fullname, arguments)

	// handle help
	if err == flag.ErrHelp {
		args.flagSet.Usage()
		return nil
	}
	if err != nil {
		return err
	}

	appname := getAppname(fullname)

	// get configuration
	cfg, err = getConfig(appname, args, quote.Sources())
	if err != nil {
		return err
	}

	return execGet(args, cfg)
}

func parseGet(fullname string, arguments []string) (*Args, error) {
	args := NewArgs(fullname)

	fs := args.flagSet

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), usageGet, fullname)
	}
	// flagx.AliasedStringVar(fs, &args.Config, app.NamesConfig, "", "config file")
	// flagx.AliasedStringVar(fs, &args.ConfigType, app.NamesConfigType, "", "used if config file does not have the extension in the name; accepted values are: YAML, TOML and JSON")
	// flagx.AliasedBoolVar(fs, &args.Dryrun, app.NamesDryrun, false, "perform a trial run with no request/updates made")
	// flagx.AliasedStringVar(fs, &args.Proxy, app.NamesProxy, "", "default proxy")
	// flagx.AliasedIntVar(fs, &args.Workers, app.NamesWorkers, app.DefaultWorkers, "number of workers")
	// flagx.AliasedStringVar(fs, &args.Database, app.NamesDatabase, "", "sqlite3 database used to save the quotes")
	// flagx.AliasedStringVar(fs, &args.Mode, app.NamesMode, app.DefaultMode, `result mode: "1" first success or last error (default), "U" all errors until first success, "A" all`)
	// flagx.AliasedStringsVar(fs, &args.Isins, app.NamesIsins, "list of isins to get the quotes")
	// flagx.AliasedStringsVar(fs, &args.Sources, app.NamesSources, "list of sources to get the quotes from")

	flagx.AliasedStringVar(fs, &args.config, namesConfig, "", "")
	flagx.AliasedStringVar(fs, &args.configType, namesConfigType, "", "")
	flagx.AliasedBoolVar(fs, &args.dryrun, namesDryrun, false, "")
	flagx.AliasedStringVar(fs, &args.proxy, namesProxy, "", "")
	flagx.AliasedIntVar(fs, &args.workers, namesWorkers, defaultWorkers, "")
	flagx.AliasedStringVar(fs, &args.database, namesDatabase, "", "")
	flagx.AliasedStringVar(fs, &args.mode, namesMode, defaultMode, "")
	flagx.AliasedStringsVar(fs, &args.isins, namesIsins, "")
	flagx.AliasedStringsVar(fs, &args.sources, namesSources, "")

	err := fs.Parse(arguments)

	return args, err
}

func execGet(args *Args, cfg *Config) error {

	if args.dryrun {
		return printDryRunInfo(args.flagSet.Output(), args, cfg)
	}

	// do retrieves the quotes
	sis := cfg.SourceIsinsList()
	return quote.Get(sis, cfg.Database, cfg.taskengMode)
}

func printDryRunInfo(w io.Writer, args *Args, cfg *Config) error {
	fmt.Fprintln(w, "Dry Run")
	if args.IsPassed(namesConfig) {
		fmt.Fprintf(w, "Using configuration file %q\n", args.config)
	}
	if cfg.Database != "" {
		fmt.Fprintf(w, "Database: %q\n", cfg.Database)
	}
	fmt.Fprintf(w, "Mode: %q (%d)\n", cfg.Mode, cfg.taskengMode)
	sis := cfg.SourceIsinsList()
	fmt.Fprint(w, "Tasks:", jsonString(sis))
	return nil
}
