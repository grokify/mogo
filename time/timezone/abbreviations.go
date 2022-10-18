package timezone

import (
	"fmt"
	"strings"
	"time"

	timezone "github.com/tkuchiki/go-timezone"
)

const errFormatInvalidTzAbbreviation = "invalid timezone abbreviation [%s]"

// LocationTzAbbreviationAndName returns a `*time.Location` given a timezone abbreviation and name.
// For example: ("PST", "Pacific Standard Time"). This is a wrapper for `github.com/tkuchiki/go-timezone`
// which defines the timezone abbreviation and names.
func LocationTzAbbreviationAndName(tzAbbr, tzName string) (*time.Location, error) {
	tzAbbr = strings.ToUpper(strings.TrimSpace(tzAbbr))
	tzName = strings.TrimSpace(tzName)
	tz := timezone.New()
	if len(tzName) > 0 {
		tzi, err := tz.GetTzAbbreviationInfoByTZName(tzAbbr, tzName)
		if err != nil {
			return nil, err
		}
		return time.FixedZone(tzi.Name(), tzi.Offset()), nil
	}
	tzis, err := tz.GetTzAbbreviationInfo(tzAbbr)
	if err != nil {
		return nil, err
	}
	switch len(tzis) {
	case 0:
		return nil, fmt.Errorf(errFormatInvalidTzAbbreviation, tzAbbr)
	case 1:
		return time.FixedZone(tzis[0].Name(), tzis[0].Offset()), nil
	default:
		return nil, timezone.ErrAmbiguousTzAbbreviations
	}
}
