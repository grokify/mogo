package location

import "slices"

const (
	SubdivisionTypeCanton     = "Canton"
	SubdivisionTypeACT        = "Australian Capital Territory (ACT)"
	SubdivisionTypeCapital    = "Capital"
	SubdivisonTypeCountry     = "Country"
	SubdivisionTypeCounty     = "County"
	SubdivisionTypePrefecture = "Prefecture"
	SubdivisionTypeProvince   = "Province"
	SubdivisionTypeSAR        = "Special Administrative Region"
	SubdivisionTypeSCR        = "Special Capital Region"
	SubdivisionTypeState      = "State"
)

type Locations []Location

func (locs Locations) ContainsRegionCode(rc string) bool {
	for _, loc := range locs {
		if rc == loc.RegionCode {
			return true
		}
	}
	return false
}

func (locs Locations) SubregionsCountByCountry(countries []string) int {
	count := 0
	for _, loc := range locs {
		if slices.Contains(countries, loc.ISO3166P1A2CountryCode) {
			if len(loc.Subregions) > 0 {
				count += len(loc.Subregions)
			} else if loc.SubregionsCount > 0 {
				count += loc.SubregionsCount
			}
		}
	}
	return count
}

type Location struct {
	CityName                     string
	UNLOCODE                     string
	ISO3166P1A2CountryCode       string
	ISO3166P2SubdivisionCodeFull string
	ISO3166P2CountryCode         string
	ISO3166P2SubdivisionCode     string
	ISO3166P2SubdivisionName     string
	ISO3166P2SubdivisionCategory string
	RegionType                   string
	RegionCode                   string
	RegionName                   string
	ReferenceURLs                []string
	Subregions                   []string
	SubregionsCount              int
}
