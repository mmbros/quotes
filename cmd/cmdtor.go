package cmd

import (
	"flag"
	"fmt"

	"github.com/mmbros/quotes/internal/tor"
)

const usageTor = `Usage:
    %s [options]

Checks if Tor network will be used to get the quote.
To use the Tor network the proxy must be defined through:
    1. proxy argument parameter
    2. proxy config file parameter
    3. HTTP_PROXY, HTTPS_PROXY and NOPROXY enviroment variables.

Options:
    -c, --config      path    config file
        --config-type string  used if config file does not have the extension in the name;
                              accepted values are: YAML, TOML and JSON 
    -p, --proxy       url     proxy to test the Tor network
`

func parseExecTor(fullname string, arguments []string) error {
	var cfg *Config

	// parse the arguments
	flags := NewFlags(fullname, fgAppTor)
	flags.SetUsage(usageTor, fullname)

	err := flags.Parse(arguments)

	// handle help
	if err == flag.ErrHelp {
		// clear error
		// note: usage already showed internally
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

	return execTor(flags, cfg)
}

func execTor(flags *Flags, cfg *Config) error {

	// if flags.IsPassed(namesConfig) {
	// 	fmt.Printf("Using configuration file %q\n", flags.config)
	// }
	fmt.Println(cfg.cfi)
	proxy := cfg.Proxy

	var proxyEnv string
	if proxy == "" {
		proxyEnv = tor.ProxyFromEnvironment()
	} else {
		proxyEnv = proxy
	}

	fmt.Print("Checking Tor connection with ")
	if proxyEnv == "" {
		fmt.Print("no proxy")
	} else {
		fmt.Printf("proxy %q", proxyEnv)
		if proxy == "" {
			fmt.Print(" from HTTPS_PROXY environment variable")
		}
	}
	fmt.Println(".")

	_, msg, err := tor.Check(proxy)
	if err == nil {
		// ok checking Tor network:
		// prints the result: it can be ok or ko
		fmt.Println(msg)
	}
	return err
}
