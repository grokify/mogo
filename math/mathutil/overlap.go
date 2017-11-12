package mathutil

func IsOverlapSortedInt(x1, x2, y1, y2 int) bool {
	return x1 <= y2 && y1 <= x2
}

func IsOverlapSortedInt32(x1, x2, y1, y2 int32) bool {
	return x1 <= y2 && y1 <= x2
}

func IsOverlapSortedInt64(x1, x2, y1, y2 int64) bool {
	return x1 <= y2 && y1 <= x2
}

func IsOverlapUnsortedInt(x1, x2, y1, y2 int) bool {
	return (x1 >= y1 && x1 <= y2) ||
		(x2 >= y1 && x2 <= y2) ||
		(y1 >= x1 && y1 <= x2) ||
		(y2 >= x1 && y2 <= x2)
}
