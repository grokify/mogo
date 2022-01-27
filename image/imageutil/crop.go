package imageutil

import (
	"image"
	"strings"

	"golang.org/x/image/draw"
)

const (
	AlignTop    = "top"
	AlignCenter = "center"
	AlignBottom = "bottom"
	AlignLeft   = "left"
	AlignRight  = "right"
)

// Crop takes an image and crops it to the specified rectangle.
func Crop(src image.Image, retain image.Rectangle) image.Image {
	new := image.NewRGBA(retain)
	draw.Draw(new, new.Bounds(), src, retain.Min, draw.Over)
	return new
}

// CropX crops an image by its width horizontally.
func CropX(src image.Image, width uint, align string) image.Image {
	if int(width) >= src.Bounds().Dx() {
		return src
	}
	var xMin int
	switch strings.ToLower(strings.TrimSpace(align)) {
	case AlignLeft:
		xMin = src.Bounds().Min.X
	case AlignRight:
		xMin = src.Bounds().Max.X - int(width)
	default:
		xMin = (src.Bounds().Max.X - int(width)) / 2
	}
	return Crop(src, image.Rect(
		xMin,
		src.Bounds().Min.Y,
		xMin+int(width),
		src.Bounds().Max.Y))
}

// CropY crops an image by its height verticaly.
func CropY(src image.Image, height uint, align string) image.Image {
	if int(height) >= src.Bounds().Dy() {
		return src
	}
	var yMin int
	switch strings.ToLower(strings.TrimSpace(align)) {
	case AlignTop:
		yMin = src.Bounds().Min.Y
	case AlignBottom:
		yMin = src.Bounds().Max.Y - int(height)
	default:
		yMin = (src.Bounds().Max.Y - int(height)) / 2
	}
	return Crop(src, image.Rect(
		src.Bounds().Min.X,
		yMin,
		src.Bounds().Max.X,
		yMin+int(height)))
}

func Square(src image.Image) image.Image {
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()
	if width == height {
		return src
	} else if height > width {
		return CropY(src, uint(width), AlignCenter)
	}
	return CropX(src, uint(height), AlignCenter)
}
