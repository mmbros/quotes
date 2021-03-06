package cmd_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/mmbros/quotes/cmd"
	"github.com/stretchr/testify/assert"
)

func Test_Execute(t *testing.T) {

	testCases := map[string]struct {
		cmdline  string
		wantCode int
		wantErr  string
	}{
		"ok": {
			cmdline:  "app",
			wantCode: 0,
		},
		"ok2": {
			cmdline:  "app -h",
			wantCode: 0,
		},
		"error": {
			cmdline:  "app --not-exists",
			wantCode: 1,
			wantErr:  "flag provided but not defined",
		},
	}

	for title, tc := range testCases {

		t.Run(title, func(t *testing.T) {

			os.Args = strings.Split(tc.cmdline, " ")

			var stdout, stderr strings.Builder

			flag.CommandLine.SetOutput(&stdout)
			gotCode := cmd.Execute(&stderr)

			if gotCode != tc.wantCode {
				t.Errorf("Execute() = %d, want %d", gotCode, tc.wantCode)
			}

			gotErr := stderr.String()
			if tc.wantErr != "" {
				// msg error is expected
				if !strings.Contains(gotErr, tc.wantErr) {
					t.Errorf("Execute(): expected error %q does not contain %q", gotErr, tc.wantErr)
				}
			} else {
				// msg error is NOT expected
				if stderr.Len() > 0 {
					t.Errorf("Execute(): unexpected error %q", gotErr)
				}
			}
		})
	}
}

func Test_Help(t *testing.T) {

	testCases := map[string]struct {
		cmdline string
		want    string
	}{
		"app (only)": {
			cmdline: "app",
			want:    "app <command>",
		},
		"app --help": {
			cmdline: "app --help",
			want:    "app <command>",
		},
		"app -h": {
			cmdline: "app --help",
			want:    "app <command>",
		},
		"app get -h": {
			cmdline: "app get -h",
			want:    "app get",
		},
		"app g --help (short version)": {
			cmdline: "app g --help",
			want:    "app get",
		},
		"app sources -h": {
			cmdline: "app sources -h",
			want:    "app sources",
		},
		"app so --help (short version)": {
			cmdline: "app so --help",
			want:    "app sources",
		},
		"app tor -h": {
			cmdline: "app tor -h",
			want:    "app tor",
		},
		"app t --help (short version)": {
			cmdline: "app t --help",
			want:    "app tor",
		},
	}

	for title, tc := range testCases {

		t.Run(title, func(t *testing.T) {

			var out strings.Builder
			flag.CommandLine.SetOutput(&out)

			os.Args = strings.Split(tc.cmdline, " ")
			cmd.Execute(&out)

			assert.Contains(t, out.String(), tc.want, "usage does not contain expected string")
		})
	}
}
