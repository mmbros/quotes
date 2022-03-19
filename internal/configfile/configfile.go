package configfile

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mmbros/quote/internal/errors"
)

type EnumSource int

// Configuration file source enum
// Note that command-line or environment source can set an empty config file.
const (
	srcNone EnumSource = iota
	srcCommandLine
	srcEnvironment
	srcDefaults
)

// Info about the config file
type Info struct {
	path   string     // path of the config file
	format string     // format of config file (json, yaml, toml)
	Source EnumSource // source of the config file
}

// Path returs the confg file epath
func (i *Info) Path() string {
	return i.path
}

// Format returns the (lowercase) data format base the following criteria:
// 1. the explicitly passed format, or
// 2. the extension of the path without the leading "."
// It returns "" in case of nil object.
func (i *Info) Format() string {
	var s string

	if i == nil {
		return s
	}
	if i.format != "" {
		s = i.format
	} else {
		s = filepath.Ext(i.path)
		if len(s) > 1 {
			s = s[1:]
		}
	}

	return strings.ToLower(s)
}

// sanitizeEnvKey ...
func sanitizeEnvKey(key string) string {
	mapping := func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return 'A' + r - 'a'
		case r >= 'A' && r <= 'Z':
			fallthrough
		case r >= '0' && r <= '9':
			return r
		}
		return '_'
	}
	return strings.Map(mapping, key)
}

func newInfoFromCommandLine(path, format string) (*Info, error) {
	var err error
	info := &Info{
		path:   path,
		format: format,
		Source: srcCommandLine,
	}
	if path != "" {
		// if it is not blank, check if exists
		_, err = os.Stat(path)
		if err != nil {
			err = errors.WrapErrorf(err, "configfile: command-line config file: %v", err.Error())
		}
	}

	// command-line cases:
	//   a. config file is blank
	//   b. config file exists
	//   c. config file does not exist
	return info, err
}

// newInfoFromEnv returns the configfile.Info using the environment parameters.
// Returns nil, if the config environment variable is not defined.
func newInfoFromEnv(appname, keyPrefix string) (*Info, error) {

	if keyPrefix == "" {
		keyPrefix = appname
	}
	keyConfigName := sanitizeEnvKey(keyPrefix + "_CONFIG")
	keyConfigValue, keyConfigExists := os.LookupEnv(keyConfigName)

	if !keyConfigExists {
		// config file environment parameter not defined
		return nil, nil
	}

	var err error
	keyConfigTypeName := sanitizeEnvKey(keyPrefix + "_CONFIG_TYPE")
	info := &Info{
		path:   keyConfigValue,
		format: os.Getenv(keyConfigTypeName),
		Source: srcEnvironment,
	}

	if info.path != "" {
		_, err = os.Stat(info.path)
		if err != nil {
			err = errors.WrapErrorf(err, "configfile: environment config file: %v", err.Error())
		}
	}

	// environment cases:
	//   a. config file is blank
	//   b. config file exists
	//   c. config file does not exist
	return info, err
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func newInfoFromDefaults(appname, home string) (*Info, error) {

	astr := []string{
		filepath.Join(home, "."+appname+"."),
		filepath.Join(home, "."+appname, "config."),
		filepath.Join(home, "."+appname, appname+"."),
	}

	for _, base := range astr {
		for _, ext := range []string{"toml", "yaml", "yml", "json"} {
			path := base + ext

			_, err := os.Stat(path)
			if err == nil {
				info := &Info{
					path:   path,
					format: ext,
					Source: srcDefaults,
				}
				return info, nil
			}
		}
	}

	return &Info{Source: srcNone}, nil
}

// NewInfo init the config file information struct from (in order):
// 1. command line
// 2. environment variables
// 3. defaults
//
// Arguments:
//   appname:   name of the app
//   keyPrefix: prefix of the environment variables. If empty, apname will be used
//   path:      path of the config file passed as command line argument. It can be empty
//   format:    format of the config file passed as command line argument
//   passed:    true if path is passed as command line argument
func NewInfo(appname, keyPrefix, path, format string, passed bool) (*Info, error) {

	// 1. checks command-line config file, if passed
	if passed {
		return newInfoFromCommandLine(path, format)
	}

	// 2. checks environment variables
	info, err := newInfoFromEnv(appname, keyPrefix)
	if info != nil {
		return info, err
	}

	// 3. checks defaults config file
	return newInfoFromDefaults(appname, userHomeDir())
}
