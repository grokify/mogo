package sortutil

import (
	"regexp"
	"sort"
	"strconv"
)

// parsePrefixNumberSafe parses a string into prefix and optional numeric suffix.
// If the suffix is not all digits, it's ignored.
func parsePrefixNumberSafe(s string) (prefix string, number *int) {
	// Match optional dash and trailing digits
	re := regexp.MustCompile(`^([A-Za-z0-9]+)-?([0-9]+)?$`)
	m := re.FindStringSubmatch(s)
	if m == nil || m[2] == "" {
		// No numeric suffix
		return s, nil
	}

	n, err := strconv.Atoi(m[2])
	if err != nil {
		// Treat as no numeric suffix
		return s, nil
	}

	return m[1], &n
}

// IntegerSuffix sorts strings like "ABC-10", "ABC-2", "XYZ", "FOO-bar".
// - Stable sort
// - Only numeric suffixes are considered
func IntegerSuffix(values []string) []string {
	type item struct {
		orig   string
		prefix string
		num    *int
	}

	items := make([]item, len(values))
	for i, v := range values {
		prefix, num := parsePrefixNumberSafe(v)
		items[i] = item{
			orig:   v,
			prefix: prefix,
			num:    num,
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].prefix == items[j].prefix {
			// Both have numeric suffix
			if items[i].num != nil && items[j].num != nil {
				return *items[i].num < *items[j].num
			}
			// Only one has numeric suffix
			if items[i].num != nil {
				return true
			}
			if items[j].num != nil {
				return false
			}
			// Neither has numeric suffix
			return items[i].orig < items[j].orig
		}
		return items[i].prefix < items[j].prefix
	})

	result := make([]string, len(values))
	for i, it := range items {
		result[i] = it.orig
	}

	return result
}
