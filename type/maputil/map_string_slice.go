package maputil

import (
	"sort"

	"github.com/grokify/mogo/type/stringsutil"
)

type MapStringSlice map[string][]string

func (mss MapStringSlice) Add(key, value string) {
	mss[key] = append(mss[key], value)
}

func (mss MapStringSlice) Sort(dedupe bool) {
	for k, vals := range mss {
		if dedupe {
			vals = stringsutil.Dedupe(vals)
		}
		sort.Strings(vals)
		mss[k] = vals
	}
}
