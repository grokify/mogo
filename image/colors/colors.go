package colors

import (
	"image/color"
)

type Colors []color.Color

func (clrs Colors) In(c color.Color) bool {
	wantHex := ColorToHex(c)
	for _, cx := range clrs {
		if ColorToHex(cx) == wantHex {
			return true
		}
	}
	return false
}
