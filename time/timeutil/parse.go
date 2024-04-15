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
		"Jan": {"January"},
		"Feb": {"February"},
		"Mar": {"March"},
		"Apr": {"April"},
		"May": {""},
		"Jun": {"June"},
		"Jul": {"July"},
		"Aug": {"August"},
		"Sep": {"September", "Sept"},
		"Oct": {"October"},
		"Nov": {"November"},
		"Dec": {"December"}}
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
