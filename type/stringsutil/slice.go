package stringsutil

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/grokify/mogo/type/number"
	"github.com/grokify/mogo/type/slicesutil"
	"golang.org/x/exp/slices"
)

type Strings []string

func (strs Strings) FilterIndexes(indexes []uint) (Strings, error) {
	n := Strings{}
	for _, idx := range indexes {
		if int(idx) >= len(strs) {
			return n, errors.New("index out of bounds")
		}
		n = append(n, strs[idx])
	}
	return n, nil
}

/*
// Unshift adds an element at the first position of the slice.
// EOL: use `slices.Insert()` instead.`
func Unshift(elems []string, x string) []string {
	return append([]string{x}, elems...)
}
*/

// SliceCondenseSpace trims space from lines and removes empty lines. `unique` dedupes lines and `sort`
// preforms a sort on the results.
func SliceCondenseSpace(elems []string, dedupeResults, sortResults bool) []string {
	results := SliceTrimSpace(elems, true)
	if dedupeResults {
		results = slicesutil.Dedupe(results)
	}
	if sortResults {
		sort.Strings(results)
	}
	return results
}

// SliceTrimSpace removes leading and trailing spaces per string. If condense
// is used, empty strings are removed.
func SliceTrimSpace(elems []string, condense bool) []string {
	var new []string
	for _, el := range elems {
		if el = strings.TrimSpace(el); el != "" || !condense {
			new = append(new, el)
		}
	}
	return new
}

// SliceTrim trims each line in a slice of lines using a provided cut string.
func SliceTrim(elems []string, cutstr string, condense bool) []string {
	var new []string
	for _, el := range elems {
		if el = strings.Trim(el, cutstr); el != "" || !condense {
			new = append(new, el)
		}
	}
	return new
}

/*
// JoinAny takes an array of `any`` and converts
// each value to a string using fmt.Sprintf("%v")
func JoinAny(a []any, sep string) string {
	strs := []string{}
	for _, item := range a {
		strs = append(strs, fmt.Sprintf("%v", item))
	}
	return strings.Join(strs, sep)
}

func JoinInt(a []int, sep string) string {
	strs := []string{}
	for _, i := range a {
		strs = append(strs, strconv.Itoa(i))
	}
	return strings.Join(strs, sep)
}

func JoinCondenseTrimSpace(slice []string, sep string) string {
	return strings.Join(SliceTrimSpace(slice, true), sep)
}
*/

func SliceCondenseRegexps(s []string, regexps []*regexp.Regexp, replacement string) []string {
	parts := []string{}
	for _, el := range s {
		for _, rx := range regexps {
			el = rx.ReplaceAllString(el, replacement)
		}
		el = strings.TrimSpace(el)
		if len(el) > 0 {
			parts = append(parts, el)
		}
	}
	return parts
}

func SliceCondensePunctuation(s []string) []string {
	parts := []string{}
	for _, part := range s {
		part = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(part, " ")
		part = regexp.MustCompile(`\s+`).ReplaceAllString(part, " ")
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			parts = append(parts, part)
		}
	}
	return parts
}

func SliceCondenseAndQuoteSpace(s []string, quoteLeft, quoteRight string) []string {
	return SliceCondenseAndQuote(s, " ", " ", quoteLeft, quoteRight)
}

func SliceCondenseAndQuote(s []string, trimLeft, trimRight, quoteLeft, quoteRight string) []string {
	newItems := []string{}
	for _, item := range s {
		item = strings.TrimLeft(item, trimLeft)
		item = strings.TrimRight(item, trimRight)
		if len(item) > 0 {
			item = quoteLeft + item + quoteRight
			newItems = append(newItems, item)
		}
	}
	return newItems
}

// SplitTrimSpace splits a string and trims spaces on remaining elements.
func SplitTrimSpace(s, sep string) []string {
	split := strings.Split(s, sep)
	strs := []string{}
	for _, str := range split {
		strs = append(strs, strings.TrimSpace(str))
	}
	return strs
}

var rxSplitLines = regexp.MustCompile(`(\r\n|\r|\n)`)

// SplitTextLines splits a string on the regxp `(\r\n|\r|\n)`.
func SplitTextLines(text string) []string {
	return rxSplitLines.Split(text, -1)
}

// SliceToSingleIntOrNeg converts a single element slice with a string to an integer or `-1`
func SliceToSingleIntOrNeg(s []string) int {
	if len(s) != 1 {
		return -1
	}
	num, err := strconv.Atoi(s[0])
	if err != nil {
		return -1
	}
	return num
}

func SliceChooseOnePreferredLowerTrimSpace(options, preferenceOrder []string) string {
	if len(options) == 0 {
		return ""
	} else if len(preferenceOrder) == 0 {
		return strings.ToLower(strings.TrimSpace(options[0]))
	}
	optMap := map[string]int{}
	for _, opt := range options {
		opt = strings.ToLower(strings.TrimSpace(opt))
		if len(opt) > 0 {
			optMap[opt] = 1
		}
	}
	for _, pref := range preferenceOrder {
		pref = strings.ToLower(strings.TrimSpace(pref))
		if _, ok := optMap[pref]; ok {
			return pref
		}
	}
	return strings.ToLower(strings.TrimSpace(options[0]))
}

