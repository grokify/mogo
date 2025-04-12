package strslices

import (
	"strings"
	"testing"

	"github.com/grokify/mogo/strconv/strconvutil"
	"golang.org/x/exp/slices"
)

var sliceTests = []struct {
	v        []string
	want     []string
	condense bool
}{
	{[]string{"  foo  ", " ", "bar"}, []string{"foo", "bar"}, true},
	{[]string{"  foo  ", " ", "bar"}, []string{"foo", "", "bar"}, false},
}

func TestSlices(t *testing.T) {
	for _, tt := range sliceTests {
		got := Trim(tt.v, " ", tt.condense)
		if !slices.Equal(tt.want, got) {
			t.Errorf("stringsutil.SliceTrim(\"%s\") want (%s) got (%s)",
				strings.Join(tt.v, ","), strings.Join(tt.want, ","), strings.Join(got, ","))
		}
	}
}

var sliceOrderExplicitTests = []struct {
	v             []string
	ord           []string
	inclUnordered bool
	want          []string
	wantIdx       []int
}{
	{
		[]string{
			"Done",
			"In Progress",
			"In Review",
			"Open",
			"Second Review",
			"To Do",
		},
		[]string{
			"Open",
			"To Do",
			"In Progress",
			"In Review",
			"Second Review",
			"Done",
		},
		true,
		[]string{
			"Open",
			"To Do",
			"In Progress",
			"In Review",
			"Second Review",
			"Done",
		},
		[]int{3, 5, 1, 2, 4, 0},
	},
	{[]string{"foo", "bar"}, []string{"bar"}, false, []string{"bar"}, []int{1}},
	{[]string{"foo", "bar"}, []string{"bar"}, true, []string{"bar", "foo"}, []int{1, 0}},
	{[]string{"foo", "bar"}, []string{"bar", "foo"}, false, []string{"bar", "foo"}, []int{1, 0}},
	{[]string{"foo", "bar"}, []string{"bar", "foo"}, true, []string{"bar", "foo"}, []int{1, 0}},
	{[]string{"foo", "bar", "baz", "qux"}, []string{"qux", "foo", "bar"}, false, []string{"qux", "foo", "bar"}, []int{3, 0, 1}},
	{[]string{"foo", "bar", "baz", "qux"}, []string{"qux", "foo", "bar"}, true, []string{"qux", "foo", "bar", "baz"}, []int{3, 0, 1, 2}},
	{[]string{"foo", "bar", "baz", "qux"}, []string{"qux", "foo", "bar", "baz"}, true, []string{"qux", "foo", "bar", "baz"}, []int{3, 0, 1, 2}},
}

func TestSliceOrderExplicit(t *testing.T) {
	for i, tt := range sliceOrderExplicitTests {
		gotStr, gotIdx := OrderExplicit(tt.v, tt.ord, tt.inclUnordered)
		if !slices.Equal(tt.want, gotStr) {
			t.Errorf("stringsutil.SliceTrim(\"%s\", \"%s\", %v) strings want (%s) got (%s) test (%d)",
				strings.Join(tt.v, ","),
				strings.Join(tt.ord, ","),
				tt.inclUnordered,
				strings.Join(tt.want, ","),
				strings.Join(gotStr, ","), i)
		}
		if !slices.Equal(tt.wantIdx, gotIdx) {
			t.Errorf("stringsutil.SliceTrim(\"%s\", \"%s\", %v) indexes want (%s) got (%s) test (%d)",
				strings.Join(tt.v, ","),
				strings.Join(tt.ord, ","),
				tt.inclUnordered,
				strings.Join(strconvutil.SliceItoa(tt.wantIdx), ","),
				strings.Join(strconvutil.SliceItoa(gotIdx), ","), i,
			)
		}
	}
}
