package logutil

import (
	"testing"
)

var severityTests = []struct {
	execSeverity string
	itemSeverity string
	include      bool
	errorIsNil   bool
}{
	{SeverityDisabled, SeverityCritical, false, true},
	{SeverityError, SeverityDisabled, false, true},
	{SeverityError, SeverityWarning, false, true},
	{SeverityError, SeverityCritical, true, true},
	{SeverityDebug, SeverityWarning, true, true},
	{SeverityDebug, SeverityCritical, true, true},
	{SeverityCritical, SeverityError, false, true},
	{"foo", "bar", false, false},
}

func TestSeverity(t *testing.T) {
	for _, tt := range severityTests {
		tryIncl, err := SeverityInclude(tt.execSeverity, tt.itemSeverity)
		if err != nil && tt.errorIsNil {
			t.Errorf("logutil.SeverityInclude(\"%s\",\"%s\") error [%v]", tt.execSeverity, tt.itemSeverity, err.Error())
		}
		if tryIncl != tt.include {
			t.Errorf("logutil.SeverityInclude(\"%s\",\"%s\") want [%v] got [%v]", tt.execSeverity, tt.itemSeverity, tt.include, tryIncl)
		}
	}
}
