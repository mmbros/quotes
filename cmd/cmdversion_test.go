package cmd

import (
	"bytes"
	"regexp"
	"testing"
)

func Test_parseExecVersion(t *testing.T) {
	type args struct {
		fullname  string
		arguments []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"app version", args{"app version", []string{}}, false},
		{"app version --build-options", args{"app version", []string{"--build-options"}}, false},
		{"app v -b", args{"app v", []string{"-b"}}, false},
		{"app v -x", args{"app v", []string{"-x"}}, true},
		{"app v --help", args{"app v", []string{"--help"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseExecVersion(tt.args.fullname, tt.args.arguments); (err != nil) != tt.wantErr {
				t.Errorf("parseExecVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_execVersion(t *testing.T) {
	type args struct {
		appname  string
		extended bool
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{"simple", args{"app", false}, "^app version.*\n$"},
		{"extended", args{"app", true}, "^app version.*\n.*\n.*\n.*\nos/arch:.*\n$"},
	}

	// Compile this regular expression.

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			execVersion(w, tt.args.appname, tt.args.extended)

			reW := regexp.MustCompile(tt.wantW)

			if gotW := w.String(); !reW.MatchString(gotW) {
				t.Errorf("execVersion() = %v, doen not match %v", gotW, tt.wantW)
			}
		})
	}
}
