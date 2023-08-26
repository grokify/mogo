package maputil

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type MapStringString map[string]string

// Encode encodes the values into “URL encoded” form
// ("bar=baz&foo=quux") sorted by key.
func (m MapStringString) Encode() string {
	if m == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(v))
	}
	return buf.String()
}

func (m MapStringString) Get(key string) string {
	if v, ok := m[key]; ok {
		return v
	} else {
		return ""
	}
}

func (m MapStringString) Gets(inclNonMatches bool, keys []string) []string {
	var ret []string
	for _, k := range keys {
		if v, ok := m[k]; ok {
			ret = append(ret, v)
		} else if inclNonMatches {
			ret = append(ret, "")
		}
	}
	return ret
}

// Subset returns a subset of a `MapStringString`. `trimSpace` removes leading/trailing spaces on from the
// source values. `inclEmpty` includes keys where the value matches the empty string. `inclUnknown` adds
// desired keys in the resulting map which are not known in the source map.
func (m MapStringString) Subset(keys []string, trimSpace, inclEmpty, inclUnknown bool) MapStringString {
	newMap := map[string]string{}
	keyMap := map[string]int{}
	for i, k := range keys {
		keyMap[k] = i
	}
	for k, v := range m {
		if _, ok := keyMap[k]; ok {
			if trimSpace {
				v = strings.TrimSpace(v)
			}
			if v == "" && !inclEmpty {
				continue
			}
			newMap[k] = v
		} else if inclUnknown && !inclEmpty {
			newMap[k] = ""
		}
	}
	return newMap
}

func ParseMapStringString(s string) (MapStringString, error) {
	mss := make(MapStringString)
	err := parseQuery(mss, s)
	return mss, err
}

func parseQuery(mss MapStringString, s string) (err error) {
	for s != "" {
		var key string
		key, s, _ = strings.Cut(s, "&")
		if strings.Contains(key, ";") {
			err = fmt.Errorf("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		mss[key] = value
	}
	return err
}
