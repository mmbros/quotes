package quotegetter

import (
	"net/http"
	"sort"
	"strings"
)

// NewQuoteGetterFunc creates a new QuoteGetter from Name and http.Client.
type NewQuoteGetterFunc func(string, *http.Client) QuoteGetter

// Sources maps the source name to the corrisponing NewQuoteGetterFunc.
type Sources map[string]NewQuoteGetterFunc

// Exists returns if the source exists.
func (sources Sources) Exists(name string) bool {
	return sources[name] != nil
}

// Names returns a list of the names of the QuoteGetters.
// Note the list is not sorted.
func (sources Sources) Names() []string {
	list := make([]string, 0, len(sources))
	for name := range sources {
		list = append(list, name)
	}
	return list
}

// String returns the sorted comma separated list of the QuoteGetters names.
func (sources Sources) String() string {
	names := sources.Names()
	sort.Strings(names)
	return strings.Join(names, ", ")
}
