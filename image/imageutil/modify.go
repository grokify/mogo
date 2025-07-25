package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/grokify/mogo/image/colors"
	"golang.org/x/image/draw"
)

// AddBackgroundColor adds a background of `color.Color` to an image.
// It is is useful when the image has a transparent background. Use
// `colornames` for more colors, e.g. `colornames.Blue`. This returns
// a `draw.Image` so it can be used as an input to `draw.Draw()`.
func AddBackgroundColor(img image.Image, clr color.Color) draw.Image {
	imgNew := image.NewRGBA(img.Bounds())
	draw.Draw(imgNew, imgNew.Bounds(), &image.Uniform{clr}, image.Point{}, draw.Src)
	draw.Draw(imgNew, imgNew.Bounds(), img, img.Bounds().Min, draw.Over)
	return imgNew
}

// AddBackgroundWhite adds a white background which is usable
// when the image has a transparent background.
func AddBackgroundWhite(img image.Image) image.Image {
	return AddBackgroundColor(img, color.White)
}

// Bleach modifies near-white values to true-white.
func Bleach(img image.Image) image.Image {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Normalize near-white pixels to true white
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := rgba.At(x, y)
			if colors.IsNearWhite(c) {
				_, _, _, a := rgba.At(x, y).RGBA()
				a8 := uint8(a >> 8) // #nosec G115 // This is intentional truncation.

				rgba.Set(x, y, color.RGBA{255, 255, 255, a8})
			}
		}
	}
	return rgba
}

// Resize scales an image to the provided size units. Use a 0
// to scale the aspect ratio. See gitub.com/nfnt/resize for Lanczos3, etc.
// https://github.com/nfnt/resize .
func Resize(width, height int, src image.Image, scale draw.Scaler) image.Image {
	if width <= 0 && height <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	} else if width == 0 && height != 0 {
		width = int(ImageAspect(src) * float64(height))
	} else if height == 0 && width != 0 {
		height = int(float64(width) / ImageAspect(src))
	}
	rect := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(rect)
	if scale == nil {
		scale = ScalerBest()
	}
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

// ResizeMaxDimension resizes an image so that the max dimension matches what is supplied.
// If `maxDimension` is `0` or less a zero value `image.Image` is returned.
func ResizeMaxDimension(maxDimension int, src image.Image, scale draw.Scaler) image.Image {
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()
	if width > height {
		if width == maxDimension {
			return src
		} else {
			return Resize(maxDimension, 0, src, scale)
		}
	} else if height == maxDimension {
		return src
	} else {
		return Resize(0, maxDimension, src, scale)
	}
}

// ResizeMax resizes an image to maximum dimensions. To resize
// to a maximum of 800 pixels width, the following can be used:
// `ResizeMax(800, 0, img, nil)`.
func ResizeMax(maxWidth, maxHeight int, src image.Image, scale draw.Scaler) image.Image {
	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()
	if srcWidth <= maxWidth && srcHeight <= maxHeight {
		return src
	}
	outWidth := 0
	outHeight := 0
	if maxHeight == 0 {
		outWidth = maxWidth
	} else if maxWidth == 0 {
		outHeight = maxHeight
	} else {
		wRatio := float64(maxWidth) / float64(srcWidth)
		hRatio := float64(maxHeight) / float64(srcHeight)
		if wRatio < hRatio {
			outHeight = maxHeight
		} else {
			outWidth = maxWidth
		}
	}
	return Resize(outWidth, outHeight, src, scale)
}

// ResizeMin resizes an image to minimum dimensions. To resize
// to a minimum of 800 pixels width, the following can be used:
// `ResizeMin(800, 0, img, nil)`.
func ResizeMin(minWidth, minHeight int, src image.Image, scale draw.Scaler) image.Image {
	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()
	if srcWidth >= minWidth && srcHeight >= minHeight {
		return src
	}
	outWidth := 0
	outHeight := 0
	if minHeight == 0 {
		outWidth = minWidth
	} else if minWidth == 0 {
		outHeight = minHeight
	} else {
		wRatio := float64(minWidth) / float64(srcWidth)
		hRatio := float64(minHeight) / float64(srcHeight)
		if wRatio > hRatio {
			outHeight = minHeight
		} else {
			outWidth = minWidth
		}
	}
	return Resize(outWidth, outHeight, src, scale)
}

