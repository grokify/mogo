package maputil

import (
	"net/url"
	"sort"
	"strings"

	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

// type MapStringSlice map[string][]string

type MapStringSlice url.Values

func NewMapStringSlice() MapStringSlice {
	return MapStringSlice{}
}

// Add adds a key and value to the `map[string][]string`. It will
// panic on a nil struct, so do not precede with `var mss MapStringSlice`
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
			vals = slicesutil.Dedupe(vals)
		}
		sort.Strings(vals)
		mss[key] = vals
	}
}

func (mss MapStringSlice) Keys() []string {
	return StringKeys(mss, nil)
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

// MapStringSliceCondenseSpace will trim spaces for keys and values, and remove empty values. It will
// also optionally dedupe and sort values.
func MapStringSliceCondenseSpace(m map[string][]string, dedupeVals, sortVals bool) map[string][]string {
	new := map[string][]string{}
	for key, vals := range m {
		new[strings.TrimSpace(key)] = stringsutil.SliceCondenseSpace(vals, dedupeVals, sortVals)
	}
	return new
}
