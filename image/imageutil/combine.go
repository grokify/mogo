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

func MergeXSameY(images []image.Image, larger bool) image.Image {
	if len(images) == 0 {
		return nil
	} else if len(images) == 1 {
		return images[0]
	}
	images = ResizeSameY(images, larger)
	_, _, minY, _, sumX, _ := SliceXY(images, -1)
	output := image.NewRGBA(image.Rect(0, 0, sumX, minY))
	sumXPrev := 0
	for i, img := range images {
		if i == 0 {
			draw.Draw(output, img.Bounds(), img, image.Point{}, draw.Src)
		} else {
			imgOffset := image.Pt(sumXPrev, 0)
			draw.Draw(output, img.Bounds().Add(imgOffset), img, image.Point{}, draw.Src)
		}
		sumXPrev += img.Bounds().Dx()
	}
	return output
}

func mergeXSameYTwo(img1, img2 image.Image, larger bool) image.Image {
	img1, img2 = ResizeSameYTwo(img1, img2, larger)
	output := image.NewRGBA(
		image.Rect(0, 0,
			img1.Bounds().Dx()+img2.Bounds().Dx(),
			img1.Bounds().Dy()))
	draw.Draw(output, img1.Bounds(), img1, image.Point{}, draw.Src)
	img2Offset := image.Pt(img1.Bounds().Dx(), 0)
	draw.Draw(output, img2.Bounds().Add(img2Offset), img2, image.Point{}, draw.Src)
	return output
}

func MergeXSameYRead(locations []string, larger bool) (image.Image, error) {
	images, err := ReadImages(locations)
	if err != nil {
		return nil, err
	}
	return MergeXSameY(images, true), nil
	/*
		img1, _, err := ReadImageAny(location1)
		if err != nil {
			return img1, err
		}
		img2, _, err := ReadImageAny(location2)
		if err != nil {
			return img2, err
		}
		return MergeXSameY(img1, img2, true), nil
	*/
}

func MergeYSameX(images []image.Image, larger bool) image.Image {
	if len(images) == 0 {
		return nil
	} else if len(images) == 1 {
		return images[0]
	}
	images = ResizeSameX(images, larger)
	minX, _, _, _, _, sumY := SliceXY(images, -1)
	output := image.NewRGBA(image.Rect(0, 0, minX, sumY))
	sumYPrev := 0
	for i, img := range images {
		if i == 0 {
			draw.Draw(output, img.Bounds(), img, image.Point{}, draw.Src)
		} else {
			imgOffset := image.Pt(0, sumYPrev)
			draw.Draw(output, img.Bounds().Add(imgOffset), img, image.Point{}, draw.Src)
		}
		sumYPrev += img.Bounds().Dy()
	}
	return output
}

func mergeYSameXTwo(img1, img2 image.Image, larger bool) image.Image {
	img1, img2 = ResizeSameXTwo(img1, img2, larger)
	output := image.NewRGBA(
		image.Rect(0, 0,
			img1.Bounds().Dx(),
			img1.Bounds().Dy()+img2.Bounds().Dy()))
	draw.Draw(output, img1.Bounds(), img1, image.Point{}, draw.Src)
	img2Offset := image.Pt(0, img1.Bounds().Dy())
	draw.Draw(output, img2.Bounds().Add(img2Offset), img2, image.Point{}, draw.Src)
	return output
}

func MergeYSameXRead(locations []string, larger bool) (image.Image, error) {
	images, err := ReadImages(locations)
	if err != nil {
		return nil, err
	}
	return MergeYSameX(images, true), nil
	/*
		img1, _, err := ReadImageAny(location1)
		if err != nil {
			return img1, err
		}
		img2, _, err := ReadImageAny(location2)
		if err != nil {
			return img2, err
		}
		return MergeYSameX(img1, img2, true), nil
	*/
}

func MatrixMergeRead(matrix [][]string, largerX, largerY bool) (image.Image, error) {
	matrixImages := [][]image.Image{}
	for _, row := range matrix {
		if len(row) == 0 {
			continue
		}
		images, err := ReadImages(row)
		if err != nil {
			return nil, err
		}
		matrixImages = append(matrixImages, images)
	}

	return MatrixMerge(matrixImages, largerX, largerY), nil
}

func MatrixMerge(matrix [][]image.Image, largerX, largerY bool) image.Image {
	rowImages := []image.Image{}
	for _, rowParts := range matrix {
		rowImages = append(rowImages, MergeXSameY(rowParts, largerY))
	}
	return MergeYSameX(rowImages, largerX)
}

func merge4Read(location1, location2, location3, location4 string, larger bool) (image.Image, error) {
	img12, err := MergeXSameYRead([]string{location1, location2}, larger)
	if err != nil {
		return img12, err
	}

	img34, err := MergeXSameYRead([]string{location3, location4}, larger)
	if err != nil {
		return img34, err
	}

	return MergeYSameX([]image.Image{img12, img34}, larger), nil
}
