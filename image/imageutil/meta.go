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

// ColorsMatrix returns colors for an image covering the pixels
// described in `image.Rectangle`: https://pkg.go.dev/image#Rectangle
func (meta *ImageMeta) ColorsMatrix() [][]color.Color {
	img := meta.Image
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

func (meta *ImageMeta) ColorsHistogram() map[string]uint {
	img := meta.Image
	cm := map[string]uint{}
	for xi := img.Bounds().Min.X; xi < img.Bounds().Max.X; xi++ {
		for yi := img.Bounds().Min.Y; yi < img.Bounds().Max.Y; yi++ {
			cm[colors.ColorString(img.At(xi, yi))]++
		}
	}
	return cm
}

type ImageMetadata struct {
	Width  int
	Height int
}

func NewImageMetadata(img image.Image) ImageMetadata {
	return ImageMetadata{
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy()}
}
