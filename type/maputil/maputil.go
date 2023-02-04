package maputil

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/grokify/mogo/sort/sortutil"
	"golang.org/x/exp/constraints"
)

// StringKeys takes a map where the keys are strings and reurns a slice of key names.
// An optional transform function, `xf`, can be supplied along with an option to sort
// the results. If both transform and sort are requested, the sort is performed on the
// transformed strings.
func StringKeys[V any](m map[string]V, xf func(s string) string) []string {
	keys := []string{}
	for k := range m {
		if xf != nil {
			k = xf(k)
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Keys returns a list of sorted keys.
func Keys[K constraints.Ordered, V any](m map[K]V) []K {
	keys := []K{}
	for k := range m {
		keys = append(keys, k)
	}
	sortutil.Slice(keys)
	return keys
}

// KeysExist checks to verify if a set of keys exists within a map with string keys.
// If `requireAll` is set, then all keys must be present for the function to return `true`.
// If `requireAll` is not set, then only one key must exist for the function to return `true`.
func KeysExist[K comparable, V any](m map[K]V, keys []K, requireAll bool) bool {
	for _, k := range keys {
		_, ok := m[k]
		if requireAll && !ok {
			return false
		} else if !requireAll && ok {
			return true
		}
	}
	if requireAll {
		return true
	} else {
		return false
	}
}

// Values returns a string slice of sorted values.
func Values[K comparable, V constraints.Ordered](m map[K]V) []V {
	vals := []V{}
	for _, val := range m {
		vals = append(vals, val)
	}
	sortutil.Slice(vals)
	return vals
}

// NumberValuesAverage returns a `float64` average of a map's values.
func NumberValuesAverage[K comparable, V constraints.Float | constraints.Integer](m map[K]V) float64 {
	if len(m) == 0 {
		return 0
	}
	sum := float64(0)
	for _, v := range m {
		sum += float64(v)
	}
	return sum / float64(len(m))
}

var ErrMapDuplicateValues = errors.New("map has duplicate values")

// UniqueValues returns `true` if the number of unique values is the same as the number of keys.
func UniqueValues[C comparable](m map[C]C) bool {
	rev := map[C]C{}
	for k, v := range m {
		if _, ok := rev[v]; ok {
			return false
		}
		rev[v] = k
	}
	return len(rev) == len(m)
}

func DuplicateValues[C comparable](m map[C]C) map[C][]C {
	dupes := map[C][]C{}
	for k, v := range m {
		if _, ok := dupes[v]; !ok {
			dupes[v] = []C{}
		}
		dupes[v] = append(dupes[v], k)
	}
	for k, v := range dupes {
		if len(v) == 1 {
			delete(dupes, k)
		}
	}
	return dupes
}

func ValueIntOrDefault[K comparable, V constraints.Integer](m map[K]V, key K, def V) V {
	if val, ok := m[key]; ok {
		return val
	}
	return def
}

func ValueStringOrDefault[K comparable](m map[K]any, key K, def string) string {
	if val, ok := m[key]; ok {
		if str, ok2 := val.(string); ok2 {
			return str
		}
	}
	return def
}

func MapSSToKeyValues(kvs map[string]string, sep string) string {
	pairs := []string{}
	for k, v := range kvs {
		k = strings.Trim(k, sep)
		v = strings.Trim(v, sep)
		if len(k) > 0 {
			pairs = append(pairs, k+"="+v)
		}
	}
	return strings.Join(pairs, sep)
}

type MapInt64Int64 map[int64]int64

func (m MapInt64Int64) KeysSorted() []int64 {
	return Keys(m)
}

func (m MapInt64Int64) ValuesSortedByKeys() []int64 {
	vals := []int64{}
	keys := m.KeysSorted()
	for _, k := range keys {
		if v, ok := m[k]; ok {
			vals = append(vals, v)
		} else {
			panic(fmt.Sprintf("key not found [%d]", k))
		}
	}
	return vals
}
