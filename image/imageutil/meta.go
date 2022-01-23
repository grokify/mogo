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

func NewImageMeta(img image.Image) ImageMeta {
	return ImageMeta{Image: img}
}

func (meta *ImageMeta) ColorAverage() color.Color {
	return colors.ColorAverageImage(meta.Image)
}

func RectanglePixelCount(r image.Rectangle) int {
	w := r.Dx()
	h := r.Dy()
	if w <= 0 || h <= 0 {
		return 0
	}
	return w * h
}

// ImageColors returns colors for an image covering the pixels
// described in `image.Rectangle`: https://pkg.go.dev/image#Rectangle
func ImageColors(img image.Image) [][]color.Color {
	rows := [][]color.Color{}
	minPt := img.Bounds().Min
	maxPt := img.Bounds().Max
	if minPt.X == maxPt.X || minPt.Y == maxPt.Y {
		return rows
	}
	for y := minPt.Y; y < maxPt.Y; y++ {
		row := []color.Color{}
		for x := minPt.X; x < maxPt.X; x++ {
			row = append(row, img.At(x, y))
		}
		rows = append(rows, row)
	}
	return rows
}
