package imageutil

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// DrawText draws text onto dst in the given color. pt is the left end of the
// text's baseline, so glyphs extend above pt.Y. If face is nil,
// basicfont.Face7x13 is used. Pixels not covered by glyphs are left unchanged.
func DrawText(dst draw.Image, text string, pt image.Point, col color.Color, face font.Face) {
	if face == nil {
		face = basicfont.Face7x13
	}
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.P(pt.X, pt.Y),
	}
	d.DrawString(text)
}
