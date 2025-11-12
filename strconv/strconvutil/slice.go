package strconvutil

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/grokify/mogo/type/slicesutil"
	"golang.org/x/exp/constraints"
)

func SliceAtof(s []string, bitSize int) ([]float64, error) {
	var out []float64
	for _, si := range s {
		if fi, err := strconv.ParseFloat(si, bitSize); err != nil {
			return out, err
		} else {
			out = append(out, fi)
		}
	}
	return out, nil
}

// ParseInts parses all integers from a string, ignoring non-digit separators
func ParseInts(input string) ([]int, error) {
	// Match sequences of digits with optional leading minus sign
	re := regexp.MustCompile(`-?\d+`)
	matches := re.FindAllString(input, -1)

	result := make([]int, 0, len(matches))
	for _, m := range matches {
		num, err := strconv.Atoi(m)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %w", m, err)
		}
		result = append(result, num)
	}
	return result, nil
}

// SliceAtoi converts a slice of string integers.
func SliceAtoi(s []string, dedupe, sort bool) ([]int, error) {
	var out []int
	for _, si := range s {
		sv, err := strconv.Atoi(si)
		if err != nil {
			return []int{}, err
		}
		out = append(out, sv)
	}
	if dedupe {
		out = slicesutil.Dedupe(out)
	}
	if sort {
		slicesutil.Sort(out)
	}
	return out, nil
}

func SliceAtotFunc(funcFormat func(s string) (time.Time, error), s []string) ([]time.Time, error) {
	var times []time.Time
	if funcFormat == nil {
		return times, errors.New("funcFormat cannot be nil")
	}
	for _, si := range s {
		if ti, err := funcFormat(si); err != nil {
			return times, err
		} else {
			times = append(times, ti)
		}
	}
	return times, nil
}

func SliceAtou(s []string, dedupe, sort bool) ([]uint, error) {
	var out []uint
	for _, si := range s {
		sv, err := Atou(si)
		if err != nil {
			return []uint{}, err
		} else {
			out = append(out, sv)
		}
	}
	if dedupe {
		out = slicesutil.Dedupe(out)
	}
	if sort {
		slicesutil.Sort(out)
	}
	return out, nil
}

// SliceItoa converts a slice of `constraints.Integer` to a slice of `string`.
func SliceItoa[S ~[]E, E constraints.Integer](s S) []string {
	var out []string
	for _, v := range s {
		out = append(out, strconv.Itoa(int(v)))
	}
	return out
}

// SliceItoaMore converts a slice of `constraints.Integer` to a slice of `string` with additional
// functionality to dedupe and sort.
func SliceItoaMore[S ~[]E, E constraints.Integer](s S, dedupe, sort bool) []string {
	out := SliceItoa(s)
	if dedupe {
		out = slicesutil.Dedupe(out)
	}
	if sort {
		slicesutil.Sort(out)
	}
	return out
}

// JoinBytes joins a slice of strings and returns a byte array.
func JoinBytes(data []string, sep []byte) []byte {
	var out []byte
	for i, r := range data {
		out = append(out, []byte(r)...)
		if len(sep) > 0 && i < len(data)-1 {
			out = append(out, sep...)
		}
	}
	return out
}
