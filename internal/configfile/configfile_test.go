package configfile

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Example_sanitizeEnvKey() {
	fmt.Println(sanitizeEnvKey("appname"))
	fmt.Println(sanitizeEnvKey("app-name"))
	fmt.Println(sanitizeEnvKey("az_AZ-09"))
	fmt.Println(sanitizeEnvKey("a=p_p n+a(m)e"))

	// Output:
	// APPNAME
	// APP_NAME
	// AZ_AZ_09
	// A_P_P_N_A_M_E
}

func Test_Path(t *testing.T) {

	cases := []struct {
		path string
	}{
		{"/home/user/config.YAML"},
		{"config.toml"},
		{"config"},
	}

	for _, c := range cases {
		i := &SourceInfo{c.path, "", configFromCommandLine}
		got := i.Path()
		want := c.path
		assert.Equal(t, want, got, c)
	}
}

func Test_Format(t *testing.T) {

	cases := []struct {
		path   string
		format string
		want   string
	}{
		{"/home/user/config.YAML", "", "yaml"},
		{"config.toml", "EXT", "ext"},
		{"config", "", ""},
	}

	for _, c := range cases {
		i := &SourceInfo{c.path, c.format, configFromCommandLine}
		got := i.Format()
		assert.Equal(t, c.want, got, c)
	}

	var nilInfo *SourceInfo
	got := nilInfo.Format()
	assert.Equal(t, "", got, "nil Info")
}

