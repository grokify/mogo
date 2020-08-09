package imageutil

import "image"

func ImageAspect(img image.Image) float64 {
	return Aspect(ImageWidthHeight(img))
}

func Aspect(width, height int) float64 {
	return float64(width) / float64(height)
}

func ImageWidthHeight(img image.Image) (int, int) {
	return RectangleWidthHeight(img.Bounds())
}

func RectangleWidthHeight(rect image.Rectangle) (int, int) {
	return rect.Max.X - rect.Min.X,
		rect.Max.Y - rect.Min.Y
}
