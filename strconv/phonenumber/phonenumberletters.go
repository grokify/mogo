package phonenumber

import (
	"regexp"
	"strconv"
	"strings"
)

const numbers = "2ABC3DEF4GHI5JKL6MNO7PQRS8TUV9WXYZ"

func LetterToNumberMap() map[string]int {
	rx := regexp.MustCompile(`(\d)([A-Z]+)`)
	m := rx.FindAllStringSubmatch(numbers, -1)

	l2n := map[string]int{}

	for _, m2 := range m {
		if len(m2) >= 3 {
			num, err := strconv.Atoi(m2[1])
			if err != nil {
				panic(err)
			}

			letters := strings.Split(m2[2], "")

			for _, item := range letters {
				l2n[item] = num
			}
		}
	}

	return l2n
}

func StringToNumbers(s string) string {
	conv := []string{}
	arr := strings.Split(s, "")
	l2n := LetterToNumberMap()
	for _, letter := range arr {
		letter = strings.ToUpper(letter)
		if num, ok := l2n[letter]; ok {
			conv = append(conv, strconv.Itoa(num))
		} else {
			conv = append(conv, letter)
		}
	}
	return strings.Join(conv, "")
}
