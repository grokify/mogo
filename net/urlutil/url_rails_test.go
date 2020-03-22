package urlutil

import (
	"encoding/json"
	"reflect"
	"testing"
)

type testObject struct {
	String      string
	Int         json.Number
	ArrayString []string
	ArrayInt    []string
	Map         map[string]string
}

var unmarshalRailsQSTests = []struct {
	v    string
	want testObject
}{
	{
		"string=foobar&map[foo]=bar&int=2&arrayInt[]=1&arrayInt[]=2&arrayString[]=foo&arrayString[]=bar",
		testObject{
			String:      "foobar",
			Int:         "2",
			ArrayInt:    []string{"1", "2"},
			ArrayString: []string{"foo", "bar"},
			Map:         map[string]string{"foo": "bar"},
		},
	},
}

func TestUnmarshalRailsQS(t *testing.T) {
	for _, tt := range unmarshalRailsQSTests {
		got := testObject{}
		err := UnmarshalRailsQS(tt.v, &got)
		if err != nil {
			t.Errorf("jsonutil.UnmarshalRailsQS(\"%s\") Error: [%s]",
				tt.v, err.Error())
		}
		if !reflect.DeepEqual(tt.want, got) {
			wantJSONBytes, err := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("jsonutil.UnmarshalRailsQS(\"%s\") Error: [%s]",
					tt.v, err.Error())
			}
			gotJSONBytes, err := json.Marshal(got)
			if err != nil {
				t.Errorf("jsonutil.UnmarshalRailsQS(\"%s\") Error: [%s]",
					tt.v, err.Error())
			}
			t.Errorf("jsonutil.UnmarshalRailsQS(\"%s\") want [%v], got [%v]",
				tt.v, string(wantJSONBytes), string(gotJSONBytes))
		}
	}
}
