package imageutil

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
)

func ImageToRGBA(src image.Image) *image.RGBA {
	/*
		// https://stackoverflow.com/questions/31463756/convert-image-image-to-image-nrgba
		switch img := img.(type) {
		case *image.NRGBA:
			return NRGBAtoRGBA(img)
		case *image.Paletted:
			return ImageWithSetToRGBA(img)
		case *image.YCbCr:
			return YCbCrToRGBA(img)
		}
	*/
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Rect, src, src.Bounds().Min, draw.Src)
	return dst
}

func ImageToPalettedPlan9(src image.Image) *image.Paletted {
	return ImageToPaletted(src, palette.Plan9)
}

func ImageToPaletted(src image.Image, p color.Palette) *image.Paletted {
	if dst, ok := src.(*image.Paletted); ok {
		return dst
	}
	dst := image.NewPaletted(src.Bounds(), p)
	draw.Draw(dst, dst.Rect, src, src.Bounds().Min, draw.Over)
	return dst
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

func ImageWithSetToRGBA(src draw.Image) *image.RGBA {
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
