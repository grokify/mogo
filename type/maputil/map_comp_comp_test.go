package maputil

import (
	"testing"
)

func TestEqual(t *testing.T) {
	var equalTests = []struct {
		vA1       map[string]string
		vA2       map[string]string
		wantEqual bool
	}{
		{map[string]string{"foo": "bar", "baz": "qux"}, map[string]string{"baz": "qux", "foo": "bar"}, true},
		{map[string]string{"foo": "bar", "baz": "qux1"}, map[string]string{"baz": "qux", "foo": "bar"}, false},
	}

	for _, tt := range equalTests {
		m1 := MapCompComp[string, string](tt.vA1)
		try := m1.Equal(tt.vA2)
		if try != tt.wantEqual {
			t.Errorf("maputil.MapCompComp(%v).Equal(%v) Mismatch: want (%v), got (%v)",
				tt.vA1, tt.vA2, tt.wantEqual, try)
		}
	}
}

func TestFilterMergeViaMap(t *testing.T) {
	var equalTests = []struct {
		vA1          map[string]string
		vA2          map[string]string
		useNOnlyVal  bool
		nOnlyDefault *string
		want         map[string]string
	}{
		{
			map[string]string{"alpha": "bravo", "charlie": "delta", "golf": "hotel"},
			map[string]string{"alpha": "bravo", "echo": "foxtrot", "golf": "foo"},
			true, nil,
			map[string]string{"alpha": "bravo", "echo": "foxtrot", "golf": "hotel"},
		},
	}

	for _, tt := range equalTests {
		m1 := MapCompComp[string, string](tt.vA1)
		try := m1.FilterMergeByMap(tt.vA2, tt.useNOnlyVal, tt.nOnlyDefault)
		tryMap := MapCompComp[string, string](try)
		if !tryMap.Equal(tt.want) {
			t.Errorf("maputil.MapCompComp(%v).FilterMergeViaMap(%v) Mismatch: want (%v), got (%v)",
				tt.vA1, tt.vA2, tt.want, try)
		}
	}
}

func TestValueKeyCounts(t *testing.T) {
	var valueKeyCountsTests = []struct {
		v    map[string]string
		want map[string]int
	}{
		{map[string]string{"foo": "bar", "baz": "bar"}, map[string]int{"bar": 2}},
	}

	for _, tt := range valueKeyCountsTests {
		m := MapCompComp[string, string](tt.v)
		try := m.ValueKeyCounts()
		tryComp := MapCompComp[string, int](try)
		if !tryComp.Equal(tt.want) {
			t.Errorf("maputil.MapCompComp(%v).ValueKeyCounts(): Mismatch: want (%v), got (%v)",
				tt.v, tt.want, try)
		}
	}
}
