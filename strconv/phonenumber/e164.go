package phonenumber

import (
	"strings"

	"github.com/nyaruka/phonenumbers"
)

func E164Format(numberToParse, defaultRegion string, numberFormat phonenumbers.PhoneNumberFormat) (string, error) {
	defaultRegion = strings.ToUpper(strings.TrimSpace(defaultRegion))
	if len(defaultRegion) == 0 {
		defaultRegion = "US"
	}
	phone, err := phonenumbers.Parse(numberToParse, defaultRegion)
	if err != nil {
		return "", err
	}
	return phonenumbers.Format(phone, numberFormat), nil
}
