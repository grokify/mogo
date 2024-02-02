package strconvutil

import (
	"strconv"

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
func SliceItoa[S ~[]E, E constraints.Integer](s S, dedupe, sort bool) []string {
	var out []string
	for _, v := range s {
		out = append(out, strconv.Itoa(int(v)))
	}
	if dedupe {
		out = slicesutil.Dedupe(out)
	}
	if sort {
		slicesutil.Sort(out)
	}
	return out
}
