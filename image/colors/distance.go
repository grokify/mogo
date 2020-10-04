package colors

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

func DistanceCIE94(color1, color2 color.RGBA) float64 {
	clr1 := ColorfulColor(color1)
	clr2 := ColorfulColor(color2)
	return clr1.DistanceCIE94(clr2)
}

func ColorfulColor(clr color.RGBA) colorful.Color {
	return colorful.Color{
		R: float64(clr.R) / 255.0,
		G: float64(clr.G) / 255.0,
		B: float64(clr.B) / 255.0}
}
