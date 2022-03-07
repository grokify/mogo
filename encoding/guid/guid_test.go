// encoding/guid supports encoding and decoding Guid values.
package guid

import (
	"testing"
)

var guidToBase58Tests = []struct {
	v    string
	want string
}{
	{"00a646d3-9c61-4cb7-bfcd-ee2522c8f633", "5ep1jGdDdWDAcGA7TWuKg"},
	{"00a646d39c614cb7bfcdee2522c8f633", "5ep1jGdDdWDAcGA7TWuKg"},
	{"c9a646d3-9c61-4cb7-bfcd-ee2522c8f633", "RuERb5XkGtKhpCYLxK2axr"},
	{"c9a646d39c614cb7bfcdee2522c8f633", "RuERb5XkGtKhpCYLxK2axr"},
}

func TestGuidToBase58(t *testing.T) {
	for _, tt := range guidToBase58Tests {
		b58, err := GUIDToBase58(tt.v)

		if err != nil {
			t.Errorf("base10.Encode(%v): want %v, err %v", tt.v, tt.want, err.Error())
		}

		if b58 != tt.want {
			t.Errorf("base10.Encode(%v): want %v, got %v", tt.v, tt.want, b58)
		}
	}
}

var base58ToGUIDHypenTests = []struct {
	v    string
	want string
}{
	{"5ep1jGdDdWDAcGA7TWuKg", "00a646d3-9c61-4cb7-bfcd-ee2522c8f633"},
	{"RuERb5XkGtKhpCYLxK2axr", "c9a646d3-9c61-4cb7-bfcd-ee2522c8f633"},
}

func TestBase58ToGUIDHyphen(t *testing.T) {
	for _, tt := range base58ToGUIDHypenTests {
		b58, err := Base58ToGUID(tt.v, true)

		if err != nil {
			t.Errorf("guid.Base58ToGuid(%v): want %v, err %v", tt.v, tt.want, err.Error())
		}

		if b58 != tt.want {
			t.Errorf("guid.Base58ToGuid(%v): want %v, got %v", tt.v, tt.want, b58)
		}
	}
}

var base58ToGuidNoHypenTests = []struct {
	v    string
	want string
}{
	{"5ep1jGdDdWDAcGA7TWuKg", "00a646d39c614cb7bfcdee2522c8f633"},
	{"RuERb5XkGtKhpCYLxK2axr", "c9a646d39c614cb7bfcdee2522c8f633"},
}

func TestBase58ToGUIDNoHyphen(t *testing.T) {
	for _, tt := range base58ToGUIDHypenTests {
		b58, err := Base58ToGUID(tt.v, true)

		if err != nil {
			t.Errorf("guid.Base58ToGuid(%v): want %v, err %v", tt.v, tt.want, err.Error())
		}

		if b58 != tt.want {
			t.Errorf("guid.Base58ToGuid(%v): want %v, got %v", tt.v, tt.want, b58)
		}
	}
}
