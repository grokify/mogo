package ratio

// RatioInt generates the missing value in a ratio calculation.
// Missing number should be in second set of coordinates.
func RatioInt(x1, y1, x2, y2 int) (int, int) {
	if x2 > 0 && y2 > 0 {
		return x2, y2
	} else if x2 <= 0 && y2 <= 0 {
		return x1, y1
	} else if x2 <= 0 {
		return int((float64(x1) / float64(y1)) * float64(y2)), y2
	} else if y2 <= 0 {
		return x2, int((float64(y1) / float64(x1)) * float64(x2))
	}
	return x2, y2
}
