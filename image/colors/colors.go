package colors

import (
	"image/color"
)

var (
	Red    = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	Green  = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	Blue   = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	Orange = color.RGBA{R: 255, G: 165, B: 0, A: 255}
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

func Equal(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}
