package imageutil

import (
	"fmt"
	"image"
	"image/color"

	"github.com/grokify/mogo/image/colors"
)

func RowsFilteredColor(img image.Image, c color.Color, cmore ...color.Color) []int {
	minPt := img.Bounds().Min
	maxPt := img.Bounds().Max
	rows := []int{}

	cmore = append([]color.Color{c}, cmore...)
	for y := minPt.Y; y <= maxPt.Y; y++ {
		for x := minPt.X; x <= maxPt.X; x++ {
			for _, ci := range cmore {
				if colors.ColorToHex(img.At(x, y)) == colors.ColorToHex(ci) {
					rows = append(rows, y)
					continue
				}
			}
		}
	}
	return rows
}

func ColsFilteredColor(img image.Image, c ...color.Color) []int {
	cols := []int{}
	if len(c) == 0 {
		return cols
	}
	clrs := colors.Colors(c)

	minPt := img.Bounds().Min
	maxPt := img.Bounds().Max
	fmt.Printf("X MIN [%d] MAX [%d]\n", minPt.X, maxPt.X)
COL:
	for x := minPt.X; x <= maxPt.X; x++ {
		for y := minPt.Y; y <= maxPt.Y; y++ {
			cx := img.At(x, y)
			if !clrs.In(cx) {
				continue COL
			}
		}
		cols = append(cols, x)
	}
	return cols
}
