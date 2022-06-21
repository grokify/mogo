package mailutil

import (
	"net/mail"
	"regexp"
	"sort"
	"strings"

	"github.com/grokify/mogo/fmt/fmtutil"
)

type Addresses []*mail.Address

func (addrs Addresses) Strings(includeAngleBrackets, toLower, sortAsc bool) []string {
	strs := []string{}
	for _, addr := range addrs {
		if addr == nil {
			continue
		}
		str := strings.TrimSpace(addr.String())
		if len(str) == 0 {
			continue
		}
		if toLower {
			str = strings.ToLower(str)
		}
		if !includeAngleBrackets {
			str = strings.Trim(str, "<>")
		}
		strs = append(strs, str)
	}
	if sortAsc {
		sort.Strings(strs)
	}
	return strs
}

var rxAddressRFC5322Capture = regexp.MustCompile(`<([^><]+?)>`)

// ParseMulti will parse multiple email addresses from a string using
// RFC 5322 angle brackets.
func ParseMulti(input string) (Addresses, error) {
	addrs := Addresses{}
	m := rxAddressRFC5322Capture.FindAllStringSubmatch(input, -1)
	if len(m) == 0 {
		return addrs, nil
	}
	for _, mx := range m {
		fmtutil.PrintJSON(mx)
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
