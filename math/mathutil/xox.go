package mathutil

// PercentChangeToXoX converts a 1.0 == 100% based `float64` to a
// XoX percentage `float64`.
func PercentChangeToXoX(v float64) float64 {
	if v < 1.0 {
		return -1 * 100.0 * (1.0 - v)
	}
	return 100.0 * (v - 1.0)
}
