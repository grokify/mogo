package phonenumber

import (
	"testing"
)

var letterToNumberTests = []struct {
	v    string
	want int
}{
	{"H", 4},
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
