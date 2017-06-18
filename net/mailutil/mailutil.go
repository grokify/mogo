package mailutil

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

const (
	EmailFull                string = `^(?i)([^\@]+)\@([^\@]+\.[a-z]+)+$`
	EmailFullNumeric         string = `(?i)^([0-9]+|test)\@[0-9]+\.[a-z]+$`
	EmailDomainExampleOrTest string = `(?i)[\@.](example|test)\.[a-z]+$`
	DomainSingleCharValid    string = `(?i)^([qxz]\.com|[iq]\.net|[cvwx]\.org|)$`
)

var (
	rxEmailFull                = regexp.MustCompile(EmailFull)
	rxEmailFullNumeric         = regexp.MustCompile(EmailFullNumeric)
	rxEmailDomainExampleOrTest = regexp.MustCompile(EmailDomainExampleOrTest)
	rxDomainSingleCharValid    = regexp.MustCompile(DomainSingleCharValid)
)

func AddressIsValidFullFuzzy(address string, excludeExampleOrTest bool, excludeNumericTestAddress bool) bool {
	address = strings.Trim(address, " ")
	if len(address) < 6 { // i@g.cn
		return false
	}
	address = strings.ToLower(address)
	valid, _, hostname := AddressIsValidFull(address)
	if !valid {
		return false
	}
	valid = govalidator.IsEmail(address)
	if !valid {
		return false
	}
	valid = HostnameIsValid(hostname)
	if !valid {
		return false
	}
	if excludeExampleOrTest {
		test := DomainIsExampleOrTest(address)
		if test {
			return false
		}
	}
	if excludeNumericTestAddress {
		rsEmailFullNumeric := rxEmailFullNumeric.FindString(address)
		if len(rsEmailFullNumeric) > 0 {
			return false
		}
	}
	return true
}

func AddressIsValidFull(address string) (bool, string, string) {
	rsEmailFull := rxEmailFull.FindStringSubmatch(address)
	if len(rsEmailFull) > 0 {
		return true, rsEmailFull[1], rsEmailFull[2]
	}
	return false, "", ""
}

func DomainIsExampleOrTest(address string) bool {
	rsEmailDomainExampleOrTest := rxEmailDomainExampleOrTest.FindStringSubmatch(address)
	if len(rsEmailDomainExampleOrTest) > 0 {
		return true
	}
	return false
}

func HostnameIsValid(hostname string) bool {
	parts := strings.Split(hostname, ".")
	if len(parts) < 2 {
		return false
	}
	if len(parts) == 2 && len(parts[0]) == 1 {
		return DomainIsValidSingleChar(hostname)
	}
	return true
}

func DomainIsValidSingleChar(domain string) bool {
	rs := rxDomainSingleCharValid.FindStringSubmatch(domain)
	if len(rs) > 0 {
		return true
	}
	return false
}

type MailAddress struct {
	SMTPUser string
	SMTPHost string
	Address  *mail.Address
}

func ParseAddress(address string) (MailAddress, error) {
	ma := MailAddress{}
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return ma, err
	}
	addr.Address = strings.ToLower(strings.TrimSpace(addr.Address))
	ma.Address = addr

	if len(addr.Address) > 2 {
		user, host, err := ParseAddressSpec(addr.Address)
		if err == nil {
			ma.SMTPUser = user
			ma.SMTPHost = host
		}
	}
	return ma, nil
}

// ParseAddressSpec parses RFC 5322 Addr-Spec Specification
func ParseAddressSpec(addrSpec string) (string, string, error) {
	rs := regexp.MustCompile(`^([^@]+)@([^@]+)$`).FindStringSubmatch(addrSpec)
	if len(rs) < 1 {
		return "", "", errors.New("RFC 5322 Address Spec not found.")
	}
	return rs[1], rs[2], nil
}
