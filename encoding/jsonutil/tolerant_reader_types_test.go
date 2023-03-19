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
		if bool(obj.ValBool) != tt.valBool {
			t.Errorf("jsonutil.Bool Unmarshal(%v): want [%v] got [%v]", tt.plaintext, tt.valBool, obj.ValBool)
		}
		if int64(obj.ValInt64) != tt.valInt64 {
			t.Errorf("jsonutil.Int64 Unmarshal(%v): want [%v] got [%v]", tt.plaintext, tt.valInt64, obj.ValInt64)
		}
	}
}

var stringTests = []struct {
	v    string
	want String
	json string
}{
	{`{"value":"mystring"}`, "mystring", `{"value":"mystring"}`},
	{`{"value":1}`, "1", `{"value":"1"}`},
	{`{"value":false}`, "false", `{"value":"false"}`},
}

type field struct {
	Value String `json:"value"`
}

func TestString(t *testing.T) {
	for _, tt := range stringTests {
		f := &field{}
		err := json.Unmarshal([]byte(tt.v), f)
		if err != nil {
			t.Errorf("json.Unmarshal(%s): err (%s)", tt.v, err.Error())
			continue
		}
		if f.Value != tt.want {
			t.Errorf("json.Unmarshal(%s): mismatch want (%v), got (%v)", tt.v, tt.want, f.Value)
			continue
		}
		m, err := json.Marshal(f)
		if err != nil {
			panic(err) // should not happen
		}
		if string(m) != tt.json {
			t.Errorf("json.Marshal(%v): mismatch want (%v), got (%v)", f, tt.json, string(m))
		}
	}
}
