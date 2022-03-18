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
			cmdline: "appNAME",
			want:    "appNAME <command>",
		},
		"app --help": {
			cmdline: "appNAME --help",
			want:    "appNAME <command>",
		},
		"app get -h": {
			cmdline: "appNAME get -h",
			want:    "appNAME get",
		},
		"app g --help (short version)": {
			cmdline: "appNAME g --help",
			want:    "appNAME get",
		},
	}
	for title, tc := range testCases {

		t.Run(title, func(t *testing.T) {

			var out strings.Builder
			flag.CommandLine.SetOutput(&out)

			os.Args = strings.Split(tc.cmdline, " ")
			cmd.Execute()

			// if diff := cmp.Diff(tc.want, out.String(), nil); diff != "" {
			// 	t.Errorf("%s: mismatch (-want +got):\n%s", title, diff)
			// }

			assert.Contains(t, out.String(), tc.want, "usage does not contain expected string")

		})
	}
}
