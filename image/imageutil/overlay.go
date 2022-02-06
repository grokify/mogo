package imageutil

import (
	"image"
	"image/draw"
)

func Overlay(src, overlay image.Image, offset image.Point) image.Image {
	output := image.NewRGBA(src.Bounds())
	draw.Draw(output, src.Bounds(), src, image.Point{}, draw.Src)
	draw.Draw(output, overlay.Bounds().Add(offset), overlay, image.Point{}, draw.Src)
	return output
}

func OverlayCenterYLeftAlign(src, overlay image.Image) image.Image {
	h1 := src.Bounds().Dy()
	h2 := overlay.Bounds().Dy()
	offset := image.Pt(0, (h1-h2)/2)
	return Overlay(src, overlay, offset)
}

func OverlayLowerLeft(src, overlay image.Image) image.Image {
	return Overlay(
		src, overlay,
		image.Pt(
			src.Bounds().Min.X,
			src.Bounds().Max.Y-overlay.Bounds().Dy()))
}
