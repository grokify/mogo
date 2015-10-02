package mailutil

import (
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
