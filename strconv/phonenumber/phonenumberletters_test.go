package phonenumber

import (
	"testing"
)

// numbers = "2ABC3DEF4GHI5JKL6MNO7PQRS8TUV9WXYZ"

var letterToNumberTests = []struct {
	v    string
	want int
}{
	{"A", 2}, {"B", 2}, {"C", 2}, {"D", 3}, {"E", 3},
	{"F", 3}, {"G", 4}, {"H", 4}, {"I", 4}, {"J", 5},
	{"K", 5}, {"L", 5}, {"M", 6}, {"N", 6}, {"O", 6},
	{"P", 7}, {"Q", 7}, {"R", 7}, {"S", 7}, {"T", 8},
	{"U", 8}, {"V", 8}, {"W", 9}, {"X", 9}, {"Y", 9}, {"Z", 9},
}

func TestLetterToNumber(t *testing.T) {
	l2n := LetterToNumberMap()

	for _, tt := range letterToNumberTests {
		if num, ok := l2n[tt.v]; ok {
			if num != tt.want {
				t.Errorf("phonenumber.LetterToNumberMap() Error: with [%v], want [%v], got [%v]",
					tt.v, tt.want, num)
			}

		} else {
			t.Errorf("phonenumber.LetterToNumberMap() Not Found Error: with [%v], want [%v]",
				tt.v, tt.want)
		}
	}
}

var stringToNumbersTests = []struct {
	v    string
	want string
}{
	{"gotmilk", "4686455"},
}

func TestStringToNumbers(t *testing.T) {
	for _, tt := range stringToNumbersTests {
		nums := StringToNumbers(tt.v)

		if nums != tt.want {
			t.Errorf("phonenumber.StringToNumbers() Error: with [%v], want [%v], got [%v]",
				tt.v, tt.want, nums)
		}
	}
}
