package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"strings"

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

func debugImagesDimensions(note string, images []image.Image) {
	dimensions := []string{note}
	for i, img := range images {
		dimensions = append(dimensions,
			fmt.Sprintf("X%d[%v]y%d[%v]",
				i+1, img.Bounds().Dx(), i+1, img.Bounds().Dy()))
	}
	fmt.Printf(strings.Join(dimensions, " ") + "\n")
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
	for _, img := range images {
		if IsNilOrEmpty(img) {
			continue
		}
		draw.Draw(output, img.Bounds().Add(image.Pt(sumXPrev, 0)),
			img, image.Point{}, draw.Src)
		sumXPrev += img.Bounds().Dx()
	}
	return output
}

func MergeXSameYRead(locations []string, larger bool) (image.Image, error) {
	images, err := ReadImages(locations)
	if err != nil {
		return nil, err
	}
	return MergeXSameY(images, true), nil
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
	for _, img := range images {
		if IsNilOrEmpty(img) {
			continue
		}
		draw.Draw(output, img.Bounds().Add(image.Pt(0, sumYPrev)),
			img, image.Point{}, draw.Src)
		sumYPrev += img.Bounds().Dy()
	}
	return output
}

func MergeYSameXRead(locations []string, larger bool) (image.Image, error) {
	images, err := ReadImages(locations)
	if err != nil {
		return nil, err
	}
	return MergeYSameX(images, true), nil
}

type Matrix [][]image.Image

// AddBackgroundColor adds a background of `color.Color` to the images.
// It is is useful when the image has a transparent background. Use
// `colornames` for more colors, e.g. `colornames.Blue`.
func (matrix Matrix) AddBackgroundColor(clr color.Color) {
	for i, row := range matrix {
		for j, img := range row {
			matrix[i][j] = AddBackgroundColor(img, clr)
		}
	}
}

// Merge combines a set of images resizing each row element's
// height for consistent rows, an each row's width for consistent
// widths.
func (matrix Matrix) Merge(largerX, largerY bool) image.Image {
	if len(matrix) == 0 {
		return nil
	}
	rowImages := []image.Image{}
	for _, rowParts := range matrix {
		if len(rowParts) > 0 {
			rowImages = append(rowImages, MergeXSameY(rowParts, largerY))
		}
	}
	if len(rowImages) == 0 {
		return nil
	}
	return MergeYSameX(rowImages, largerX)
}

func MatrixRead(imglocations [][]string) (Matrix, error) {
	matrixImages := Matrix{}
	for _, row := range imglocations {
		if len(row) == 0 {
			continue
		}
		images, err := ReadImages(row)
		if err != nil {
			return nil, err
		}
		matrixImages = append(matrixImages, images)
	}
	return matrixImages, nil
}

func MatrixMergeRead(imglocations [][]string, largerX, largerY bool) (image.Image, error) {
	matrixImages, err := MatrixRead(imglocations)
	if err != nil {
		return nil, err
	}

	return matrixImages.Merge(largerX, largerY), nil
}
