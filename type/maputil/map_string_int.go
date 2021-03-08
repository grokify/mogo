package maputil

import (
	"sort"
	"strconv"
	"strings"
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

const (
	SortNameAsc   = "name asc"
	SortNameDesc  = "name desc"
	SortValueAsc  = "value asc"
	SortValueDesc = "value desc"
)

// Sorted returns a set of key names and values sorted by
// the sort type.
func (msi MapStringInt) Sorted(sortBy string) RecordSet {
	sortBy = strings.ToLower(strings.TrimSpace(sortBy))
	records := []Record{}
	for name, count := range msi {
		records = append(records, Record{Name: name, Value: count})
	}
	switch sortBy {
	case SortNameAsc:
		sort.Slice(records, func(i, j int) bool {
			return records[i].Name < records[j].Name
		})
	case SortNameDesc:
		sort.Slice(records, func(i, j int) bool {
			return records[i].Name > records[j].Name
		})
	case SortValueAsc:
		sort.Slice(records, func(i, j int) bool {
			return records[i].Value < records[j].Value
		})
	case SortValueDesc:
		sort.Slice(records, func(i, j int) bool {
			return records[i].Value > records[j].Value
		})
	}
	return records
}

type Record struct {
	Name  string
	Value int
}

type RecordSet []Record

func (rs RecordSet) Total() int {
	total := 0
	for _, rec := range rs {
		total += rec.Value
	}
	return total
}

func (rs RecordSet) Markdown(prefix, sep string, countFirst, addTotal bool) string {
	lines := []string{}
	for _, rec := range rs {
		if countFirst {
			lines = append(lines, prefix+strconv.Itoa(rec.Value)+sep+rec.Name)
		} else {
			lines = append(lines, prefix+rec.Name+sep+strconv.Itoa(rec.Value))
		}
	}
	if addTotal {
		total := rs.Total()
		if countFirst {
			lines = append(lines, prefix+strconv.Itoa(total)+sep+"Total")
		} else {
			lines = append(lines, prefix+"Total"+sep+strconv.Itoa(total))
		}
	}
	return strings.Join(lines, "\n")
}
