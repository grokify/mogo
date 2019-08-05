package timeutil

func DayofmonthToEnglish(i uint16) string {
	days := []string{
		"zeorth",
		"first", "second", "third", "fourth", "fifth",
		"sixth", "seventh", "eighth", "ninth", "tenth",
		"eleventh", "twelfth", "thirteenth", "fourteenth", "fifteenth",
		"sixteenth", "seventeenth", "eighteenth", "nineteenth", "twentieth",
	}
	tenZero := []string{"tenth", "twentieth", "thirtieth", "fourtieth", "fiftieth"}
	tenPlus := []string{"ten", "twenty", "thirty", "fourty", "fifty"}
	if i < 21 {
		return days[i]
	}
	if i > 59 {
		panic("E_OUT_OF_RANGE")
	}
	quotient, remainder := i/10, i%10
	if remainder == 0 {
		return tenZero[quotient-1]
	}
	return tenPlus[quotient-1] + " " + days[remainder]
}
