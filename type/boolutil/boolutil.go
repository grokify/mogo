package boolutil

func Flip(b bool) bool {
	return !b
}

func ToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
