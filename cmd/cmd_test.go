package cmd_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/mmbros/quote/cmd"
	"github.com/stretchr/testify/assert"
)

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
		"app s --help (short version)": {
			cmdline: "app s --help",
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
			cmd.Execute()

			assert.Contains(t, out.String(), tc.want, "usage does not contain expected string")
		})
	}
}
