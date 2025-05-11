// Package fmtutil implements some formatting utility functions.
package fmtutil

import (
	"strconv"

	"github.com/grokify/mogo/type/number"
)

func SprintfFormatLeadingCharLength(char string, length int) string {
	if length < 0 {
		length = 0
	}
	return "%" + char + strconv.Itoa(length) + "d"
}

func SprintfFormatLeadingCharMaxIntVal(char string, value int) string {
	return SprintfFormatLeadingCharLength(char, number.IntLength(value))
}
