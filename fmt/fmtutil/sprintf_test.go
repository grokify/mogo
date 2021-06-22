package fmtutil

import (
	"fmt"
	"testing"
)

var sprintfFormatLeadingCharMaxIntValTests = []struct {
	prefix    string
	intval    int
	format    string
	formatInt int
	formatted string
}{
	{"0", 0, "%01d", 0, "0"},
	{"0", 9, "%01d", 1, "1"},
	{"0", 10, "%02d", 1, "01"},
	{"0", 14, "%02d", 1, "01"},
	{"0", 100, "%03d", 1, "001"},
	{"0", 100, "%03d", 11, "011"},
	{"0", 999, "%03d", 11, "011"},
	{"0", 1000, "%04d", 11, "0011"},
}

func TestSprintfFormatLeadingCharMaxIntVal(t *testing.T) {
	for _, tt := range sprintfFormatLeadingCharMaxIntValTests {
		format := SprintfFormatLeadingCharMaxIntVal(tt.prefix, tt.intval)
		if tt.format != format {
			t.Errorf("fmtutil.SprintfFormatLeadingCharMaxIntVal(\"%s\", %d) Error: want [%s], got [%s]",
				tt.prefix, tt.intval, tt.format, format)
		}
		formatted := fmt.Sprintf(format, tt.formatInt)
		if tt.formatted != formatted {
			t.Errorf("TestSprintfFormatLeadingCharMaxIntVal: Error: want [%s], got [%s]",
				tt.formatted, formatted)
		}
	}
}
