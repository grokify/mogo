package imageutil

import (
	"fmt"
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

// Crop takes an image and crops it to the specified rectangle. `CropImage`
// is from: https://stackoverflow.com/a/63256403.
func Crop(img image.Image, retain image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}
	return simg.SubImage(retain), nil
}

/*
// CropX crops an image by its width horizontally.
func CropX(img image.Image, width uint, align string) (image.Image, error) {
	if int(width) >= img.Bounds().Dx() {
		return img, nil
	}
	var xMin, xMax int
	switch strings.ToLower(strings.TrimSpace(align)) {
	case AlignLeft:
		xMin = img.Bounds().Min.X
		xMax = xMin + int(width)
	case AlignCenter:
		xMin = (img.Bounds().Max.Y - int(width)) / 2
		xMax = xMin + int(width)
	case AlignRight:
		xMax = img.Bounds().Max.X
		xMin = xMax - int(width)
	default:
		return nil, fmt.Errorf("alignment not supported [%s]", align)
	}
	return Crop(img, image.Rect(
		xMin,
		img.Bounds().Min.Y,
		xMax,
		img.Bounds().Max.Y))
}

// CropY crops an image by its height verticaly.
func CropY(img image.Image, height uint, align string) (image.Image, error) {
	if int(height) >= img.Bounds().Dy() {
		return img, nil
	}
	var yMin, yMax int
	switch strings.ToLower(strings.TrimSpace(align)) {
	case AlignTop:
		yMin = img.Bounds().Min.Y
		yMax = yMin + int(height)
	case AlignCenter:
		yMin = (img.Bounds().Max.Y - int(height)) / 2
		yMax = yMin + int(height)
	case AlignBottom:
		yMax = img.Bounds().Max.Y
		yMin = yMax - int(height)
	default:
		return nil, fmt.Errorf("alignment not supported [%s]", align)
	}
	return Crop(img, image.Rect(
		img.Bounds().Min.X,
		yMin,
		img.Bounds().Max.X,
		yMax))
}
*/

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

// CropX crops an image by its width horizontally.
func CropX(src image.Image, width uint, align string) image.Image {
	if int(width) > src.Bounds().Dx() {
		return src
	}
	var xMin int
	switch strings.ToLower(strings.TrimSpace(align)) {
	case AlignLeft:
		xMin = 0
	case AlignRight:
		xMin = -1 * int(width)
	default:
		xMin = -1 * (src.Bounds().Max.X - int(width)) / 2
	}
	offset := image.Point{
		X: xMin,
		Y: src.Bounds().Min.Y,
	}
	new := image.NewRGBA(image.Rect(0, 0, int(width), src.Bounds().Dy()))
	draw.Draw(new, src.Bounds().Add(offset), src, image.Point{}, draw.Over)
	return new
}

// CropY crops an image by its height verticaly.
func CropY(src image.Image, height uint, align string) image.Image {
	if int(height) > src.Bounds().Dy() {
		return src
	}
	var yMin int
	switch strings.ToLower(strings.TrimSpace(align)) {
	case AlignTop:
		yMin = 0
	case AlignBottom:
		yMin = -1 * int(height)
	default:
		yMin = -1 * (src.Bounds().Max.Y - int(height)) / 2
	}
	offset := image.Point{
		X: src.Bounds().Min.X,
		Y: yMin,
	}
	new := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx(), int(height)))
	draw.Draw(new, src.Bounds().Add(offset), src, image.Point{}, draw.Over)
	return new
}
