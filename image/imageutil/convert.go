package imageutil

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

func ImageToRGBA(img image.Image) *image.RGBA {
	if dst, ok := img.(*image.RGBA); ok {
		return dst
	}
	// https://stackoverflow.com/questions/31463756/convert-image-image-to-image-nrgba
	switch img := img.(type) {
	case *image.NRGBA:
		return NRGBAtoRGBA(img)
	case *image.Paletted:
		return ImageWithSetToRGBA(img)
	case *image.YCbCr:
		return YCbCrToRGBA(img)
	}
	// Use the image/draw package to convert to *image.RGBA.
	b := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), img, b.Min, draw.Src)
	return dst
}

type ImageWithSet interface {
	// ColorModel returns the Image's color model.
	ColorModel() color.Model
	// Bounds returns the domain for which At can return non-zero color.
	// The bounds do not necessarily contain the point (0, 0).
	Bounds() image.Rectangle
	// At returns the color of the pixel at (x, y).
	// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
	// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
	At(x, y int) color.Color
	Set(x, y int, c color.Color)
}

func NRGBAtoRGBA(imgNRGBA *image.NRGBA) *image.RGBA {
	rect := imgNRGBA.Bounds()
	imgRGBA := image.NewRGBA(rect)
	for x := rect.Min.X; x <= rect.Max.X; x++ {
		for y := rect.Min.Y; y <= rect.Max.Y; y++ {
			imgRGBA.Set(x, y, imgNRGBA.At(x, y))
		}
	}
	return imgRGBA
}

func ImageWithSetToRGBA(src ImageWithSet) *image.RGBA {
	rect := src.Bounds()
	imgRGBA := image.NewRGBA(rect)
	for x := rect.Min.X; x <= rect.Max.X; x++ {
		for y := rect.Min.Y; y <= rect.Max.Y; y++ {
			imgRGBA.Set(x, y, src.At(x, y))
		}
	}
	return imgRGBA
}

func ImageAnyToRGBA(src image.Image) *image.RGBA {
	// https://stackoverflow.com/questions/31463756/convert-image-image-to-image-nrgba
	b := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)
	return img
}

func YCbCrToRGBA(src *image.YCbCr) *image.RGBA {
	// https://stackoverflow.com/questions/31463756/convert-image-image-to-image-nrgba
	b := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)
	return img
}
