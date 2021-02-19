package jsonutil

import (
	"encoding/json"
	"testing"
)

type testObject struct {
	ValBool  Bool
	ValInt64 Int64
}

var tolerantReaderTests = []struct {
	plaintext string
	valBool   bool
	valInt64  int64
}{
	{`{"ValBool":" 1 ","ValInt64":"    "}`, true, 0},
	{`{"ValBool":"true","ValInt64":"    "}`, true, 0},
	{`{"ValBool":"true","ValInt64":"  314  "}`, true, 314},
	{`{"ValBool":1,"ValInt64":314}`, true, 314}}

func TestTolerantReader(t *testing.T) {
	for _, tt := range tolerantReaderTests {
		raw := []byte(tt.plaintext)
		obj := testObject{}
		err := json.Unmarshal(raw, &obj)
		if err != nil {
			t.Errorf("jsonutil.Bool Unmarshal(%v): err [%v]", tt.plaintext, err.Error())
		}
		if obj.ValBool.Value() != tt.valBool {
			t.Errorf("jsonutil.Bool Unmarshal(%v): want [%v] got [%v]", tt.plaintext, tt.valBool, obj.ValBool)
		}
		if obj.ValInt64.Value() != tt.valInt64 {
			t.Errorf("jsonutil.Int64 Unmarshal(%v): want [%v] got [%v]", tt.plaintext, tt.valInt64, obj.ValInt64)
		}
	}
}
