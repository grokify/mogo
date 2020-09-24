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

func MergeXSameYHttp(url1, url2 string, larger bool) (image.Image, error) {
	img1, _, err := ReadImageHttp(url1)
	if err != nil {
		return img1, err
	}
	img2, _, err := ReadImageHttp(url2)
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

func MergeYSameXHttp(url1, url2 string, larger bool) (image.Image, error) {
	img1, _, err := ReadImageHttp(url1)
	if err != nil {
		return img1, err
	}
	img2, _, err := ReadImageHttp(url2)
	if err != nil {
		return img2, err
	}
	return MergeYSameX(img1, img2, true), nil
}

func Merge4Http(url1, url2, url3, url4 string, larger bool) (image.Image, error) {
	img12, err := MergeXSameYHttp(url1, url2, larger)
	if err != nil {
		return img12, err
	}

	img34, err := MergeXSameYHttp(url3, url4, larger)
	if err != nil {
		return img34, err
	}

	return MergeYSameX(img12, img34, larger), nil
}
