package cmd_test

// Lisa Melone
// lisa.melone@wail.ch
//
// andiamo SPE$$0 a pescare
//
// https://woelklimail.com/

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/mmbros/quote/cmd"
	"github.com/stretchr/testify/assert"
)

func Test_GetDryRun(t *testing.T) {

	testCases := map[string]struct {
		cmdline string
		want    string
	}{
		"app get -n": {
			cmdline: "app get --dry-run",
			want:    `Mode: "1"`,
		},
		"app get -n (with config)": {
			cmdline: "app get -n --config /home/mau/mauro.quote.yaml",
			want:    "Using configuration file",
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

			// t.Log(out.String())
			// t.FailNow()

		})
	}
}
