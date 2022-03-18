package cmd

import (
	"flag"
	"fmt"

	"github.com/mmbros/flagx"
	"github.com/mmbros/quote/internal/quote"
)

func parseExecTor(fullname string, arguments []string) error {
	var cfg *Config

	// parse the arguments
	args, err := parseTor(fullname, arguments)

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

	return execTor(args, cfg)
}

func parseTor(fullname string, arguments []string) (*Args, error) {
	// it is used a module level declaration for test porpouses.
	// normally do: args := &appArgs{}
	args := &Args{}

	fs := flag.NewFlagSet(fullname, flag.ContinueOnError)
	// use the same output as flag.CommandLine
	fs.SetOutput(flag.CommandLine.Output())

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), usageTor, fullname)
	}
	flagx.AliasedStringVar(fs, &args.config, namesConfig, "", "")
	flagx.AliasedStringVar(fs, &args.configType, namesConfigType, "", "")
	flagx.AliasedStringVar(fs, &args.proxy, namesProxy, "", "")

	err := fs.Parse(arguments)

	return args, err
}

func execTor(args *Args, cfg *Config) error {
	if args.IsPassed(namesConfig) {
		fmt.Printf("Using configuration file %q\n", args.config)
	}
	proxy := cfg.Proxy
	// proxy = "x://\\"
	fmt.Printf("Checking Tor connection with proxy %q\n", proxy)
	_, msg, err := quote.TorCheck(proxy)
	if err == nil {
		// ok checking Tor network:
		// prints the result: it can be ok or ko
		fmt.Println(msg)
	}
	return err
}
