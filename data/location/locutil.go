package location

import (
	"github.com/grokify/mogo/strconv/strconvutil"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// LatLngString returns a string.
func LatLngString(loc latlng.LatLng, sep string, precision uint) string {
	return strconvutil.FormatDecimal(loc.Latitude, 0) + sep + strconvutil.FormatDecimal(loc.Longitude, precision)
}
