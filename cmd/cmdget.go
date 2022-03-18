package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/mmbros/quote/internal/quote"
)

func parseExecGet(fullname string, arguments []string) error {
	var cfg *Config

	// parse the arguments
	flags := NewFlags(fullname, fgAppGet)
	flags.SetUsage(usageGet, fullname)

	err := flags.Parse(arguments)

	// handle help
	if err == flag.ErrHelp {
		flags.Usage()
		return nil
	}
	if err != nil {
		return err
	}

	// get configuration
	cfg, err = getConfig(flags, quote.Sources())
	if err != nil {
		return err
	}

	return execGet(flags, cfg)
}

// func parseGet(fullname string, arguments []string) (*Args, error) {
// 	args := NewArgs(fullname, fgAppGet)
// 	args.Usage(usageGet, fullname)

// 	err := fs.Parse(arguments)

// 	return args, err
// }

func execGet(flags *Flags, cfg *Config) error {

	if flags.dryrun {
		return printDryRunInfo(flags.Output(), flags, cfg)
	}

	// do retrieves the quotes
	sis := cfg.SourceIsinsList()
	return quote.Get(sis, cfg.Database, cfg.taskengMode)
}

func printDryRunInfo(w io.Writer, flags *Flags, cfg *Config) error {
	fmt.Fprintln(w, "Dry Run")
	if flags.IsPassed(namesConfig) {
		fmt.Fprintf(w, "Using configuration file %q\n", flags.config)
	}
	if cfg.Database != "" {
		fmt.Fprintf(w, "Database: %q\n", cfg.Database)
	}
	fmt.Fprintf(w, "Mode: %q (%d)\n", cfg.Mode, cfg.taskengMode)
	sis := cfg.SourceIsinsList()
	fmt.Fprint(w, "Tasks:", jsonString(sis))
	return nil
}
