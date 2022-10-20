package maputil

import (
	"sort"

	"github.com/grokify/mogo/type/stringsutil"
)

type MapStringSlice map[string][]string

func NewMapStringSlice() MapStringSlice {
	return MapStringSlice{}
}

// Add adds a key and value to the `map[string][]string`. It will
// panic on a nil struct, so do not preceed with `var mss MapStringSlice`
func (mss MapStringSlice) Add(key, value string) {
	if _, ok := mss[key]; !ok {
		mss[key] = []string{value}
	} else {
		mss[key] = append(mss[key], value)
	}
}

func (mss MapStringSlice) Sort(dedupe bool) {
	for key, vals := range mss {
		if dedupe {
			vals = stringsutil.Dedupe(vals)
		}
		sort.Strings(vals)
		mss[key] = vals
	}
}

// KeysByValueCounts returns a `map[int][]string` where the key is the
// count of values and the values are the keys with that value count.
func (mss MapStringSlice) KeysByValueCounts() map[int][]string {
	byCount := map[int][]string{}
	for key, vals := range mss {
		count := len(vals)
		byCount[count] = append(byCount[count], key)
	}
	return byCount
}
