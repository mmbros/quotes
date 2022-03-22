package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const usageVersion = `Usage:
    %s

Prints version informations.

Options:
    -b, --build-options        also print build options
`

// set at compile time with
//   -ldflags="-X 'github.com/mmbros/quotes/cmd.Version=x.y.z' -X 'github.com/mmbros/quotes/cmd.GitCommit=...'"
var (
	Version   string // git tag ...
	GitCommit string // git rev-parse --short HEAD
	GoVersion string // go version
	BuildTime string // when the executable was built
	OsArch    string // uname -s -m
)

func parseExecVersion(fullname string, arguments []string) error {

	// parse the arguments
	flags := NewFlags(fullname, fgAppVersion)
	flags.SetUsage(usageVersion, fullname)

	err := flags.Parse(arguments)

	// handle help
	if err == flag.ErrHelp {
		flags.Usage()
		return nil
	}
	if err == nil {
		// NOTE: flags.dryrun contains the build-options flag
		execVersion(os.Stdout, flags.Appname(), flags.dryrun)
	}
	return err
}

func execVersion(w io.Writer, appname string, extended bool) {
	fmt.Fprintf(w, "%s version %s\n", appname, Version)
	if extended {
		fmt.Fprintf(w, `%s
build date: %s
git commit: %s
os/arch: %s
`, GoVersion, BuildTime, GitCommit, OsArch)
	}
}
