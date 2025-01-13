package mailutil

import (
	"net/mail"
	"sort"
	"strings"

	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/type/stringsutil"
)

type Addresses []mail.Address

func ParseAddressList(list string) (Addresses, error) {
	list = strings.TrimSpace(stringsutil.RemoveNonPrintable(list))
	if list == "" {
		return Addresses{}, nil
	} else if addrs, err := mail.ParseAddressList(list); err != nil {
		return Addresses{}, err
	} else if len(addrs) == 0 {
		return Addresses{}, nil
	} else {
		addrs := Addresses(pointer.DereferenceSlice(addrs))
		addrs.CanonialAddresses()
		return addrs.FilterExclEmpty(), nil
	}
}

func (addrs Addresses) CanonialAddresses() {
	for i, addr := range addrs {
		if addr2 := strings.TrimSpace(stringsutil.RemoveNonPrintable(addr.Address)); addr2 == addr.Address {
			addr.Address = addr2
			addrs[i] = addr
		}
	}
}

func (addrs Addresses) FilterExclEmpty() Addresses {
	var out Addresses
	for _, addr := range addrs {
		if addr.Address != "" || addr.Name != "" {
			out = append(out, addr)
		}
	}
	return out
}

func (addrs Addresses) FilterInclWithoutAddress() Addresses {
	var out Addresses
	for _, addr := range addrs {
		if addr.Address == "" {
			out = append(out, addr)
		}
	}
	return out
}

func (addrs Addresses) FilterInclWithAddress() Addresses {
	var out Addresses
	for _, addr := range addrs {
		if addr.Address != "" {
			out = append(out, addr)
		}
	}
	return out
}

func (addrs Addresses) TrimSpace(trimName, trimAddress bool) {
	for i := range addrs {
		if trimName {
			addrs[i].Name = strings.TrimSpace(addrs[i].Name)
		}
		if trimAddress {
			addrs[i].Address = strings.TrimSpace(addrs[i].Address)
		}
	}
}

func (addrs Addresses) Strings(smtpOnly, smtpToLower, sortAsc bool) []string {
	strs := []string{}
	for _, addr := range addrs {
		var str string
		if smtpOnly {
			if smtpToLower {
				str = strings.ToLower(addr.Address)
			} else {
				str = addr.Address
			}
		} else {
			str = addr.String()
		}
		strs = append(strs, str)
	}
	if sortAsc {
		sort.Strings(strs)
	}
	return strs
}

func (addrs Addresses) String(smtpOnly, smtpToLower, sortAsc bool) string {
	strs := addrs.Strings(smtpOnly, smtpToLower, sortAsc)
	return strings.Join(strs, AddrSep)
}

/*
var rxAddressRFC5322Capture = regexp.MustCompile(`<([^><]+?)>`)

// ParseMulti will parse multiple email addresses from a string using
// RFC 5322 angle brackets. `mail.ParseAddressList()` will only handle
// comma delimiters.
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