func Test_newInfoFromCommandLine(t *testing.T) {
	type args struct {
		path   string
		format string
	}
	tests := []struct {
		name    string
		args    args
		want    *SourceInfo
		wantErr bool
	}{
		{
			"blank path and format",
			args{"", ""},
			&SourceInfo{"", "", configFromCommandLine},
			false,
		},
		{
			"blank path",
			args{"", "format"},
			&SourceInfo{"", "format", configFromCommandLine},
			false,
		},
		{
			"path exists",
			args{"/dev/null", "json"},
			&SourceInfo{"/dev/null", "json", configFromCommandLine},
			false,
		},
		{
			"path not exist",
			args{"/dev/null$@#", "toml"},
			&SourceInfo{"/dev/null$@#", "toml", configFromCommandLine},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newInfoFromCommandLine(tt.args.path, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("newInfoFromCommandLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newInfoFromCommandLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newInfoFromEnv(t *testing.T) {
	const (
		appname = "app"
	)
	keyConfig := sanitizeEnvKey(appname + "_CONFIG")
	keyConfigType := sanitizeEnvKey(appname + "_CONFIG_TYPE")

	tests := []struct {
		name string

		valConfig     string
		defConfig     bool
		valConfigType string
		defConfigType bool

		want    *SourceInfo
		wantErr bool
	}{
		{
			"env not defined",
			"", false, "", false,
			nil,
			false,
		},
		{
			"config not defined and config-type defined",
			"", false, "yml", true,
			nil,
			false,
		},
		{
			"config defined",
			"/dev/null", true, "", false,
			&SourceInfo{"/dev/null", "", configFromEnvironment},
			false,
		},
		{
			"blank config and config-type defined",
			"", true, "yml", true,
			&SourceInfo{"", "yml", configFromEnvironment},
			false,
		},
		{
			"config defined but non exists",
			"/dev/null#$%", true, "toml", true,
			&SourceInfo{"/dev/null#$%", "toml", configFromEnvironment},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// handle environment
			os.Clearenv()
			if tt.defConfig {
				os.Setenv(keyConfig, tt.valConfig)
			}
			if tt.defConfigType {
				os.Setenv(keyConfigType, tt.valConfigType)
			}

			got, err := newInfoFromEnv(appname, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("newInfoFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newInfoFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userHomeDir(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"linux",
			os.Getenv("HOME"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := userHomeDir(); got != tt.want {
				t.Errorf("userHomeDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func fileTouchAll(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {

		err = os.MkdirAll(filepath.Dir(path), 0700) // Create your file
		if err != nil {
			return nil
		}
		file, err := os.Create(path)
		if err != nil {
			file.Close()
		}
	}
	return err
}

func Test_newInfoFromDefaults(t *testing.T) {

	const appname = "quotes-temp-config"

	home := filepath.Join(os.TempDir(), "quotes")

	tests := []struct {
		name    string
		path    string
		want    *SourceInfo
		wantErr bool
	}{
		{
			"no config",
			"",
			&SourceInfo{"", "", configNone},
			false,
		},
		{
			".appname.json",
			filepath.Join(home, "."+appname+".json"),
			&SourceInfo{filepath.Join(home, "."+appname+".json"), "json", configFromDefaults},
			false,
		},
		{
			".appname/config.yaml",
			filepath.Join(home, "."+appname, "config.yaml"),
			&SourceInfo{filepath.Join(home, "."+appname, "config.yaml"), "yaml", configFromDefaults},
			false,
		},
		{
			".appname/appname.toml",
			filepath.Join(home, "."+appname, appname+".toml"),
			&SourceInfo{filepath.Join(home, "."+appname, appname+".toml"), "toml", configFromDefaults},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.path != "" {
				err := fileTouchAll(tt.path)
				if err != nil {
					t.Errorf("fileTouchAll() error = %v", err)
					t.Fail()
				} else {
					defer func() {
						os.Remove(tt.path)
					}()
				}
			}

			got, err := newInfoFromDefaults(appname, home)
			if (err != nil) != tt.wantErr {
				t.Errorf("newInfoFromDefaults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newInfoFromDefaults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInfo(t *testing.T) {

	const (
		appname     = "quotes-temp-config"
		pathCmdline = "/cmdline/config/not/exists"
		pathEnv     = "/env/config/not/exists"
		pathDefault = "/def/config/not/exists"
	)

	keyConfig := sanitizeEnvKey("APP_CONFIG")
	// home := filepath.Join(os.TempDir(), "quotes")

	type args struct {
		path   string
		format string
		passed bool
	}
	tests := []struct {
		name string
		args args

		envValue  string
		envSetted bool
		defPath   string

		want    *SourceInfo
		wantErr bool
	}{
		{
			"cmdline",
			args{pathCmdline, "json", true},
			pathEnv, true,
			"",
			&SourceInfo{pathCmdline, "json", configFromCommandLine},
			true,
		},
		{
			"cmdline blank",
			args{"", "toml", true},
			pathEnv, true,
			"",
			&SourceInfo{"", "toml", configFromCommandLine},
			false,
		},
		{
			"env",
			args{},
			pathEnv, true,
			"",
			&SourceInfo{pathEnv, "", configFromEnvironment},
			true,
		},
		{
			"env blank",
			args{},
			"", true,
			"",
			&SourceInfo{"", "", configFromEnvironment},
			false,
		},
		{
			"no default found (no home dir)",
			args{},
			"", false,
			pathDefault,
			&SourceInfo{"", "", configNone},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// handle default config file
			if tt.defPath != "" {
				err := fileTouchAll(tt.defPath)
				if err != nil {
					t.Errorf("fileTouchAll() error = %v", err)
					t.Fail()
				} else {
					defer func() {
						os.Remove(tt.defPath)
					}()
				}
			}

			// handle environment
			os.Clearenv()
			if tt.envSetted {
				os.Setenv(keyConfig, tt.envValue)
			}

			got, err := NewSourceInfo(appname, "app", tt.args.path, tt.args.format, tt.args.passed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_String(t *testing.T) {
	tests := map[string]struct {
		info *SourceInfo
		want string
	}{
		"nil":  {nil, "not defined"},
		"zero": {&SourceInfo{}, "not defined"},

		"none empty":     {&SourceInfo{"", "", configNone}, "not defined"},
		"none type only": {&SourceInfo{"", "toml", configNone}, "not defined"},

		"cmdline empty":       {&SourceInfo{"", "", configFromCommandLine}, "skipped by command"},
		"cmdline type only":   {&SourceInfo{"", "json", configFromCommandLine}, "skipped by command"},
		"cmdline path":        {&SourceInfo{"path/to/file.toml", "", configFromCommandLine}, `"path/to/file.toml" from command-line`},
		"cmdline path & type": {&SourceInfo{"path/to/file.toml", "yaml", configFromCommandLine}, `"path/to/file.toml" (type="yaml") from command-line`},

		"env empty":       {&SourceInfo{"", "", configFromEnvironment}, "skipped by environment"},
		"env type only":   {&SourceInfo{"", "json", configFromEnvironment}, "skipped by environment"},
		"env path":        {&SourceInfo{"path/to/file.ext", "", configFromEnvironment}, `"path/to/file.ext" from environment`},
		"env path & type": {&SourceInfo{"path/to/file.ext", "toml", configFromEnvironment}, `"path/to/file.ext" (type="toml") from environment`},

		"defaults path": {&SourceInfo{"path/to/file.ext", "", configFromDefaults}, `"path/to/file.ext" from def`},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			got := tt.info.String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("Fprintln() = %q does not contain %q", got, tt.want)
			}

		})
	}
}
