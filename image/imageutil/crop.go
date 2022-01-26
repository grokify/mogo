package imageutil

import (
	"fmt"
	"image"
	"strings"
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

// CropX crops an image by its width horizontally.
func CropX(img image.Image, w uint, align string) (image.Image, error) {
	if int(w) > img.Bounds().Dx() {
		return img, nil
	}
	var xMin, xMax int
	switch strings.ToLower(strings.TrimSpace(align)) {
	case AlignLeft:
		xMin = img.Bounds().Min.X
		xMax = img.Bounds().Min.X + int(w)
	case AlignRight:
		xMin = img.Bounds().Max.X - int(w)
		xMax = img.Bounds().Max.X
	case AlignCenter:
		xMin = (img.Bounds().Max.Y - int(w)) / 2
		xMax = xMin + int(w)
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
func CropY(img image.Image, h uint, align string) (image.Image, error) {
	if int(h) > img.Bounds().Dy() {
		return img, nil
	}
	var yMin, yMax int
	switch strings.ToLower(strings.TrimSpace(align)) {
	case AlignTop:
		yMin = img.Bounds().Min.Y
		yMax = img.Bounds().Min.Y + int(h)
	case AlignBottom:
		yMin = img.Bounds().Max.Y - int(h)
		yMax = img.Bounds().Max.Y
	case AlignCenter:
		yMin = (img.Bounds().Max.Y - int(h)) / 2
		yMax = yMin + int(h)
	default:
		return nil, fmt.Errorf("alignment not supported [%s]", align)
	}
	return Crop(img, image.Rect(
		img.Bounds().Min.X,
		yMin,
		img.Bounds().Max.X,
		yMax))
}
