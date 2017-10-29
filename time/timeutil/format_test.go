package timeutil

import (
	"encoding/json"
	"testing"
)

var rfc3339YMDTimeTests = []struct {
	v    string
	want string
}{
	{`{"MyTime":"2001-02-03"}`, `{"MyTime":"2001-02-03"}`},
	{`{"MyTime":"0001-01-01"}`, `{"MyTime":"0001-01-01"}`},
	{`{"MyTime":""}`, `{"MyTime":"0001-01-01"}`},
	{`{}`, `{"MyTime":"0001-01-01"}`}}

type myStruct struct {
	MyTime RFC3339YMDTime
}

func TestRfc3339YMDTime(t *testing.T) {
	for _, tt := range rfc3339YMDTimeTests {
		my := myStruct{}
		//fmt.Println(tt.v)
		err := json.Unmarshal([]byte(tt.v), &my)
		if err != nil {
			t.Errorf("Rfc3339YMDTime.Unmarshal: with %v, want %v, err %v", tt.v, tt.want, err)
		}

		bytes, err := json.Marshal(my)
		if err != nil {
			t.Errorf("Rfc3339YMDTime.Marshal: with %v, want %v, err %v", tt.v, tt.want, err)
		}

		got := string(bytes)

		if got != tt.want {
			t.Errorf("Rfc3339YMDTime(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}
