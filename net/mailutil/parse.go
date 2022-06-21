package mailutil

import (
	"net/mail"
	"regexp"
	"sort"
	"strings"
)

var rxAddressRFC5322Capture = regexp.MustCompile(`<([^><]+?)>`)

type Addresses []*mail.Address

func (addrs Addresses) Strings(smtpOnly, includeAngleBrackets, smtpToLower, sortAsc bool) []string {
	strs := []string{}
	for _, addr := range addrs {
		if addr == nil {
			continue
		}
		str := strings.TrimSpace(addr.String())
		if len(str) == 0 {
			continue
		}
		if smtpOnly {
			m := rxAddressRFC5322Capture.FindStringSubmatch(str)
			if len(m) > 0 {
				if includeAngleBrackets {
					str = strings.TrimSpace(m[1])
				} else {
					str = strings.TrimSpace(m[0])
				}
				if smtpToLower {
					str = strings.ToLower(str)
				}
			}
		}
		strs = append(strs, str)
	}
	if sortAsc {
		sort.Strings(strs)
	}
	return strs
}

/*
// ParseMulti will parse multiple email addresses from a string using
// RFC 5322 angle brackets.
func ParseMulti(input string) (Addresses, error) {
	addrs := Addresses{}
	m := rxAddressRFC5322Capture.FindAllStringSubmatch(input, -1)
	if len(m) == 0 {
		return addrs, nil
	}
	for _, mx := range m {
		try := strings.TrimSpace(mx[1])
		if len(try) == 0 {
			continue
		}
		addr, err := mail.ParseAddress(try)
		if err != nil {
			return addrs, err
		}
		if addr != nil {
			addrs = append(addrs, addr)
		}
	}
	return addrs, nil
}
*/
