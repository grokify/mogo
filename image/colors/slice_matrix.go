package colors

import (
	"fmt"
	"image/color"
	"sort"
)

func MatrixColumn(m [][]color.Color, colIdx uint) ([]color.Color, error) {
	colIdxInt := int(colIdx)
	clrs := []color.Color{}
	for i, row := range m {
		if colIdxInt >= len(row) {
			return clrs, fmt.Errorf("col [%d] not found on row [%d]", colIdx, i)
		}
		clrs = append(clrs, row[colIdxInt])
	}
	return clrs, nil
}

func MatrixUnique(c [][]color.Color) []color.Color {
	if len(c) == 0 {
		return []color.Color{}
	}
	uniques := map[string]color.Color{}
	keys := []string{}
	for _, row := range c {
		for _, cx := range row {
			cxs := ColorString(cx)
			if _, ok := uniques[cxs]; !ok {
				uniques[cxs] = cx
				keys = append(keys, cxs)
			}
		}
	}
	return uniqueColorsToSlice(keys, uniques)
}

func SliceUnique(c []color.Color) []color.Color {
	if len(c) == 0 {
		return c
	}
	uniques := map[string]color.Color{}
	keys := []string{}
	for _, cx := range c {
		cxs := ColorString(cx)
		if _, ok := uniques[cxs]; !ok {
			uniques[cxs] = cx
			keys = append(keys, cxs)
		}
	}
	return uniqueColorsToSlice(keys, uniques)
}

func uniqueColorsToSlice(keys []string, uniques map[string]color.Color) []color.Color {
	sort.Strings(keys)
	cslice := []color.Color{}
	for _, cxs := range keys {
		if c, ok := uniques[cxs]; ok {
			cslice = append(cslice, c)
		} else {
			panic(fmt.Sprintf("color [%s] not found", cxs))
		}
	}
	return cslice
}
