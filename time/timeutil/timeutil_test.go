package timeutil

import (
	"testing"
	"time"
)

var dt6ForDt14Tests = []struct {
	v    int64
	want int32
}{
	{int64(20060102150405), int32(20060102)}}

func TestDt6ForDt14(t *testing.T) {
	for _, tt := range dt6ForDt14Tests {
		got := Dt6ForDt14(tt.v)
		if got != tt.want {
			t.Errorf("Dt6ForDt14(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}

var dt8ForStringTests = []struct {
	v    string
	want int32
}{
	{"2006-01-02T15:04:05Z", int32(20060102)}}

func TestDt8ForString(t *testing.T) {
	for _, tt := range dt8ForStringTests {
		got, err := Dt8ForString(time.RFC3339, tt.v)
		if err != nil {
			t.Errorf("Dt8ForString(%v): want %v, error ", tt.v, tt.want, err)
		}
		if got != tt.want {
			t.Errorf("Dt8ForString(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}
