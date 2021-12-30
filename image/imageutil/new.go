package imageutil

import (
	"image"
	"image/color"
)

func NewRGBAColor(rect image.Rectangle, clr color.RGBA) *image.RGBA {
	img := image.NewRGBA(rect)
	PaintColorRGBA(img, clr)
	return img
}

func NewRGBATransparent(rect image.Rectangle) *image.RGBA {
	img := image.NewRGBA(rect)
	PaintColorRGBA(img, color.RGBA{255, 255, 255, 0})
	return img
}

func NewRGBAWhite(rect image.Rectangle) *image.RGBA {
	return NewRGBAColor(rect, color.RGBA{255, 255, 255, 255})
}
