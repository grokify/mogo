package imageutil

import (
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/grokify/mogo/image/padding"
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
	out := image.NewRGBA(retain)
	draw.Draw(out, out.Bounds(), src, retain.Min, draw.Over)
	return out
}

// CropX crops an image by its width horizontally.
func CropX(src image.Image, width int, align string) image.Image {
	if width < 0 || width >= src.Bounds().Dx() {
		return src
	}
	var xMin int
	switch strings.ToLower(strings.TrimSpace(align)) {
	case AlignLeft:
		xMin = src.Bounds().Min.X
	case AlignRight:
		xMin = src.Bounds().Max.X - width
	default:
		xMin = (src.Bounds().Max.X - width) / 2
	}
	return Crop(src, image.Rect(
		xMin,
		src.Bounds().Min.Y,
		xMin+width,
		src.Bounds().Max.Y))
}

// CropY crops an image by its height vertically.
func CropY(src image.Image, height int, align string) image.Image {
	if height < 0 || height >= src.Bounds().Dy() {
		return src
	}
	var yMin int
	switch strings.ToLower(strings.TrimSpace(align)) {
	case AlignTop:
		yMin = src.Bounds().Min.Y
	case AlignBottom:
		yMin = src.Bounds().Max.Y - height
	default:
		yMin = (src.Bounds().Max.Y - height) / 2
	}
	return Crop(src, image.Rect(
		src.Bounds().Min.X,
		yMin,
		src.Bounds().Max.X,
		yMin+height))
}

func CropPadding(src image.Image, isPadding padding.IsPaddingFunc) image.Image {
	if src == nil {
		return nil
	} else if isPadding == nil {
		return src
	} else {
		return Crop(src, padding.NonPaddingRectangle(src, isPadding))
	}
}

// SquareLarger returns a square image that is cropped to where the height and weight are equal
// to the larger of the source image. Additional padding is added, if necessary.
func (im Image) SquareLarger(bgcolor color.Color) image.Image { return squareLarger(im.Image, bgcolor) }

func squareLarger(src image.Image, bgcolor color.Color) image.Image {
	if src == nil {
		return nil
	}
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()
	switch {
	case width > height:
		new := AddBackgroundColor(image.NewRGBA(image.Rect(0, 0, width, width)), bgcolor)
		draw.Draw(new, new.Bounds(), src, image.Point{
			Y: src.Bounds().Min.Y + ((height - width) / 2),
			X: src.Bounds().Min.X}, draw.Over)
		return new
	case width < height:
		new := AddBackgroundColor(image.NewRGBA(image.Rect(0, 0, height, height)), bgcolor)
		draw.Draw(new, new.Bounds(), src, image.Point{
			X: src.Bounds().Min.X + ((width - height) / 2),
			Y: src.Bounds().Min.Y}, draw.Over)
		return new
	default:
		return src
	}
}

// SquareSmaller returns a square image that is cropped to where the height and weight
// are equal to the smaller of the source image. Source image pixes may be cropped and
// no additional pixels are added.
func (im Image) SquareSmaller() image.Image { return squareSmaller(im.Image) }

func squareSmaller(src image.Image) image.Image {
	if src == nil {
		return nil
	}
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()
	switch {
	case width > height:
		return CropX(src, height, AlignCenter)
	case width < height:
		return CropY(src, width, AlignCenter)
	default:
		return src
	}
}
