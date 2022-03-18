package cmd

import (
	"flag"
	"fmt"

	"github.com/mmbros/quote/internal/quote"
)

func parseExecTor(fullname string, arguments []string) error {
	var cfg *Config

	// parse the arguments
	flags := NewFlags(fullname, fgAppTor)
	flags.SetUsage(usageTor, fullname)

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

	return execTor(flags, cfg)
}

func execTor(flags *Flags, cfg *Config) error {
	if flags.IsPassed(namesConfig) {
		fmt.Printf("Using configuration file %q\n", flags.config)
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
