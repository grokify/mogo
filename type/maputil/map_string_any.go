package maputil

import (
	"errors"

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
