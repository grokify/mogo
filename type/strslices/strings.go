package strslices

import (
	"errors"
)

type Strings []string

func (strs Strings) FilterIndexes(indexes []int) (Strings, error) {
	n := Strings{}
	for _, idx := range indexes {
		if idx < 0 || idx >= len(strs) {
			return n, errors.New("index out of bounds")
		}
		n = append(n, strs[idx])
	}
	return n, nil
}
