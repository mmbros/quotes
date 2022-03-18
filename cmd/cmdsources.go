package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mmbros/quote/internal/quote"
)

func parseExecSources(fullname string, arguments []string) error {

	// parse the arguments
	flags := NewFlags(fullname, fgAppSources)
	flags.SetUsage(usageSources, fullname)

	err := flags.Parse(arguments)

	// handle help
	if err == flag.ErrHelp {
		flags.Usage()
		return nil
	}
	if err != nil {
		return err
	}

	execSources(os.Stdout)
	return nil
}

func execSources(w io.Writer)  {
	sources := quote.Sources()
	fmt.Fprintf(w, "Available sources: \"%s\"\n", strings.Join(sources, "\", \""))

}
