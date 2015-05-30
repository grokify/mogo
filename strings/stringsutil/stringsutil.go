package stringsutil

import (
	"fmt"
)

func PadRight(str string, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}

func PadLeft(str string, pad string, length int) string {
	for {
		str = pad + str
		fmt.Println(str)
		if len(str) >= length {
			return str[0:length]
		}
	}
}
