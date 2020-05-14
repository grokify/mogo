package mathutil

func IntsToInt32s(ints []int) []int32 {
	int32s := []int32{}
	for _, val := range ints {
		int32s = append(int32s, int32(val))
	}
	return int32s
}

func IntsToUints(ints []int) []uint {
	uints := []uint{}
	for _, val := range ints {
		uints = append(uints, uint(val))
	}
	return uints
}
