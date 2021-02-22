package boolutil

func Flip(b bool) bool {
	if b {
		return false
	}
	return true
}

func ToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
