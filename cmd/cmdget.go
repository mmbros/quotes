package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/mmbros/quote/internal/quote"
)

const usageGet = `Usage: 
    %[1]s [options]

Options:
    -c, --config      path     config file
        --config-type string   used if config file does not have the extension in the name;
                               accepted values are: YAML, TOML and JSON 
    -i, --isins       strings  list of isins to get the quotes
    -n, --dry-run              perform a trial run with no request/updates made
    -p, --proxy       url      default proxy
    -s, --sources     strings  list of sources to get the quotes from
    -w, --workers     int      number of workers (default %[2]d)
    -d, --database    dns      sqlite3 database used to save the quotes
    -m, --mode        char     result mode (default %[3]q): 
                                  "1" first success or last error
                                  "U" all errors until first success 
                                  "A" all 
`

func parseExecGet(fullname string, arguments []string) error {
	var cfg *Config

	// parse the arguments
	flags := NewFlags(fullname, fgAppGet)
	flags.SetUsage(usageGet, fullname, defaultWorkers, defaultMode)

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

	fmt.Fprintf(w, "%s: Dry Run\n", flags.fullname)

	// prints config file info
	cfg.cfi.Fprintln(w)

	if cfg.Database != "" {
		fmt.Fprintf(w, "Database: %q\n", cfg.Database)
	}
	fmt.Fprintf(w, "Mode: %q\n", cfg.Mode)
	sis := cfg.SourceIsinsList()
	fmt.Fprint(w, "Tasks: ", jsonString(sis))
	return nil
}
