package colors

import (
	"image/color"
)

type Colors []color.Color

func (clrs Colors) In(c color.Color) bool {
	r, g, b, a := c.RGBA()
	for _, cx := range clrs {
		rx, gx, bx, ax := cx.RGBA()
		if r == rx && g == gx && b == bx && a == ax {
			return true
		}
	}
	return false
}
