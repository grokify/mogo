package maputil

import (
	"sort"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
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

// Counts returns key counts regardless of value.
func (msmsi MapStringMapStringInt) Counts(sep string) (int, int) {
	count1Info := map[string]int{}
	count2Info := map[string]int{}
	if len(sep) == 0 {
		sep = "~~~"
	}
	for str1, map1 := range msmsi {
		for str2 := range map1 {
			count1Info[str1] = 1
			count2Info[str1+sep+str2] = 1
		}
	}
	return len(count1Info), len(count2Info)
}

// CountsWithVal returns key counts with the desired value.
func (msmsi MapStringMapStringInt) CountsWithVal(wantVal int, sep string) (int, int) {
	count1Info := map[string]int{}
	count2Info := map[string]int{}
	if len(sep) == 0 {
		sep = "~~~"
	}
	for str1, map1 := range msmsi {
		for str2, val := range map1 {
			if wantVal == val {
				count1Info[str1] = 1
				count2Info[str1+sep+str2] = 1
			}
		}
	}
	return len(count1Info), len(count2Info)
}

// MapStringMapStringIntFuncExactMatch returns a match function
// that matches an exact vaule.
func MapStringMapStringIntFuncExactMatch(valWant int) func(string, string, int) bool {
	return func(str1, str2 string, val int) bool {
		return valWant == val
	}
}

// MapStringMapStringIntFuncIncludeAll returns match function
// that matches all values.
func MapStringMapStringIntFuncIncludeAll(str1, str2 string, val int) bool {
	return true
}

// Flatten returns the values in the `MapStringMapStringInt` in an string slice.
func (msmsi MapStringMapStringInt) Flatten(prefix, sep string, fnInclude func(string, string, int) bool, dedupe, sortResults bool) []string {
	outvals := []string{}
	for str1, map1 := range msmsi {
		for str2, val := range map1 {
			if !fnInclude(str1, str2, val) {
				continue
			}
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
