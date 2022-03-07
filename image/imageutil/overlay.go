package imageutil

import (
	"image"
	"image/draw"
	"strings"
)

func Overlay(src, overlay image.Image, offset image.Point) image.Image {
	output := image.NewRGBA(src.Bounds())
	draw.Draw(output, src.Bounds(), src, image.Point{}, draw.Src)
	draw.Draw(output, overlay.Bounds().Add(offset), overlay, image.Point{}, draw.Src)
	return output
}

const (
	LocUpper      = "upper"
	LocMiddle     = "middle"
	LocLower      = "lower"
	LocLeft       = "left"
	LocCenter     = "center"
	LocRight      = "right"
	LocUpperLeft  = "upperleft"
	LocUpperRight = "upperright"
	LocLowerLeft  = "lowerleft"
	LocLowerRight = "lowerright"
)

func OverlayMore(src, overlay image.Image, overlayLocation string, padX, padY int) image.Image {
	return Overlay(src, overlay, OverlayOffset(src.Bounds(), overlay.Bounds(), overlayLocation, padX, padY))
}

func OverlayOffset(src, overlay image.Rectangle, overlayLocation string, padX, padY int) image.Point {
	pt := image.Point{}
	if strings.Contains(overlayLocation, LocUpper) {
		pt.Y = src.Min.Y + padY
	} else if strings.Contains(overlayLocation, LocLower) {
		pt.Y = src.Max.Y - overlay.Dy() - padY
	} else {
		pt.Y = src.Max.Y - ((src.Dy() - overlay.Dy()) / 2) + int(padY)
	}
	if strings.Contains(overlayLocation, LocLeft) {
		pt.X = src.Min.X + int(padX)
	} else if strings.Contains(overlayLocation, LocRight) {
		pt.X = src.Max.X - overlay.Dx() - padX
	} else {
		pt.X = src.Max.X - ((src.Dx() - overlay.Dx()) / 2) + padX
	}
	return pt
}

/*

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

*/
