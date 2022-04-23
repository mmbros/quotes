// Package configfile provides the configfile.SourceInfo struct
// that gives the config file properties: path, format and source.
package configfile

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type enumConfigSource int

// Configuration file source enum.
// Note that command-line or environment source can set an empty config file.
const (
	configNone enumConfigSource = iota
	configFromCommandLine
	configFromEnvironment
	configFromDefaults
)

// SourceInfo representes the information about the config file.
type SourceInfo struct {
	path   string           // path of the config file
	format string           // format of config file (json, yaml, toml)
	source enumConfigSource // source of the config file
}

// Path returns the confg file epath
func (i *SourceInfo) Path() string {
	return i.path
}

// Format returns the (lowercase) config file format base the following criteria:
// 1. the explicitly passed format, or
// 2. the extension of the path without the leading "."
// It returns "" in case of nil object.
func (i *SourceInfo) Format() string {
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

func (n enumConfigSource) String() string {
	switch n {
	case configNone:
		return "none"
	case configFromCommandLine:
		return "command-line"
	case configFromEnvironment:
		return "environment"
	case configFromDefaults:
		return "default paths"
	default:
		return "unknown source"
	}
}

// String returns a representation of the configfile.Info object
func (i *SourceInfo) String() string {
	var s string

	if i == nil || i.source == configNone {
		s = "not defined"
	} else if i.path == "" {
		s = fmt.Sprintf("skipped by %s", i.source)
	} else if i.format != "" {
		s = fmt.Sprintf("using %q (type=%q) from %s\n", i.path, i.format, i.source)
	} else {
		s = fmt.Sprintf("using %q from %s\n", i.path, i.source)
	}
	return "Configuration file: " + s
}

// sanitizeEnvKey return the input string modified:
// - UPPERCASE letters
// - substitute every not digit or letter char with "_"
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

func newInfoFromCommandLine(path, format string) (*SourceInfo, error) {
	var err error
	info := &SourceInfo{
		path:   path,
		format: format,
		source: configFromCommandLine,
	}
	if path != "" {
		// if it is not blank, check if exists
		_, err = os.Stat(path)
		if err != nil {
			err = fmt.Errorf("configfile: command-line config file: %w", err)
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
func newInfoFromEnv(appname, keyPrefix string) (*SourceInfo, error) {

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
	info := &SourceInfo{
		path:   keyConfigValue,
		format: os.Getenv(keyConfigTypeName),
		source: configFromEnvironment,
	}

	if info.path != "" {
		_, err = os.Stat(info.path)
		if err != nil {
			err = fmt.Errorf("configfile: environment config file: %w", err)
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

func newInfoFromDefaults(appname, home string) (*SourceInfo, error) {

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
				info := &SourceInfo{
					path:   path,
					format: ext,
					source: configFromDefaults,
				}
				return info, nil
			}
		}
	}

	return &SourceInfo{source: configNone}, nil
}

// NewSourceInfo init the config file information struct from (in order):
// 1. command line
// 2. environment variables
// 3. defaults
//
// Arguments:
//   appname:   name of the app
//   keyPrefix: prefix of the environment variables. If empty, apname will be used
//   path:      path of the config file passed as command line argument. It can be empty
//   format:    format of the config file passed as command line argument. It can be empty
//   passed:    true if path is passed as command line argument
func NewSourceInfo(appname, keyPrefix, path, format string, passed bool) (*SourceInfo, error) {

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
