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

func KeyValueSliceCounts[K comparable, V any](m map[K][]V) map[K]int {
	r := map[K]int{}
	for k, vals := range m {
		r[k] = len(vals)
	}
	return r
}

// Flatten converts a `map[string][]string{}` to a `map[string]string{}`. The default is to use the first value
// unless `useLast` is set to true, in which case the last value is used. The default is to use a key with an
// empty string to represent a key with an empty slice, unless `skipEmpty` is set to `true`, in which case the
// key with an empty slice is not represented.
func (mss MapStringSlice) Flatten(useLast, skipEmpty bool) map[string]string {
	simple := map[string]string{}
	for k, vals := range mss {
		if len(vals) == 0 {
			if !skipEmpty {
				simple[k] = ""
			}
		} else if useLast {
			simple[k] = vals[len(vals)-1]
		} else {
			simple[k] = vals[0]
		}
	}
	return simple
}

// FlattenAny converts the results of `Flatten()` to a `map[string]any{}`.`
func (mss MapStringSlice) FlattenAny(useLast, skipEmpty bool) map[string]any {
	msa := map[string]any{}
	mssimple := mss.Flatten(useLast, skipEmpty)
	for k, v := range mssimple {
		msa[k] = v
	}
	return msa
}

func (mss MapStringSlice) Lines(m map[string][]string, keyPrefix, valPrefix string) []string {
	var lines []string
	for k, vals := range mss {
		lines = append(lines, keyPrefix+k)
		for _, v := range vals {
			lines = append(lines, valPrefix+v)
		}
	}
	return lines
}
