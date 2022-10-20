package timezone

import (
	"testing"
)

var locationTzAbbreviationAndNameTests = []struct {
	tzAbbrReq  string
	tzNameReq  string
	tzNameResp string
	tzOffet    int
	isErr      bool
}{
	{"PDT", "Pacific Daylight Time", "Pacific Daylight Time", -7 * 60, false},
	{"PST", "Pacific Standard Time", "Pacific Standard Time", -8 * 60, false},
	{"CAT", "", "Central Africa Time", 60, false},
	{"CAT", "Central Africa Time", "Central Africa Time", 60, false},
	{"CET", "", "Central European Time/Central European Standard Time", 60, false},
}

// TestLoadLocationByAbbreviation tests retrieving time zones.
func TestLoadLocationByAbbreviation(t *testing.T) {
	for _, tt := range locationTzAbbreviationAndNameTests {
		loc, err := LoadLocationByAbbreviation(tt.tzAbbrReq, tt.tzNameReq)
		if err != nil {
			if tt.isErr {
				continue
			} else {
				t.Errorf("timezone.LoadLocationFromAbbreviation(\"%s\", \"%s\") Error [%s]", tt.tzAbbrReq, tt.tzNameReq, err.Error())
			}
		}
		if loc.String() != tt.tzNameResp {
			t.Errorf("timezone.LoadLocationFromAbbreviation(\"%s\", \"%s\") Want [%s] Got [%s]", tt.tzAbbrReq, tt.tzNameReq, tt.tzNameResp, loc.String())
		}
	}
}
