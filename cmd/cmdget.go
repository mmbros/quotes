package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmbros/quotes/internal/quotes"
)

const usageGet = `Usage: 
    %[1]s [options]

Options:
    -c, --config      path     config file
        --config-type string   used if config file does not have the extension in the name;
                               accepted values are: YAML, TOML and JSON 
    -d, --database    dns      sqlite3 database used to save the quotes
	-f, --force       bool     overwrite already existing output file
    -i, --isins       strings  list of isins to get the quotes
    -m, --mode        char     result mode (default %[3]q): 
                                  "1" first success or last error
                                  "U" all errors until first success 
                                  "A" all 
    -n, --dry-run              perform a trial run with no request/updates made
    -o, --output      path     pathname of the output file (default stdout)
    -p, --proxy       url      default proxy
    -s, --sources     strings  list of sources to get the quotes from
    -w, --workers     int      number of workers (default %[2]d)
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
	cfg, err = getConfig(flags, mAvailableSources.Names())
	if err != nil {
		return err
	}

	return execGet(flags, cfg)
}

func execGet(flags *Flags, cfg *Config) error {

	if flags.dryrun {
		return printDryRunInfo(flags.Output(), flags, cfg)
	}

	// handle the output
	wInfo := os.Stdout
	wOutput := os.Stdout // default
	if flags.output != "" {
		// overwrite existing file only if --force is specified
		flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		if !flags.force {
			flag |= os.O_EXCL // file must not exists
		}

		// create the file
		fout, err := os.OpenFile(flags.output, flag, 0666)
		if err != nil {
			return err
		}
		defer fout.Close()
		wOutput = fout
	}

	// do retrieves the quotes
	sis := cfg.SourceIsinsList()
	results, err := quotes.Get(mAvailableSources, sis, cfg.taskengMode, wInfo)
	if err != nil {
		return err
	}

	// fmt.Fprintf(wInfo, "\n%d task completed (%d success, %d error) in %v\n",
	// 	stats.TaskCompleted(),
	// 	stats.TaskSuccess,
	// 	stats.TaskError,
	// 	timeEnd.Sub(timeStart))
	// fmt.Fprintf(wInfo, "elapsed: %v\n", timeEnd.Sub(timeStart))

	// fmt.Fprintf(wInfo, "\n%d task completed (%d success, %d error) in %v\n",

	stats := quotes.NewStats(results)
	stats.Fprintln(wInfo)

	// prints the results in json format
	bytes, _ := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	wOutput.Write(bytes)

	return nil
}

func printDryRunInfo(w io.Writer, flags *Flags, cfg *Config) error {

	fmt.Fprintf(w, "%s: Dry Run\n", flags.fullname)

	// prints config file info
	fmt.Fprintln(w, cfg.cfi)

	if cfg.Database != "" {
		fmt.Fprintf(w, "Database: %q\n", cfg.Database)
	}
	fmt.Fprintf(w, "Mode: %q\n", cfg.Mode)
	sis := cfg.SourceIsinsList()
	fmt.Fprint(w, "Tasks: ", jsonString(sis))
	return nil
}
