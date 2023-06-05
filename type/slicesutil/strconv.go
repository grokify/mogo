package slicesutil

import (
	"github.com/grokify/mogo/strconv/strconvutil"
	"golang.org/x/exp/constraints"
)

func Itoa[S ~[]E, E constraints.Integer](s S) []string {
	out := []string{}
	for _, e := range s {
		out = append(out, strconvutil.Itoa(e))
	}
	return out
}
