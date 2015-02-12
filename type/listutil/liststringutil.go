package listutil

import (
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
