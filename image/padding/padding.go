package padding

import (
	"fmt"
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

// CreateIsPaddingFuncMagicWand creates a padding func that matches a color with a `tolerance`
// for matching nearby colors.
func CreateIsPaddingFuncMagicWand(paddingColor color.Color, tolerance uint32) func(testColor color.Color) bool {
	return func(testColor color.Color) bool {
		return colors.MagicWandMatch(testColor, paddingColor, tolerance)
	}
}

// IsPaddingFuncWhite returns a function that fulfills the `IsPaddingFunc`
// interface.
func IsPaddingFuncWhite() func(testColor color.Color) bool {
	return CreateIsPaddingFuncSimple(color.White)
}

func IsPaddingFuncNearWhite() func(testColor color.Color) bool {
	return func(testColor color.Color) bool {
		return colors.IsNearWhite(testColor)
	}
}

func AddPaddingUniform(im image.Image, paddingWidth uint32, paddingColor color.Color) image.Image {
	out := image.NewRGBA(image.Rect(0, 0, im.Bounds().Dx()+int(2*paddingWidth), im.Bounds().Dy()+int(2*paddingWidth)))
	draw.Draw(out, out.Bounds(), &image.Uniform{paddingColor}, image.Point{}, draw.Src)
	draw.Draw(out, image.Rect(
		int(paddingWidth),
		int(paddingWidth),
		out.Bounds().Max.X-1-int(paddingWidth),
		out.Bounds().Max.Y-1-int(paddingWidth)), im, image.Point{}, draw.Over)
	return out
}

func NonPaddingRectangle(im image.Image, isPadding IsPaddingFunc, retainPaddingWidth uint32) image.Rectangle {
	if isPadding == nil {
		isPadding = IsPaddingFuncWhite()
	}
	topP, rhtP, botP, lftP := PaddingWidths(im, isPadding, retainPaddingWidth)
	m := map[string]int{
		"top":   topP,
		"right": rhtP,
		"bot":   botP,
		"left":  lftP,
	}
	return image.Rect(
		lftP+im.Bounds().Min.X,
		topP+im.Bounds().Min.Y,
		rhtP*-1+im.Bounds().Max.X,
		botP*-1+im.Bounds().Max.Y)
}

func PaddingWidths(im image.Image, isPadding IsPaddingFunc, retainPaddingWidth uint32) (top, right, bottom, left int) {
	top = PaddingWidthTop(im, isPadding, retainPaddingWidth)
	right = PaddingWidthRight(im, isPadding, retainPaddingWidth)
	bottom = PaddingWidthBottom(im, isPadding, retainPaddingWidth)
	left = PaddingWidthLeft(im, isPadding, retainPaddingWidth)
	return
}

func PaddingWidthTop(im image.Image, isPadding IsPaddingFunc, retainPaddingWidth uint32) int {
	if im == nil {
		return 0
	}
	for yi := im.Bounds().Min.Y; yi < im.Bounds().Max.Y; yi++ {
		for xi := im.Bounds().Min.X; xi < im.Bounds().Max.X; xi++ {
			if !isPadding(im.At(xi, yi)) {
				if out := yi - im.Bounds().Min.Y; out < 0 {
					panic("padding width cannot be less than zero")
				} else {
					return paddingWidthReduce(out, int(retainPaddingWidth))
				}
			}
		}
	}
	return im.Bounds().Dy()
}

func PaddingWidthBottom(im image.Image, isPadding IsPaddingFunc, retainPaddingWidth uint32) int {
	if im == nil {
		return 0
	}
	for yi := im.Bounds().Max.Y - 1; yi >= im.Bounds().Min.Y; yi-- {
		for xi := im.Bounds().Min.X; xi < im.Bounds().Max.X; xi++ {
			if !isPadding(im.At(xi, yi)) {
				if out := im.Bounds().Max.Y - 1 - yi; out < 0 {
					panic("padding width cannot be less than zero")
				} else {
					return paddingWidthReduce(out, int(retainPaddingWidth))
				}
			}
		}
	}
	return im.Bounds().Dy()
}

func PaddingWidthLeft(im image.Image, isPadding IsPaddingFunc, retainPaddingWidth uint32) int {
	if im == nil {
		return 0
	}
	for xi := im.Bounds().Min.X; xi < im.Bounds().Max.X; xi++ {
		for yi := im.Bounds().Min.Y; yi < im.Bounds().Max.Y; yi++ {
			if !isPadding(im.At(xi, yi)) {
				if out := xi - im.Bounds().Min.X; out < 0 {
					panic("padding width cannot be less than zero")
				} else {
					fmt.Printf("X %d Y %d\n", xi, yi)
					return paddingWidthReduce(out, int(retainPaddingWidth))
				}
			}
		}
	}
	return im.Bounds().Dx()
}

func PaddingWidthRight(im image.Image, isPadding IsPaddingFunc, retainPaddingWidth uint32) int {
	if im == nil {
		return 0
	}
	for xi := im.Bounds().Max.X - 1; xi >= im.Bounds().Min.X; xi-- {
		for yi := im.Bounds().Min.Y; yi < im.Bounds().Max.Y; yi++ {
			if !isPadding(im.At(xi, yi)) {
				if out := im.Bounds().Max.X - 1 - xi; out < 0 {
					panic("padding width cannot be less than zero")
				} else {
					return paddingWidthReduce(out, int(retainPaddingWidth))
				}
			}
		}
	}
	return im.Bounds().Dx()
}

func paddingWidthReduce(paddingWidth int, retainPadding int) int {
	if paddingWidth < 0 {
		paddingWidth = 0
	}
	if retainPadding < 0 {
		retainPadding = 0
	}
	if retainPadding == 0 {
		return paddingWidth
	} else if modWidth := paddingWidth - retainPadding; modWidth < 0 {
		return 0
	} else {
		return modWidth
	}
}