// Scale will resize the image to the provided rectangle using the
// provided interpolation function.
func Scale(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

func SliceXY(images []image.Image, maxIdx int) (minX, maxX, minY, maxY, sumX, sumY int) {
	for i, img := range images {
		if maxIdx >= 0 && i > maxIdx {
			break
		}
		sumX += img.Bounds().Dx()
		sumY += img.Bounds().Dy()
		if i == 0 {
			minX = img.Bounds().Dx()
			maxX = img.Bounds().Dx()
			minY = img.Bounds().Dy()
			maxY = img.Bounds().Dy()
			continue
		}
		if minX > img.Bounds().Dx() {
			minX = img.Bounds().Dx()
		}
		if maxX < img.Bounds().Dx() {
			maxX = img.Bounds().Dx()
		}
		if minY > img.Bounds().Dy() {
			minY = img.Bounds().Dy()
		}
		if maxY < img.Bounds().Dy() {
			maxY = img.Bounds().Dy()
		}
	}
	return
}

func ResizeSameX(images []image.Image, larger bool) []image.Image {
	// minX, maxX, _, _, _, _ := SliceXY(images, -1)
	imgs := Images(images)
	minX := imgs.DxMin()
	maxX := imgs.DxMax()
	for i, img := range images {
		if larger && img.Bounds().Dx() != maxX {
			images[i] = Resize(maxX, 0, img, ScalerBest())
		} else if !larger && img.Bounds().Dx() != minX {
			images[i] = Resize(minX, 0, img, ScalerBest())
		}
	}
	return images
}

func ResizeSameY(images []image.Image, larger bool) []image.Image {
	// _, _, minY, maxY, _, _ := SliceXY(images, -1)
	imgs := Images(images)
	minY := imgs.DyMin()
	maxY := imgs.DyMax()
	for i, img := range images {
		if larger && img.Bounds().Dy() != maxY {
			images[i] = Resize(0, maxY, img, ScalerBest())
		} else if !larger && img.Bounds().Dy() != minY {
			images[i] = Resize(0, minY, img, ScalerBest())
		}
	}
	return images
}

// ScalerDefault returns a general best results interpolation
// algorithm. See more here https://blog.codinghorror.com/better-image-resizing/ ,
// https://support.esri.com/en/technical-article/000005606 ,
// https://stackoverflow.com/questions/384991/what-is-the-best-image-downscaling-algorithm-quality-wise/6171860 .
func ScalerDefault() draw.Scaler { return draw.BiLinear }

func ScalerBest() draw.Scaler { return draw.CatmullRom }

const (
	AlgNearestNeighbor = "nearestneighbor"
	AlgApproxBiLinear  = "approxbilinear"
	AlgBiLinear        = "bilinear"
	AlgCatmullRom      = "catmullrom"
)

func ParseScaler(rawInterpolation string) (draw.Scaler, error) {
	rawInterpolation = strings.ToLower(strings.TrimSpace(rawInterpolation))
	switch rawInterpolation {
	case AlgNearestNeighbor:
		return draw.NearestNeighbor, nil
	case AlgApproxBiLinear:
		return draw.ApproxBiLinear, nil
	case AlgBiLinear:
		return draw.BiLinear, nil
	case AlgCatmullRom:
		return draw.CatmullRom, nil
	}
	return draw.NearestNeighbor, fmt.Errorf("cannot parse Scalar [%s]", rawInterpolation)
}

func PaintColor(img draw.Image, clr color.Color, area image.Rectangle) {
	if img == nil {
		return
	}
	rectImg := img.Bounds()

	for x := area.Min.X; x < area.Max.X; x++ {
		if x < rectImg.Min.X || x > rectImg.Max.X {
			continue
		}
		for y := area.Min.Y; y < area.Max.Y; y++ {
			if y < rectImg.Min.Y || y > rectImg.Max.Y {
				continue
			}
			img.Set(x, y, clr)
		}
	}
}

// AddBorder adds a border to a `draw.Image`. If you have an `image.Image`,
// first convert it with `ImageToRGBA(img)`.
func AddBorder(img draw.Image, clr color.Color, width int) draw.Image {
	if img == nil || width <= 0 {
		return img
	}
	border := width
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	w2 := w + border*2
	h2 := h + border*2
	i2 := image.NewRGBA(image.Rect(0, 0, w2-1, h2-1))
	for x := 0; x < w2; x++ {
		for y := 0; y < h2; y++ {
			if x < border || x > h+border {
				i2.Set(x, y, clr)
			} else if y < border || y > w+border {
				i2.Set(x, y, clr)
			} else {
				i2.Set(x, y, img.At(x-border, y-border))
			}
		}
	}
	return i2
}

func AddBorderAverageColor(img image.Image, width int) image.Image {
	if img == nil || width <= 0 {
		return img
	}
	imgRGBA := ImageToRGBA(img)
	return AddBorder(imgRGBA, colors.ColorAverageImage(imgRGBA), width)
}

// Information on rotation:
//
// https://www.golangprograms.com/how-to-rotate-an-image-in-golang.html
// https://code.google.com/archive/p/graphics-go/
// https://github.com/BurntSushi/graphics-go
// graphics.Rotate(dstImage, srcImage, &graphics.RotateOptions{math.Pi / 2.0}
