package cmd_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/mmbros/quotes/cmd"
	"github.com/stretchr/testify/assert"
)

func Test_GetDryRun(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "temp-app-config-")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	fname := tmpFile.Name()
	defer os.Remove(fname)

	tests := map[string]struct {
		cmdline string
		want    string
	}{
		"app get -n": {
			cmdline: "app get --dry-run",
			want:    `Mode: "A"`,
		},
		"app get -n (with config)": {
			cmdline: "app get -n --config " + fname,
			want:    fname,
		},
	}
	for title, tt := range tests {

		t.Run(title, func(t *testing.T) {

			var out strings.Builder
			flag.CommandLine.SetOutput(&out)

			os.Args = strings.Split(tt.cmdline, " ")
			cmd.Execute(&out)

			// if diff := cmp.Diff(tc.want, out.String(), nil); diff != "" {
			// 	t.Errorf("%s: mismatch (-want +got):\n%s", title, diff)
			// }

			assert.Contains(t, out.String(), tt.want, "usage does not contain expected string")

			// t.Log(out.String())
			// t.FailNow()

		})
	}
}
