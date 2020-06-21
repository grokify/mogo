package maputil

import (
	"sort"
)

type MapStringInt map[string]int

func (msi MapStringInt) AddKey(key string, val int) {
	if _, ok := msi[key]; !ok {
		msi[key] = 0
	}
	msi[key] += val
}

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

func (msi MapStringInt) KeyValue(key string, defaultValue int) int {
	if val, ok := msi[key]; ok {
		return val
	}
	return defaultValue
}

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