// SlicesCompare returns 3 slices given 2 slices which represent intersection
// sets. The first set is present in slice A but not B, second for both and
// third present in slice B but not A.
func SlicesCompare(sliceA, sliceB []string) ([]string, []string, []string) {
	anob := []string{}
	both := []string{}
	bnoa := []string{}
	mapA := map[string]int{}
	for _, s := range sliceA {
		mapA[s] = 1
	}
	mapB := map[string]int{}
	for _, b := range sliceB {
		if _, ok := mapA[b]; ok {
			both = append(both, b)
		} else {
			bnoa = append(both, b)
		}
	}
	for _, a := range sliceA {
		if _, ok := mapB[a]; !ok {
			anob = append(anob, a)
		}
	}
	return SliceCondenseSpace(anob, true, true),
		SliceCondenseSpace(both, true, true),
		SliceCondenseSpace(bnoa, true, true)
}

/*
func SliceJoinQuoteMaxLength(slice []string, begQuote, endQuote, sep string, maxLength int) []string {
	words := []string{}
	curWords := []string{}
	curLength := 0
	for _, word := range slice {
		if maxLength > 0 && curLength+len(begQuote+word+endQuote+sep) > maxLength {
			words = append(words, strings.Join(curWords, sep))
			curWords = []string{}
			curLength = 0
		} else {
			curWords = append(curWords, begQuote+word+endQuote)
			curLength += len(begQuote + word + endQuote + sep)
		}
	}
	if len(curWords) > 0 {
		words = append(words, strings.Join(curWords, sep))
	}
	return words
}

type JoinCustomConfig struct {
	QuoteBegin string
	QuoteEnd   string
	Separator  string
	MaxLength  int
	TrimSpace  bool
	SkipEmpty  bool
}

func JoinCustom(slice []string, cfg JoinCustomConfig) []string {
	lines := []string{}
	curWords := []string{}
	curLength := 0
	for _, word := range slice {
		if cfg.TrimSpace {
			word = strings.TrimSpace(word)
		}
		if cfg.SkipEmpty && len(word) == 0 {
			continue
		}
		if cfg.MaxLength > 0 {
			if curLength+len(cfg.QuoteBegin+word+cfg.QuoteEnd+cfg.Separator) > cfg.MaxLength {
				lines = append(lines, strings.Join(curWords, cfg.Separator))
				curWords = []string{}
				curLength = 0
			} else {
				curWords = append(curWords, cfg.QuoteBegin+word+cfg.QuoteEnd)
				curLength += len(cfg.QuoteBegin + word + cfg.QuoteEnd + cfg.Separator)
			}
		} else {
			curWords = append(curWords, cfg.QuoteBegin+word+cfg.QuoteEnd)
		}
	}
	if len(curWords) > 0 {
		lines = append(lines, strings.Join(curWords, cfg.Separator))
	}
	return lines
}

func SliceJoinQuoteMaxLengthTrimSpaceSkipEmpty(slice []string, begQuote, endQuote, sep string, maxLength int) []string {
	words := []string{}
	curWords := []string{}
	curLength := 0
	for _, word := range slice {
		if maxLength > 0 && curLength+len(begQuote+word+endQuote+sep) > maxLength {
			words = append(words, strings.Join(curWords, sep))
			curWords = []string{}
			curLength = 0
		} else {
			curWords = append(curWords, begQuote+word+endQuote)
			curLength += len(begQuote + word + endQuote + sep)
		}
	}
	if len(curWords) > 0 {
		words = append(words, strings.Join(curWords, sep))
	}
	return words
}

func SliceJoinQuoted(slice []string, begQuote, endQuote, sep string) string {
	words := []string{}
	for _, word := range slice {
		words = append(words, begQuote+word+endQuote)
	}
	return strings.Join(words, sep)
}
*/

