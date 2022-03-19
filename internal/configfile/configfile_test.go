package configfile

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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
		i := &Info{c.path, "", srcCommandLine}
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
		i := &Info{c.path, c.format, srcCommandLine}
		got := i.Format()
		assert.Equal(t, c.want, got, c)
	}

	var nilInfo *Info
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
		want    *Info
		wantErr bool
	}{
		{
			"blank path and format",
			args{"", ""},
			&Info{"", "", srcCommandLine},
			false,
		},
		{
			"blank path",
			args{"", "format"},
			&Info{"", "format", srcCommandLine},
			false,
		},
		{
			"path exists",
			args{"/dev/null", "json"},
			&Info{"/dev/null", "json", srcCommandLine},
			false,
		},
		{
			"path not exist",
			args{"/dev/null$@#", "toml"},
			&Info{"/dev/null$@#", "toml", srcCommandLine},
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

		want    *Info
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
			&Info{"/dev/null", "", srcEnvironment},
			false,
		},
		{
			"blank config and config-type defined",
			"", true, "yml", true,
			&Info{"", "yml", srcEnvironment},
			false,
		},
		{
			"config defined but non exists",
			"/dev/null#$%", true, "toml", true,
			&Info{"/dev/null#$%", "toml", srcEnvironment},
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
		want    *Info
		wantErr bool
	}{
		{
			"no config",
			"",
			&Info{"", "", srcNone},
			false,
		},
		{
			".appname.json",
			filepath.Join(home, "."+appname+".json"),
			&Info{filepath.Join(home, "."+appname+".json"), "json", srcDefaults},
			false,
		},
		{
			".appname/config.yaml",
			filepath.Join(home, "."+appname, "config.yaml"),
			&Info{filepath.Join(home, "."+appname, "config.yaml"), "yaml", srcDefaults},
			false,
		},
		{
			".appname/appname.toml",
			filepath.Join(home, "."+appname, appname+".toml"),
			&Info{filepath.Join(home, "."+appname, appname+".toml"), "toml", srcDefaults},
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

		want    *Info
		wantErr bool
	}{
		{
			"cmdline",
			args{pathCmdline, "json", true},
			pathEnv, true,
			"",
			&Info{pathCmdline, "json", srcCommandLine},
			true,
		},
		{
			"cmdline blank",
			args{"", "toml", true},
			pathEnv, true,
			"",
			&Info{"", "toml", srcCommandLine},
			false,
		},
		{
			"env",
			args{},
			pathEnv, true,
			"",
			&Info{pathEnv, "", srcEnvironment},
			true,
		},
		{
			"env blank",
			args{},
			"", true,
			"",
			&Info{"", "", srcEnvironment},
			false,
		},
		{
			"no default found (no home dir)",
			args{},
			"", false,
			pathDefault,
			&Info{"", "", srcNone},
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

			got, err := NewInfo(appname, "app", tt.args.path, tt.args.format, tt.args.passed)
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
