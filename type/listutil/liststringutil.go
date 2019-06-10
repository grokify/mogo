package listutil

import (
	"strings"

	"github.com/grokify/gotilla/type/maputil"
)

func ListStringsToLowerUniqueSorted(list []string) []string {
	myMap := map[string]bool{}
	for _, myString := range list {
		myMap[myString] = true
	}
	listOut := maputil.StringKeysToLowerSorted(myMap)
	return listOut
}

func Include(haystack []string, needle string) bool {
	for _, try := range haystack {
		if try == needle {
			return true
		}
	}
	return false
}

func IncludeCaseInsensitive(haystack []string, needle string) bool {
	needleLower := strings.ToLower(needle)
	for _, try := range haystack {
		if strings.ToLower(try) == needleLower {
			return true
		}
	}
	return false
}

func StripSliceElements(in []string, strip []string) []string {
	out := []string{}
WORDS:
	for _, s := range in {
		for _, try := range strip {
			if s == try {
				continue WORDS
			}
		}
		out = append(out, s)
	}
	return out
}

func SplitCount(slice []string, size int) [][]string {
	slices := [][]string{}
	if size < 1 {
		return slices
	}
	current := []string{}
	for _, item := range slice {
		current = append(current, item)
		if len(current) == size {
			slices = append(slices, current)
			current = []string{}
		}
	}
	if len(current) > 0 {
		slices = append(slices, current)
	}
	return slices
}

func Unshift(a []string, x string) []string { return append([]string{x}, a...) }
