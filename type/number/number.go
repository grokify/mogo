package number

func SliceInt64ToFloat64(src []int64) []float64 {
	out := []float64{}
	for _, in := range src {
		out = append(out, float64(in))
	}
	return out
}
