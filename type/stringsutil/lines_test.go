package stringsutil

import (
	"testing"
)

var toLineFeedsTests = []struct {
	inputString    string
	linefeedString string
}{
	{"hello\r\nworld\nagain", "hello\nworld\nagain"},
	{"hello\rworld\ragain", "hello\nworld\nagain"},
	{"hello\nworld\nagain", "hello\nworld\nagain"},
}

func TestToLineFeeds(t *testing.T) {
	for _, tt := range toLineFeedsTests {
		tryLineFeeds := ToLineFeeds(tt.inputString)
		if tryLineFeeds != tt.linefeedString {
			t.Errorf("stringsutil.ToLineFeeds(\"%s\") Error: want [%s], got [%s]",
				tt.inputString, tt.linefeedString, tryLineFeeds)
		}
	}
}
