package imageutil

import (
	"image"
	"image/color"
	"image/draw"
)

// PaintBorder colorizes a border of size `width` in the existing
// canvas, overwriting any colors in that space.
func PaintBorder(img draw.Image, clr color.Color, width uint32) {
	if img == nil || width == 0 {
		return
	}
	PaintColor(img, clr, RectangleBorderXMin(img.Bounds(), width))
	PaintColor(img, clr, RectangleBorderXMax(img.Bounds(), width))
	PaintColor(img, clr, RectangleBorderYMin(img.Bounds(), width))
	PaintColor(img, clr, RectangleBorderYMax(img.Bounds(), width))
}

func RectangleBorderXMin(rect image.Rectangle, pixels uint32) image.Rectangle {
	maxY := rect.Min.Y + int(pixels)
	if maxY > rect.Max.Y {
		maxY = rect.Max.Y
	}
	return image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, maxY)
}

func RectangleBorderXMax(rect image.Rectangle, pixels uint32) image.Rectangle {
	minY := rect.Max.Y - int(pixels)
	if minY < rect.Min.Y {
		minY = rect.Min.Y
	}
	return image.Rect(rect.Min.X, minY, rect.Max.X, rect.Max.Y)
}

func RectangleBorderYMin(rect image.Rectangle, pixels uint32) image.Rectangle {
	maxX := rect.Min.X + int(pixels)
	if maxX > rect.Max.X {
		maxX = rect.Max.X
	}
	return image.Rect(rect.Min.X, rect.Min.Y, maxX, rect.Max.Y)
}

func RectangleBorderYMax(rect image.Rectangle, pixels uint32) image.Rectangle {
	minX := rect.Max.X - int(pixels)
	if minX < rect.Min.X {
		minX = rect.Min.X
	}
	return image.Rect(minX, rect.Min.Y, rect.Max.X, rect.Max.Y)
}
