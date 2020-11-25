package maputil

import (
	"sort"
	"strings"

	"github.com/grokify/gotilla/type/stringsutil"
)

type MapStrMapStrInt map[string]map[string]int

func (msmsi MapStrMapStrInt) Add(str1, str2 string, val int) {
	if _, ok := msmsi[str1]; !ok {
		msmsi[str1] = map[string]int{}
	}
	if _, ok := msmsi[str1][str2]; !ok {
		msmsi[str1][str2] = 0
	}
	msmsi[str1][str2] = val
}
func (msmsi MapStrMapStrInt) Flatten(prefix, sep string, dedupe, sortResults bool) []string {
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