// SliceSubtract uses Set math to remove elements of filter from real.
func SliceSubtract(real, filter []string) []string {
	filtered := []string{}
	filterMap := map[string]int{}
	for _, f := range filter {
		filterMap[f] = 1
	}
	for _, r := range real {
		if _, ok := filterMap[r]; !ok {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// SliceToMap returns the slide where the slice elements are the keys of the map,
// and the value is the number of times it appears.
func SliceToMap(s []string) map[string]int {
	strmap := map[string]int{}
	for _, s := range s {
		strmap[s]++
	}
	return strmap
}

// SliceToDoc converts a slice to a map, trimming the values if desired. The `cfg` keys are
// the document property names or keys and the values are the index location of the slice.
func SliceToDoc(s []string, cfg map[string]int, trimSpace, inclEmpty bool) map[string]string {
	dat := map[string]string{}
	for k, idx := range cfg {
		if idx <= 0 || idx >= len(s) {
			continue
		}
		v := s[idx]
		if trimSpace {
			v = strings.TrimSpace(v)
		}
		if v == "" {
			continue
		}
		dat[k] = v
	}
	return dat
}

func SliceIntersection(list1, list2 []string) []string {
	map1 := map[string]int{}
	map2 := map[string]int{}
	for _, item1 := range list1 {
		map1[item1] = 1
	}
	for _, item2 := range list2 {
		if _, ok := map1[item2]; ok {
			map2[item2] = 1
		}
	}
	intersection := []string{}
	for key := range map2 {
		intersection = append(intersection, key)
	}
	return intersection
}

func SliceIntersectionCondenseSpace(slice1, slice2 []string) []string {
	return SliceIntersection(
		SliceCondenseSpace(slice1, true, false),
		SliceCondenseSpace(slice2, true, false))
}

// SliceIsEmpty checks to see if a slice is empty. If `skipEmptyStrings`
// it will also return empty if all elements are empty strings or
// only contain spaces.
func SliceIsEmpty(s []string, skipEmptyStrings bool) bool {
	if len(s) == 0 {
		return true
	}
	if !skipEmptyStrings {
		return false
	}
	for _, s := range s {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			return false
		}
	}
	return true
}

// SliceJoinFunc joins a slice passing each elemen through the supplied function `f`.
func SliceJoinFunc(s []string, sep string, f func(string) string) string {
	if f == nil {
		return strings.Join(s, sep)
	}
	var n []string
	for _, el := range s {
		n = append(n, f(el))
	}
	return strings.Join(n, sep)
}

// SliceOrderExplicit reoders the values of a slice using a requested input order
// where the output is ordered by the requested order, minus missing requests, and
// followed by non-ordered items. In addition to the output slide, an output slice
// of index locations is also provided.
func SliceOrderExplicit(s, order []string, inclUnordered bool) ([]string, []int) {
	if len(s) == 0 {
		return []string{}, []int{}
	} else if len(s) == 1 {
		return []string{s[0]}, []int{0}
	} else if len(order) == 0 {
		return slices.Clone(s), number.SliceIntBuildBeginEnd(0, len(s)-1)
	}
	smap := map[string]int{}
	for i, si := range s {
		smap[si] = i
	}
	var strs []string
	var idxs []int
	omap := map[string]int{}
	for i, ord := range order {
		omap[ord] = i
		if idx, ok := smap[ord]; ok {
			strs = append(strs, ord)
			idxs = append(idxs, idx)
		}
	}
	if inclUnordered {
		for si, sv := range s {
			if _, ok := omap[sv]; ok {
				continue
			}
			strs = append(strs, sv)
			idxs = append(idxs, si)
		}
	}
	if len(strs) != len(idxs) {
		panic("strs and idxs length mismatch")
	} else {
		return strs, idxs
	}
}

/*
// SliceDedupeOrdered is a function with similar functionality to SliceOrderExplicit().
func SliceDedupeOrdered(s, ordered []string, sortAsc bool, inclUnordered bool) []string {
	s = slicesutil.Dedupe(s)
	ordered = slicesutil.Dedupe(ordered)
	if len(s) == 0 {
		return []string{}
	} else if len(ordered) == 0 {
		if !inclUnordered {
			return []string{}
		} else {
			var out []string
			out = append(out, s...)
			if sortAsc {
				sort.Strings(out)
			}
			return out
		}
	}

	srcMap := map[string]int{}
	for _, s := range s {
		srcMap[s]++
	}
	var out []string
	outMap := map[string]int{}
	for _, ord := range ordered {
		if _, ok := srcMap[ord]; ok {
			out = append(out, ord)
			outMap[ord]++
		}
	}
	if inclUnordered {
		var unordered []string
		for _, si := range s {
			if _, ok := outMap[si]; !ok {
				unordered = append(unordered, si)
			}
		}
		if sortAsc {
			sort.Strings(unordered)
		}
		if len(unordered) > 0 {
			out = append(out, unordered...)
		}
	}

	return out
}
*/

// SliceSplitLengthStats returns a `map[int]int` indicating how many
// strings of which length are present.
func SliceSplitLengthStats(s []string, sep string) map[int]int {
	stats := map[int]int{}
	for _, s := range s {
		p := strings.Split(s, sep)
		stats[len(p)]++
	}
	return stats
}

// SliceBySplitLength returns lines by split length. This is useful for analyzing
// what types of data exist with different lengths.
func SliceBySplitLength(s []string, sep string) map[int][]string {
	bylen := map[int][]string{}
	for _, s := range s {
		p := strings.Split(s, sep)
		if _, ok := bylen[len(p)]; !ok {
			bylen[len(p)] = []string{}
		}
		bylen[len(p)] = append(bylen[len(p)], s)
	}
	return bylen
}
