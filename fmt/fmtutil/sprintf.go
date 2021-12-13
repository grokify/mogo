// Package fmtutil implements some formatting utility functions.
package fmtutil

import (
	"strconv"

	"github.com/grokify/mogo/type/number"
)

func SprintfFormatLeadingCharLength(char string, length uint) string {
	return "%" + char + strconv.Itoa(int(length)) + "d"
}

func SprintfFormatLeadingCharMaxIntVal(char string, value int) string {
	return SprintfFormatLeadingCharLength(char, number.IntLength(value))
}
