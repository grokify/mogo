package location

import (
	"github.com/grokify/mogo/strconv/strconvutil"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// LatLngString returns a string.
func LatLngString(loc *latlng.LatLng, sep string, precision uint) string {
	if loc == nil {
		return strconvutil.FormatDecimal(0, precision) + sep + strconvutil.FormatDecimal(0, precision)
	}
	return strconvutil.FormatDecimal(loc.Latitude, precision) + sep + strconvutil.FormatDecimal(loc.Longitude, precision)
}
