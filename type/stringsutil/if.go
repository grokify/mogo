package stringsutil

func IfBoolString(boolVal bool, valueA, valueB string) string {
	if boolVal {
		return valueA
	}
	return valueB
}
