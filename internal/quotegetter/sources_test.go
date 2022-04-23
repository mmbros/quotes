package quotegetter

import (
	"net/http"
	"testing"
)

func TestString(t *testing.T) {

	dummy := func(name string, client *http.Client) QuoteGetter { return nil }

	tests := []struct {
		title string
		names []string
		want  string
	}{
		{
			title: "empty",
		},
		{
			title: "one",
			names: []string{"X"},
			want:  "X",
		},
		{
			title: "two",
			names: []string{"X", "AAA"},
			want:  "AAA, X",
		},
		{
			title: "three",
			names: []string{"000", "X", "AAA"},
			want:  "000, AAA, X",
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			// init sources
			sources := Sources{}
			for _, n := range tt.names {
				sources[n] = dummy
			}
			// check
			got := sources.String()
			if got != tt.want {
				t.Errorf("want %q, got %q!", tt.want, got)
			}
		})
	}
}
