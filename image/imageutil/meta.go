package imageutil

import (
	"image"
	"os"
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

type ImageStats struct {
	Width  int
	Height int
}

func ImageStatsNil() ImageStats {
	return ImageStats{
		Width:  -1,
		Height: -1}
}

func ImageStatsRect(rect image.Rectangle) ImageStats {
	return ImageStats{
		Width:  rect.Max.X - rect.Min.X,
		Height: rect.Max.Y - rect.Min.Y}
}
