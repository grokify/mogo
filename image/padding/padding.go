package padding

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/grokify/mogo/image/colors"
)

type IsPaddingFunc func(c color.Color) bool

// CreateIsPaddingFuncSimple returns a function that fulfills the
// `IsPaddingFunc` inteface.
func CreateIsPaddingFuncSimple(paddingColor color.Color) func(testColor color.Color) bool {
	return func(testColor color.Color) bool {
		return colors.Equal(paddingColor, testColor)
	}
}

// IsPaddingFuncWhite returns a function that fulfills the `IsPaddingFunc`
// interface.
func IsPaddingFuncWhite() func(testColor color.Color) bool {
	return CreateIsPaddingFuncSimple(color.White)
}

func AddPaddingUniform(im image.Image, paddingWidth uint, paddingColor color.Color) image.Image {
	out := image.NewRGBA(image.Rect(0, 0, im.Bounds().Dx()+int(2*paddingWidth), im.Bounds().Dy()+int(2*paddingWidth)))
	draw.Draw(out, out.Bounds(), &image.Uniform{paddingColor}, image.Point{}, draw.Src)
	draw.Draw(out, image.Rect(
		int(paddingWidth),
		int(paddingWidth),
		out.Bounds().Max.X-1-int(paddingWidth),
		out.Bounds().Max.Y-1-int(paddingWidth)), im, image.Point{}, draw.Over)
	return out
}

func NonPaddingRectangle(im image.Image, isPadding IsPaddingFunc) image.Rectangle {
	if isPadding == nil {
		isPadding = IsPaddingFuncWhite()
	}
	topP, rightP, bottomP, leftP := PaddingWidths(im, isPadding)
	return image.Rect(
		int(leftP)+im.Bounds().Min.X,
		int(topP)+im.Bounds().Min.Y,
		int(rightP)*-1+im.Bounds().Max.X,
		int(bottomP)*-1+im.Bounds().Max.Y)
}

func PaddingWidths(im image.Image, isPadding IsPaddingFunc) (top, right, bottom, left int) {
	top = PaddingWidthTop(im, isPadding)
	right = PaddingWidthRight(im, isPadding)
	bottom = PaddingWidthBottom(im, isPadding)
	left = PaddingWidthLeft(im, isPadding)
	return
}

func PaddingWidthTop(im image.Image, isPadding IsPaddingFunc) int {
	if im == nil {
		return 0
	}
	for yi := im.Bounds().Min.Y; yi < im.Bounds().Max.Y; yi++ {
		for xi := im.Bounds().Min.X; xi < im.Bounds().Max.X; xi++ {
			if !isPadding(im.At(xi, yi)) {
				if out := yi - im.Bounds().Min.Y; out < 0 {
					panic("padding width cannot be less than zero")
				} else {
					return out
				}
			}
		}
	}
	return im.Bounds().Dy()
}

func PaddingWidthBottom(im image.Image, isPadding IsPaddingFunc) int {
	if im == nil {
		return 0
	}
	for yi := im.Bounds().Max.Y - 1; yi >= im.Bounds().Min.Y; yi-- {
		for xi := im.Bounds().Min.X; xi < im.Bounds().Max.X; xi++ {
			if !isPadding(im.At(xi, yi)) {
				if out := im.Bounds().Max.Y - 1 - yi; out < 0 {
					panic("padding width cannot be less than zero")
				} else {
					return out
				}
			}
		}
	}
	return im.Bounds().Dy()
}

func PaddingWidthLeft(im image.Image, isPadding IsPaddingFunc) int {
	if im == nil {
		return 0
	}
	for xi := im.Bounds().Min.X; xi < im.Bounds().Max.X; xi++ {
		for yi := im.Bounds().Min.Y; yi < im.Bounds().Max.Y; yi++ {
			if !isPadding(im.At(xi, yi)) {
				if out := xi - im.Bounds().Min.X; out < 0 {
					panic("padding width cannot be less than zero")
				} else {
					return out
				}
			}
		}
	}
	return im.Bounds().Dx()
}

func PaddingWidthRight(im image.Image, isPadding func(c color.Color) bool) int {
	if im == nil {
		return 0
	}
	for xi := im.Bounds().Max.X - 1; xi >= im.Bounds().Min.X; xi-- {
		for yi := im.Bounds().Min.Y; yi < im.Bounds().Max.Y; yi++ {
			if !isPadding(im.At(xi, yi)) {
				if out := im.Bounds().Max.X - 1 - xi; out < 0 {
					panic("padding width cannot be less than zero")
				} else {
					return out
				}
			}
		}
	}
	return im.Bounds().Dx()
}
