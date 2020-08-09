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
// to scale the aspect ratio. See gitub.com/nfnt/resize for Lanczos3, etc.
// https://github.com/nfnt/resize .
func Resize(width, height uint, src image.Image, scale draw.Scaler) image.Image {
	if width == 0 && height != 0 {
		width = uint(ImageAspect(src) * float64(height))
	} else if height == 0 && width != 0 {
		height = uint(float64(width) / ImageAspect(src))
	}
	rect := image.Rect(0, 0, int(width), int(height))
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

// Scale will resize the image to the provided rectangle using the
// provided interpolation function.
func Scale(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

// DefaultScaler returns a general best results interpolation
// algorithm. See more here https://blog.codinghorror.com/better-image-resizing/ ,
// https://support.esri.com/en/technical-article/000005606 ,
// https://stackoverflow.com/questions/384991/what-is-the-best-image-downscaling-algorithm-quality-wise/6171860 .
func DefaultScaler() draw.Scaler { return draw.BiLinear }

func BestScaler() draw.Scaler { return draw.CatmullRom }

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
