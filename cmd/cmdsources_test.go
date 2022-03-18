package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func Test_parseExecSources(t *testing.T) {
	type args struct {
		fullname  string
		arguments []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseExecSources(tt.args.fullname, tt.args.arguments); (err != nil) != tt.wantErr {
				t.Errorf("parseExecSources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_execSources(t *testing.T) {
	tests := []struct {
		name  string
		wantW string
	}{
		{"contains", "morningstarit"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			execSources(w)
			if gotW := w.String(); !strings.Contains(gotW, tt.wantW) {
				t.Errorf("execSources() =\n\t%vdoes not cointains:\n\t%v", gotW, tt.wantW)
			}
		})
	}
}
