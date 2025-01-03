package maputil

import (
	"errors"

	"github.com/grokify/mogo/sort/sortutil"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

var (
	ErrKeyNotExist = errors.New("key not found")
	ErrNotString   = errors.New("value not string")
)

// MapStrAny represents a `map[string]any`
type MapStrAny map[string]any

func (msa MapStrAny) ValueString(k string, errOnNotExist bool) (string, error) {
	v, ok := msa[k]
	if !ok {
		if errOnNotExist {
			return "", ErrKeyNotExist
		} else {
			return "", nil
		}
	}
	s, ok := v.(string)
	if !ok {
		return "", ErrNotString
	}
	return s, nil
}

func KeysEqual[K constraints.Ordered, V any](m1, m2 map[K]V) bool {
	m1Keys, m2Keys := Keys(m1), Keys(m2)
	return slices.Equal(m1Keys, m2Keys)
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
