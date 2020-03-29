package codegen

import (
	"testing"
)

var goCodeNestedstructsToPointersTests = []struct {
	v    string
	want string
}{
	{"type FooBar struct {\nFoo Bar}\n", "type FooBar struct {\nFoo *Bar}\n"},
	{"type FooBar struct {\nFoo []Bar}\n", "type FooBar struct {\nFoo []*Bar}\n"},
	{"type FooBar struct {\nFoo map[string]Bar}\n", "type FooBar struct {\nFoo map[string]*Bar}\n"},
}

func TestGoCodeNestedstructsToPointers(t *testing.T) {
	for _, tt := range goCodeNestedstructsToPointersTests {
		got := GoCodeNestedstructsToPointers(tt.v)
		if got != tt.want {
			t.Errorf("GoCodeNestedstructsToPointers(\"%v\"): want [%v], got [%v]", tt.v, tt.want, got)
		}
	}
}
