package location

import (
	"github.com/grokify/mogo/strconv/strconvutil"
	"google.golang.org/genproto/googleapis/type/latlng"
)

type LatLong struct {
	Latitude  float64
	Longitude float64
}

// LatLngString returns a string.
func LatLngString(loc *latlng.LatLng, sep string, precision int) string {
	if loc == nil {
		return strconvutil.FormatDecimal(0, precision) + sep + strconvutil.FormatDecimal(0, precision)
	}
	return strconvutil.FormatDecimal(loc.Latitude, precision) + sep + strconvutil.FormatDecimal(loc.Longitude, precision)
}
