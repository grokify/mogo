package colors

import (
	"image/color"
)

var (
	Black  = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	Red    = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	Orange = color.RGBA{R: 255, G: 165, B: 0, A: 255}
	Green  = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	Blue   = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	White  = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	ANSIBlackHex   = "000000"
	ANSIRedHex     = "800000"
	ANSIGreenHEx   = "008000"
	ANSIYellowHex  = "808000"
	ANSIBlueHex    = "000080"
	ANSIMagentaHex = "800080"
	ANSICyanHex    = "008080"
	ANSIWhiteHex   = "c0c0c0"

	// colors from: https://github.com/badges/shields/blob/18e17233c49cf94f9a32115c3fcdc439cb495086/badge-maker/lib/color.js#L7
	ShieldBlueHex        = "007ec6"
	ShieldBrightGreenHex = "44cc11"
	ShieldGreenHex       = "97ca00"
	ShieldGreyHex        = "555555"
	ShieldLightGreyHex   = "9f9f9f"
	ShieldOrangeHex      = "fe7d37"
	ShieldRedHex         = "e05d44"
	ShieldYellowHex      = "dfb317"
	ShieldYellowGreenHex = "a4a61d"
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

func IsNearWhite(c color.Color) bool {
	// Normalize near-white pixels to true white
	nearWhite := uint8(230)
	r, g, b, _ := c.RGBA()
	// RGBA returns values in 16-bit (0-65535), so normalize to 8-bit
	r8 := uint8(r >> 8) // #nosec G115 // This is intentional truncation.
	g8 := uint8(g >> 8) // #nosec G115 // This is intentional truncation.
	b8 := uint8(b >> 8) // #nosec G115 // This is intentional truncation.

	// If it's close to white (e.g., light gray), snap to white
	if r8 > nearWhite && g8 > nearWhite && b8 > nearWhite {
		return true
	} else {
		return false
	}
}
