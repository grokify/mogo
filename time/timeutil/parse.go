package timeutil

import (
	"strings"
	"time"
)

func ParseTimeCanonical(layout, value string) (time.Time, error) {
	return time.Parse(layout, ReplaceMonthCanonical(value))
}

func ParseTimeCanonicalFunc(layout string) func(s string) (time.Time, error) {
	return func(s string) (time.Time, error) {
		return ParseTimeCanonical(layout, s)
	}
}

func CanonicalMonthMap() map[string][]string {
	return map[string][]string{
		"Jan": []string{"January"},
		"Feb": []string{"February"},
		"Mar": []string{"March"},
		"Apr": []string{"April"},
		"May": []string{""},
		"Jun": []string{"June"},
		"Jul": []string{"July"},
		"Aug": []string{"August"},
		"Sep": []string{"September", "Sept"},
		"Oct": []string{"October"},
		"Nov": []string{"November"},
		"Dec": []string{"December"}}
}

func ReplaceMonthCanonical(s string) string {
	mm := CanonicalMonthMap()
	for k, vals := range mm {
		for _, v := range vals {
			if v == "" {
				continue
			}
			s = strings.ReplaceAll(s, v, k)
			s = strings.ReplaceAll(s, strings.ToLower(v), k)
		}
	}
	return s
}
