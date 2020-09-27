package imageutil

import (
	"image"

	"golang.org/x/image/draw"
)

func OverlayCenterYLeftAlign(imgBg, imgOver image.Image) image.Image {
	output := image.NewRGBA(imgBg.Bounds())
	draw.Draw(output, imgBg.Bounds(), imgBg, image.ZP, draw.Src)

	h1 := imgBg.Bounds().Dy()
	h2 := imgOver.Bounds().Dy()
	offset := image.Pt(0, (h1-h2)/2)

	draw.Draw(output, imgOver.Bounds().Add(offset), imgOver, image.Point{}, draw.Src)
	return output
}

func MergeXSameY(img1, img2 image.Image, larger bool) image.Image {
	img1, img2 = ResizeSameY(img1, img2, larger)
	output := image.NewRGBA(
		image.Rect(0, 0,
			img1.Bounds().Dx()+img2.Bounds().Dx(),
			img1.Bounds().Dy()))
	draw.Draw(output, img1.Bounds(), img1, image.Point{}, draw.Src)
	img2Offset := image.Pt(img1.Bounds().Dx(), 0)
	draw.Draw(output, img2.Bounds().Add(img2Offset), img2, image.Point{}, draw.Src)
	return output
}

func MergeXSameYRead(location1, location2 string, larger bool) (image.Image, error) {
	img1, _, err := ReadImageAny(location1)
	if err != nil {
		return img1, err
	}
	img2, _, err := ReadImageAny(location2)
	if err != nil {
		return img2, err
	}
	return MergeXSameY(img1, img2, true), nil
}

func MergeYSameX(img1, img2 image.Image, larger bool) image.Image {
	img1, img2 = ResizeSameX(img1, img2, larger)
	output := image.NewRGBA(
		image.Rect(0, 0,
			img1.Bounds().Dx(),
			img1.Bounds().Dy()+img2.Bounds().Dy()))
	draw.Draw(output, img1.Bounds(), img1, image.Point{}, draw.Src)
	img2Offset := image.Pt(0, img1.Bounds().Dy())
	draw.Draw(output, img2.Bounds().Add(img2Offset), img2, image.Point{}, draw.Src)
	return output
}

func MergeYSameXRead(location1, location2 string, larger bool) (image.Image, error) {
	img1, _, err := ReadImageAny(location1)
	if err != nil {
		return img1, err
	}
	img2, _, err := ReadImageAny(location2)
	if err != nil {
		return img2, err
	}
	return MergeYSameX(img1, img2, true), nil
}

func Merge4Read(location1, location2, location3, location4 string, larger bool) (image.Image, error) {
	img12, err := MergeXSameYRead(location1, location2, larger)
	if err != nil {
		return img12, err
	}

	img34, err := MergeXSameYRead(location3, location4, larger)
	if err != nil {
		return img34, err
	}

	return MergeYSameX(img12, img34, larger), nil
}
