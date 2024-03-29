package imageutil

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/grokify/mogo/image/colors"
)

func MergeXSameY(images []image.Image, larger bool) image.Image {
	if len(images) == 0 {
		return nil
	} else if len(images) == 1 {
		return images[0]
	}
	images = ResizeSameY(images, larger)
	// _, _, minY, _, sumX, _ := SliceXY(images, -1)
	imgs := Images(images)
	minY := imgs.DyMin()
	sumX := imgs.DxSum(-1)
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
	// minX, _, _, _, _, sumY := SliceXY(images, -1)
	imgs := Images(images)
	minX := imgs.DxMin()
	sumY := imgs.DySum(-1)
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

func (matrix Matrix) AddBackgroundColorHex(hexcolor string) error {
	clr, err := colors.ParseHex(hexcolor)
	if err != nil {
		return err
	}
	matrix.AddBackgroundColor(clr)
	return nil
}

// Merge combines a set of images resizing each row element's height
// for consistent rows, an each row's width for consistent widths.
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

func (matrix Matrix) DxSumMax() int {
	dxMax := 0
	for _, row := range matrix {
		rowDxSum := Images(row).DxSum(-1)
		if rowDxSum > dxMax {
			dxMax = rowDxSum
		}
	}
	return dxMax
}

func (matrix Matrix) DxSumMin() int {
	dxMin := 0
	for i, row := range matrix {
		if i == 0 {
			dxMin = Images(row).DxSum(-1)
		}
		rowDxSum := Images(row).DxSum(-1)
		if rowDxSum < dxMin {
			dxMin = rowDxSum
		}
	}
	return dxMin
}

func (matrix Matrix) Dimensions() [][]image.Point {
	points := [][]image.Point{}
	for r, row := range matrix {
		points = append(points, []image.Point{})
		for _, img := range row {
			points[r] = append(points[r], img.Bounds().Canon().Max)
		}
	}
	return points
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

// MergeXRGBA returns an image that is composed
// of non-overlapping imagers in ImageMetaSet. The code here
// is adapted from https://stackoverflow.com/a/35965499/1908967
func MergeXRGBA(imgs []image.Image) image.Image {
	imgs2 := Images(imgs)
	mergedRGBA := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{imgs2.DxSum(-1), imgs2.DyMax()}})

	for i, img := range imgs2 {
		if i == 0 {
			draw.Draw(mergedRGBA, img.Bounds(), img, image.Point{0, 0}, draw.Src)
		} else {
			startingPostionI := image.Point{imgs2.DxSum(i - 1), 0}
			rectangleI := image.Rectangle{
				startingPostionI,
				startingPostionI.Add(img.Bounds().Size())}
			draw.Draw(mergedRGBA, rectangleI, img, image.Point{0, 0}, draw.Src)
		}
	}
	return mergedRGBA
}
