package imageutil

import (
	"image"
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
