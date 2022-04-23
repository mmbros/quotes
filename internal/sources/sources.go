package sources

import (
	"sort"
	"strings"

	"github.com/mmbros/quotes/internal/quotegetter"
)

type QuoteGetterSources struct {
	items map[string]quotegetter.NewQuoteGetterFunc
}

func NewQuoteGetterSources() *QuoteGetterSources {
	return &QuoteGetterSources{
		items: make(map[string]quotegetter.NewQuoteGetterFunc),
	}
}

func (ss *QuoteGetterSources) Add(name string, fn quotegetter.NewQuoteGetterFunc) {
	ss.items[name] = fn
}

func (ss *QuoteGetterSources) Get(name string) quotegetter.NewQuoteGetterFunc {
	return ss.items[name]
}

func (ss *QuoteGetterSources) Exists(name string) bool {
	return ss.items[name] != nil
}

// Names returns a list of the names of the QuoteGetters.
// Note the list is not sorted.
func (ss *QuoteGetterSources) Names() []string {
	list := make([]string, 0, len(ss.items))
	for name := range ss.items {
		list = append(list, name)
	}
	return list
}

// String returns the sorted comma separated list of the QuoteGetters names.
func (ss *QuoteGetterSources) String() string {
	names := ss.Names()
	sort.Strings(names)
	return strings.Join(names, ", ")
}
