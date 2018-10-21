package maputil

import (
	"reflect"
	"sort"
	"strings"

	"github.com/grokify/gotilla/sort/sortutil"
)

func StringKeys(mp interface{}) []string {
	keysVal := reflect.ValueOf(mp).MapKeys()
	keysArr := []string{}
	for _, key := range keysVal {
		keysArr = append(keysArr, key.String())
	}
	return keysArr
}

func StringKeysSorted(mp interface{}) []string {
	keysArr := StringKeys(mp)
	sort.Strings(keysArr)
	return keysArr
}

func StringKeysToLowerSorted(mp interface{}) []string {
	keysArr := StringKeys(mp)
	for i := range keysArr {
		keysArr[i] = strings.ToLower(keysArr[i])
	}
	sort.Strings(keysArr)
	return keysArr
}

func MapSSMerge(first map[string]string, second map[string]string) map[string]string {
	newMap := map[string]string{}
	for k1, v1 := range first {
		newMap[k1] = v1
	}
	for k2, v2 := range second {
		newMap[k2] = v2
	}
	return newMap
}

func MapSSValOrEmpty(data map[string]string, key string) string {
	if val, ok := data[key]; ok {
		return val
	}
	return ""
}

func MapSSEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}

	return true
}

type MapInt64Int64 map[int64]int64

func (m MapInt64Int64) KeysSorted() []int64 {
	keys := []int64{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(sortutil.Int64Slice(keys))
	return keys
}

func (m MapInt64Int64) ValuesSortedByKeys() []int64 {
	vals := []int64{}
	keys := m.KeysSorted()
	for _, k := range keys {
		vals = append(vals, k)
	}
	return vals
}
