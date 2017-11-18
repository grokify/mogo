package colorutil

import (
	"github.com/lucasb-eyer/go-colorful"
	"google.golang.org/api/slides/v1"
)

// ChartColor1 is the color palette for Google Charts as collected by
// Craig Davis here: https://gist.github.com/there4/2579834
var ChartColor1 = [...]string{
	"#3366CC",
	"#DC3912",
	"#FF9900",
	"#109618",
	"#990099",
	"#3B3EAC",
	"#0099C6",
	"#DD4477",
	"#66AA00",
	"#B82E2E",
	"#316395",
	"#994499",
	"#22AA99",
	"#AAAA11",
	"#6633CC",
	"#E67300",
	"#8B0707",
	"#329262",
	"#5574A6",
	"#3B3EAC",
}

func GoogleSlidesRgbColorParseHex(hexColor string) (slides.RgbColor, error) {
	rgbColor := slides.RgbColor{}
	c, err := colorful.Hex(hexColor)
	if err != nil {
		return rgbColor, err
	}
	return slides.RgbColor{Red: c.R, Green: c.G, Blue: c.B}, nil
}
