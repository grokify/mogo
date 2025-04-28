package maputil

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/grokify/mogo/sort/sortutil"
	"github.com/grokify/mogo/strconv/strconvutil"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

var (
	ErrKeyNotExist = errors.New("key not found")
	ErrNotString   = errors.New("value not string")
)

// MapStringAny represents a `map[string]any`
type MapStringAny map[string]any

func (msa MapStringAny) MustValueString(k string, def string) string {
	if v, err := msa.ValueString(k); err != nil {
		return def
	} else {
		return v
	}
}

func (msa MapStringAny) ValueString(k string) (string, error) {
	if v, ok := msa[k]; !ok {
		return "", ErrKeyNotExist
	} else {
		return strconvutil.AnyToString(v), nil
	}
}

func (msa MapStringAny) MustMarshal() []byte {
	if b, err := json.Marshal(msa); err != nil {
		panic(err)
	} else {
		return b
	}
}

type MapStringAnys []MapStringAny

func (msas MapStringAnys) UniqueKeys() []string {
	var keys []string
	m := map[string]int{}
	for _, msa := range msas {
		for k := range msa {
			if _, ok := m[k]; ok {
				continue
			} else {
				keys = append(keys, k)
				m[k]++
			}
		}
	}
	sort.Strings(keys)
	return keys
}

func KeysEqual[K constraints.Ordered, V any](m1, m2 map[K]V) bool {
	return slices.Equal(Keys(m1), Keys(m2))
}

func KeysSubtract[K constraints.Ordered, V any](m1, m2 map[K]V) []K {
	var out []K
	for kx := range m1 {
		if _, ok := m2[kx]; ok {
			continue
		}
		out = append(out, kx)
	}
	slices.Sort(out)
	return out
}

func Values[K constraints.Ordered, V any](m map[K]V) []V {
	keys := Keys(m)
	var vals []V
	for _, k := range keys {
		if v, ok := m[k]; !ok {
			panic("key not found")
		} else {
			vals = append(vals, v)
		}
	}
	return vals
}

func ValuesByKeys[K comparable, V any](m map[K]V, keys []K, def V) []V {
	var out []V
	for _, k := range keys {
		if v, ok := m[k]; ok {
			out = append(out, v)
		} else {
			out = append(out, def)
		}
	}
	return out
}

// ValuesSorted returns a string slice of sorted values.
func ValuesSorted[K comparable, V constraints.Ordered](m map[K]V) []V {
	var vals []V
	for _, val := range m {
		vals = append(vals, val)
	}
	sortutil.Slice(vals)
	return vals
}

// ValuesSorted returns a string slice of sorted values.
func ValuesString[K comparable](m map[K]string) []string {
	var vals []string
	for _, val := range m {
		vals = append(vals, val)
	}
	sort.Strings(vals)
	return vals
}
