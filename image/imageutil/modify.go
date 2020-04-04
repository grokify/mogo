package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"golang.org/x/image/draw"
)

// AddBackgroundWhite adds a white background which is usable
// when the image has a transparent background.
func AddBackgroundWhite(imgSrc image.Image) image.Image {
	imgNew := image.NewRGBA(imgSrc.Bounds())
	draw.Draw(imgNew, imgNew.Bounds(),
		&image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(imgNew, imgNew.Bounds(),
		imgSrc, imgSrc.Bounds().Min, draw.Over)
	return imgNew
}

// Resize scales an image to the provided size units. Use a 0
// to scale the aspect ratio.
func Resize(width, height uint, src image.Image, scale draw.Scaler) image.Image {
	if width == 0 && height != 0 {
		aspect := ImageAspect(src)
		width = uint(aspect * float64(height))
	} else if height == 0 && width != 0 {
		aspect := ImageAspect(src)
		height = uint(float64(width) / aspect)
	}
	rect := image.Rect(0, 0, int(width), int(height))
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

func ImageAspect(img image.Image) float64 {
	return Aspect(ImageWidthHeight(img))
}

func Aspect(width, height int) float64 {
	return float64(width) / float64(height)
}

func ImageWidthHeight(img image.Image) (int, int) {
	return img.Bounds().Max.X - img.Bounds().Min.X,
		img.Bounds().Max.Y - img.Bounds().Min.Y
}

// Scale will resize the image to thee provided rectangle using the
// provided interpolation function.
func Scale(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

func DefaultScaler() draw.Scaler {
	return draw.NearestNeighbor
}

func ParseScalar(rawInterpolation string) (draw.Scaler, error) {
	rawInterpolation = strings.ToLower(strings.TrimSpace(rawInterpolation))
	switch rawInterpolation {
	case "nearestneighbor":
		return draw.NearestNeighbor, nil
	case "approxbilinear":
		return draw.ApproxBiLinear, nil
	case "bilinear":
		return draw.BiLinear, nil
	case "catmullrom":
		return draw.CatmullRom, nil
	}
	return draw.NearestNeighbor, fmt.Errorf("Cannot parse Scalar [%s]", rawInterpolation)
}
