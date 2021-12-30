package imageutil

import (
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
