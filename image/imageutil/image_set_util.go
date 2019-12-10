package imageutil

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
)

func MergeHorizontalRGBA(ims ImageMetaSet) image.Image {
	mergedRGBA := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{ims.SumX(-1), ims.MaxY()}})

	for i, im := range ims.ImageMetas {
		if i == 0 {
			draw.Draw(mergedRGBA, im.Image.Bounds(), im.Image, image.Point{0, 0}, draw.Src)
		} else {
			startingPostionI := image.Point{ims.SumX(i - 1), 0}
			rectangleI := image.Rectangle{
				startingPostionI,
				startingPostionI.Add(im.Image.Bounds().Size())}
			draw.Draw(mergedRGBA, rectangleI, im.Image, image.Point{0, 0}, draw.Src)
		}
	}
	return mergedRGBA
}
