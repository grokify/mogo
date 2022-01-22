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

// CropVertical takes an image and crops it verticaly.
func CropVertical(img image.Image, h uint, align string) (image.Image, error) {
	align = strings.ToLower(strings.TrimSpace(align))
	if align != AlignTop && align != AlignBottom {
		align = AlignCenter
	}
	hMin := 0
	hMax := int(h)
	if align == AlignCenter {
		hMin = (img.Bounds().Max.Y - int(h)) / 2
		hMax = hMin + int(h)
	} else if align == AlignBottom {
		hMin = img.Bounds().Max.Y - int(h)
		hMax = img.Bounds().Max.Y
	}
	return Crop(img, image.Rect(img.Bounds().Min.X, hMin, img.Bounds().Max.X, hMax))
}

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
