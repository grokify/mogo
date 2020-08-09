package imageutil

import (
	"image"

	"golang.org/x/image/draw"
)

func ImageToRGBA(img image.Image) *image.RGBA {
	// https://stackoverflow.com/questions/31463756/convert-image-image-to-image-nrgba
	switch img := img.(type) {
	case *image.NRGBA:
		return NRGBAtoRGBA(img)
	case *image.RGBA:
		return img
	case *image.YCbCr:
		return YCbCrToRGBA(img)
	}
	panic("Format not supported.")
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

func YCbCrToRGBA(src *image.YCbCr) *image.RGBA {
	// https://stackoverflow.com/questions/31463756/convert-image-image-to-image-nrgba
	b := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)
	return img
}
