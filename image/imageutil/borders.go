package imageutil

import (
	"image"
	"image/color"
)

func PaintBorderRGBA(img *image.RGBA, clr color.RGBA, width uint) {
	if img == nil {
		return
	}
	if width <= 0 {
		PaintColorRGBA(img, clr)
		return
	}
	borders := []image.Rectangle{}
	borders = append(borders, RectangleBorderXMin(img.Bounds(), width))
	borders = append(borders, RectangleBorderXMax(img.Bounds(), width))
	borders = append(borders, RectangleBorderYMin(img.Bounds(), width))
	borders = append(borders, RectangleBorderYMax(img.Bounds(), width))

	for _, borderRect := range borders {
		PaintColorRGBARectangle(img, clr, borderRect)
	}
}

func RectangleBorderXMin(rect image.Rectangle, pixels uint) image.Rectangle {
	maxY := rect.Min.Y + int(pixels)
	if maxY > rect.Max.Y {
		maxY = rect.Max.Y
	}
	return image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, maxY)
}

func RectangleBorderXMax(rect image.Rectangle, pixels uint) image.Rectangle {
	minY := rect.Max.Y - int(pixels)
	if minY < rect.Min.Y {
		minY = rect.Min.Y
	}
	return image.Rect(rect.Min.X, minY, rect.Max.X, rect.Max.Y)
}

func RectangleBorderYMin(rect image.Rectangle, pixels uint) image.Rectangle {
	maxX := rect.Min.X + int(pixels)
	if maxX > rect.Max.X {
		maxX = rect.Max.X
	}
	return image.Rect(rect.Min.X, rect.Min.Y, maxX, rect.Max.Y)
}

func RectangleBorderYMax(rect image.Rectangle, pixels uint) image.Rectangle {
	minX := rect.Max.X - int(pixels)
	if minX < rect.Min.X {
		minX = rect.Min.X
	}
	return image.Rect(minX, rect.Min.Y, rect.Max.X, rect.Max.Y)
}
