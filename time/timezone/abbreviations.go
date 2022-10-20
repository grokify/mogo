package timezone

import (
	"fmt"
	"strings"
	"time"

	timezone "github.com/tkuchiki/go-timezone"
)

const errFormatInvalidTzAbbreviation = "invalid timezone abbreviation [%s]"

// LoadLocationByAbbreviation returns a `*time.Location` given a timezone abbreviation and name.
// For example: ("PST", "Pacific Standard Time"). If the abbreviation only represents one timezone
// the timezone name can be omitted. An error is returned if more than one location is present.
// This is a wrapper for `github.com/tkuchiki/go-timezone`
// which defines the timezone abbreviation and names.
func LoadLocationByAbbreviation(tzAbbr, tzName string) (*time.Location, error) {
	tzAbbr = strings.ToUpper(strings.TrimSpace(tzAbbr))
	tzName = strings.TrimSpace(tzName)
	tz := timezone.New()
	if len(tzName) > 0 {
		tzi, err := tz.GetTzAbbreviationInfoByTZName(tzAbbr, tzName)
		if err != nil {
			return nil, err
		} else if tzi == nil || strings.TrimSpace(tzi.Name()) == "" {
			return nil, fmt.Errorf(errFormatInvalidTzAbbreviation, tzAbbr)
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

// LoadLocationsByAbbreviation loads all the matching locations by timezone
// abbreviation and name. Errors are returned as an empty slice.
func LoadLocationsByAbbreviation(tzAbbr, tzName string) []*time.Location {
	locs := []*time.Location{}
	tzAbbr = strings.ToUpper(strings.TrimSpace(tzAbbr))
	tzName = strings.TrimSpace(tzName)
	tz := timezone.New()
	if len(tzName) > 0 {
		tzi, err := tz.GetTzAbbreviationInfoByTZName(tzAbbr, tzName)
		if err != nil {
			return []*time.Location{}
		} else if tzi == nil || strings.TrimSpace(tzi.Name()) == "" {
			return []*time.Location{}
		}
		return append(locs, time.FixedZone(tzi.Name(), tzi.Offset()))
	}
	tzis, err := tz.GetTzAbbreviationInfo(tzAbbr)
	if err != nil {
		return locs
	}
	for _, tzi := range tzis {
		locs = append(locs, time.FixedZone(tzi.Name(), tzi.Offset()))
	}
	return locs
}
