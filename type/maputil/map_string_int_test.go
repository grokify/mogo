package maputil

import (
	"testing"
)

var msiSortTests = []struct {
	data          map[string]int
	sortBy        string
	wantNameFirst string
}{
	{
		data: map[string]int{
			"ARecord": 2, "BRecord": 1, "CRecord": 4, "DRecord": 3},
		sortBy:        SortNameAsc,
		wantNameFirst: "ARecord"},
	{
		data: map[string]int{
			"ARecord": 2, "BRecord": 1, "CRecord": 4, "DRecord": 3},
		sortBy:        SortNameDesc,
		wantNameFirst: "DRecord"},
	{
		data: map[string]int{
			"ARecord": 2, "BRecord": 1, "CRecord": 4, "DRecord": 3},
		sortBy:        SortValueAsc,
		wantNameFirst: "BRecord"},
	{
		data: map[string]int{
			"ARecord": 2, "BRecord": 1, "CRecord": 4, "DRecord": 3},
		sortBy:        SortValueDesc,
		wantNameFirst: "CRecord"},
}

func TestMapStringIntSort(t *testing.T) {
	for _, tt := range msiSortTests {
		msi := MapStringInt(tt.data)
		sorted := msi.Sorted(tt.sortBy)
		if len(sorted) > 0 {
			gotNameFirst := sorted[0].Name
			if tt.wantNameFirst != gotNameFirst {
				t.Errorf("maputil.MapStringInt.Sorted() sort [%s] Error: want [%s], got [%s]",
					tt.sortBy, tt.wantNameFirst, gotNameFirst)
			}
		}
	}
}
