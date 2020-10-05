package colors

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

func DistanceCIE2K(color1, color2 color.RGBA) float64 {
	return ColorfulColor(color1).DistanceCIEDE2000(ColorfulColor(color2))
}

func DistanceCIE94(color1, color2 color.RGBA) float64 {
	return ColorfulColor(color1).DistanceCIE94(ColorfulColor(color2))
}

func DistanceCIE76(color1, color2 color.RGBA) float64 {
	return ColorfulColor(color1).DistanceCIE76(ColorfulColor(color2))
}

func ColorfulColor(clr color.RGBA) colorful.Color {
	return colorful.Color{
		R: float64(clr.R) / 255.0,
		G: float64(clr.G) / 255.0,
		B: float64(clr.B) / 255.0}
}
