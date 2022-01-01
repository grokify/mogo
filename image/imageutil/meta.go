package imageutil

import (
	"image"
	"image/color"
	"os"

	"github.com/grokify/mogo/image/colors"
)

func ImageAspect(img image.Image) float64 {
	return Aspect(img.Bounds().Dx(), img.Bounds().Dy())
}

func Aspect(width, height int) float64 {
	return float64(width) / float64(height)
}

func NegativeOffset(width, height, offset uint) image.Point {
	return image.Pt(int(width-offset), int(height-offset))
}

func IsNilOrEmpty(img image.Image) bool {
	if img == nil || img.Bounds().Dx() == 0 || img.Bounds().Dy() == 0 {
		return true
	}
	return false
}

type ImageMeta struct {
	File       *os.File
	FormatName string
	Image      image.Image
}

func (meta *ImageMeta) Stats() ImageStats {
	if meta.Image == nil {
		return ImageStatsNil()
	}
	return ImageStatsRect(meta.Image.Bounds())
}

func (meta *ImageMeta) Width() int {
	if meta.Image == nil {
		return -1
	}
	return meta.Image.Bounds().Max.X - meta.Image.Bounds().Min.X
}

func (meta *ImageMeta) Height() int {
	if meta.Image == nil {
		return -1
	}
	return meta.Image.Bounds().Max.Y - meta.Image.Bounds().Min.Y
}

func (meta *ImageMeta) ColorAverage() color.Color {
	return colors.ColorAverageImage(meta.Image)
}

type ImageStats struct {
	Width  int
	Height int
}

func ImageStatsNil() ImageStats {
	return ImageStats{
		Width:  -1,
		Height: -1}
}

func ImageStatsRect(r image.Rectangle) ImageStats {
	return ImageStats{
		Width:  r.Max.X - r.Min.X,
		Height: r.Max.Y - r.Min.Y}
}

func ImagePixelCount(img image.Image) int {
	m := ImageMeta{Image: img}
	w, h := m.Width(), m.Height()
	if w <= 0 || h <= 0 {
		return 0
	}
	return w * h
}

func ImageColors(img image.Image) [][]color.Color {
	rows := [][]color.Color{}
	minPt := img.Bounds().Min
	maxPt := img.Bounds().Max
	if minPt.X == maxPt.X || minPt.Y == maxPt.Y {
		return rows
	}
	for y := minPt.Y; y <= maxPt.Y; y++ {
		row := []color.Color{}
		for x := minPt.X; x <= maxPt.X; x++ {
			row = append(row, img.At(x, y))
		}
		rows = append(rows, row)
	}
	return rows
}
