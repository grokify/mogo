package maputil

import (
	"sort"
	"strings"

	"github.com/grokify/gotilla/type/stringsutil"
)

// MapStringMapStringInt represents a `map[string]map[string]int`
type MapStringMapStringInt map[string]map[string]int

// Set supports setting values to `MapStringMapStringInt` easily.
func (msmsi MapStringMapStringInt) Set(str1, str2 string, val int) {
	if _, ok := msmsi[str1]; !ok {
		msmsi[str1] = map[string]int{}
	}
	if _, ok := msmsi[str1][str2]; !ok {
		msmsi[str1][str2] = 0
	}
	msmsi[str1][str2] = val
}

// Add supports adding values to `MapStringMapStringInt` easily.
func (msmsi MapStringMapStringInt) Add(str1, str2 string, val int) {
	if _, ok := msmsi[str1]; !ok {
		msmsi[str1] = map[string]int{}
	}
	if _, ok := msmsi[str1][str2]; !ok {
		msmsi[str1][str2] = 0
	}
	msmsi[str1][str2] += val
}

// Flatten returns the values in the `MapStringMapStringInt` in an string slice.
func (msmsi MapStringMapStringInt) Flatten(prefix, sep string, dedupe, sortResults bool) []string {
	outvals := []string{}
	for str1, map1 := range msmsi {
		for str2 := range map1 {
			parts := []string{}
			if len(prefix) > 0 {
				parts = append(parts, prefix)
			}
			parts = append(parts, str1, str2)
			outval := strings.Join(parts, sep)
			outvals = append(outvals, outval)
		}
	}
	if dedupe {
		outvals = stringsutil.Dedupe(outvals)
	}
	if sortResults {
		sort.Strings(outvals)
	}
	return outvals
}
