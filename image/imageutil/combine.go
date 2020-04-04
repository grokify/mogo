package imageutil

import (
	"image"

	"golang.org/x/image/draw"
)

func OverlayCenterYLeftAlign(imgBg, imgOver image.Image) image.Image {
	output := image.NewRGBA(imgBg.Bounds())
	draw.Draw(output, imgBg.Bounds(), imgBg, image.ZP, draw.Src)

	_, h1 := ImageWidthHeight(imgBg)
	_, h2 := ImageWidthHeight(imgOver)
	offset := image.Pt(0, (h1-h2)/2)

	draw.Draw(output, imgOver.Bounds().Add(offset), imgOver, image.ZP, draw.Src)
	return output
}
