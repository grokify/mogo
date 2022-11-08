package smtputil

import (
	"testing"
)

var addrFormatTests = []struct {
	regular string
	swapped string
	reverse string
}{
	{"info", "info", "info"},
	{"info@example.com", "example.com@@info", "com.example@@@info"},
	{"info@www.example.com", "www.example.com@@info", "com.example.www@@@info"},
	{"info@blog.www.example.com", "blog.www.example.com@@info", "com.example.www.blog@@@info"},
}

func TestAddressFormats(t *testing.T) {
	for _, tt := range addrFormatTests {
		trySwapped, err := EmailAddrToSwapped(tt.regular)
		if err != nil {
			t.Errorf("smtputil.EmailAddrToSwapped(\"%s\") Error: input [%v], want [%v], got error [%v]",
				tt.regular, tt.regular, tt.swapped, err.Error())
		}
		if trySwapped != tt.swapped {
			t.Errorf("smtputil.EmailAddrToSwapped(\"%s\") Failure: input [%v], want [%v], got [%v]",
				tt.regular, tt.regular, tt.swapped, trySwapped)
		}
		tryReverse, err := EmailAddrToReverse(tt.regular)
		if err != nil {
			t.Errorf("smtputil.EmailAddrToReverse(\"%s\") Error: input [%v], want [%v], got error [%v]",
				tt.regular, tt.regular, tt.reverse, err.Error())
		}
		if tryReverse != tt.reverse {
			t.Errorf("smtputil.EmailAddrToReverse(\"%s\") Failure: input [%v], want [%v], got [%v]",
				tt.regular, tt.regular, tt.reverse, tryReverse)
		}
	}
}
