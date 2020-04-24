package ratio

const (
	// Ratios are always width/height
	RatioTelevisionEarly float64 = 5 / 4 // 1.25/1
	RatioTelevision      float64 = 4 / 3 // 1.3/1
	RatioAcademy         float64 = 1.375 / 1.0
	RatioPhoto           float64 = 3 / 2
	RatioGolden          float64 = 1.6180 / 1
)

func WidthToHeight(width, ratio float64) float64 {
	return width / ratio
}

func HeightToWidth(height, ratio float64) float64 {
	return height * ratio
}
