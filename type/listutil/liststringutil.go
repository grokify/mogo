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
