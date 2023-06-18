package timeutil

import (
	"testing"
)

var dt6ForDT14Tests = []struct {
	v    int64
	want int32
}{
	{int64(20060102150405), int32(20060102)}}

func TestDT6ForDT14(t *testing.T) {
	for _, tt := range dt6ForDT14Tests {
		got := DT6ForDT14(tt.v)
		if got != tt.want {
			t.Errorf("Dt6ForDt14(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}

/*
var fromToTests = []struct {
	v    string
	want string
}{
	{"Wed, 25 May 2016 11:07:20 +0000", "2016-05-25T11:07:20Z"}}

func FromToTest(t *testing.T) {
	for _, tt := range fromToTests {
		got, err := FromTo(tt.v, time.RFC1123Z, time.RFC3339)
		if err != nil {
			t.Errorf("FromTo(%v): want %v, error ", tt.v, tt.want, err)
		}
		if got != tt.want {
			t.Errorf("FromTo(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}
*/
