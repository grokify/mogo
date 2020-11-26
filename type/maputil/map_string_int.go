package maputil

import (
	"sort"
)

// MapStringInt represents a `map[string]int`
type MapStringInt map[string]int

// Set sets the value of `val` to `key`.
func (msi MapStringInt) Set(key string, val int) {
	msi[key] = val
}

// Add adds the value of `val` to `key`.
func (msi MapStringInt) Add(key string, val int) {
	if _, ok := msi[key]; !ok {
		msi[key] = 0
	}
	msi[key] += val
}

// Keys returns a string slice of the map's keys.
func (msi MapStringInt) Keys(sortKeys bool) []string {
	keys := []string{}
	for key := range msi {
		keys = append(keys, key)
	}
	if sortKeys {
		sort.Strings(keys)
	}
	return keys
}

// MustGet returns a value or a default.
func (msi MapStringInt) MustGet(key string, defaultValue int) int {
	if val, ok := msi[key]; ok {
		return val
	}
	return defaultValue
}

// MinMaxValues returns the minium and maximum values
// of the `map[string]int`.
func (msi MapStringInt) MinMaxValues() (int, int) {
	min := 0
	max := 0
	i := 0
	for _, val := range msi {
		if i == 0 {
			min = val
			max = val
		} else {
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
		}
		i++
	}
	return min, max
}
