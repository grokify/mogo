package stringsutil

// PadLeft prepends a string to a base string until the string
// length is greater or equal to the desired length.
func PadLeft(str string, pad string, length int) string {
	for {
		str = pad + str
		if len(str) >= length {
			return str[0:length]
		}
	}
}

// PadRight appends a string to a base string until the string
// length is greater or equal to the desired length.
func PadRight(str string, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}
