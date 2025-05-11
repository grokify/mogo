package imageutil

import (
	"image"

	"golang.org/x/image/draw"
)

type Images []image.Image

func (imgs Images) Dimensions() []image.Point {
	points := []image.Point{}
	for _, img := range imgs {
		points = append(points, image.Point{
			X: img.Bounds().Canon().Max.X,
			Y: img.Bounds().Canon().Max.Y})
	}
	return points
}

func (imgs Images) Dxs() []int {
	dxs := []int{}
	for _, img := range imgs {
		dxs = append(dxs, img.Bounds().Dx())
	}
	return dxs
}

func (imgs Images) DxMax() int {
	dxMax := 0
	for _, img := range imgs {
		if dx := img.Bounds().Dx(); dx > dxMax {
			dxMax = dx
		}
	}
	return dxMax
}

func (imgs Images) DxMin() int {
	dxMin := 0
	for i, img := range imgs {
		if i == 0 {
			dxMin = img.Bounds().Dx()
		} else if dx := img.Bounds().Dx(); dx < dxMin {
			dxMin = dx
		}
	}
	return dxMin
}

// DxSum returns the sum of widths up to and including
// `maxIndexInclusive`. Use a negative value for `maxIndexInclusive`
// to include all elements.
func (imgs Images) DxSum(maxIndexInclusive int) int {
	dxSum := 0
	for i, img := range imgs {
		if maxIndexInclusive >= 0 && i > maxIndexInclusive {
			break
		}
		dxSum += img.Bounds().Dx()
	}
	return dxSum
}

func (imgs Images) Dys() []int {
	dys := []int{}
	for _, img := range imgs {
		dys = append(dys, img.Bounds().Dy())
	}
	return dys
}

func (imgs Images) DyMax() int {
	dyMax := 0
	for _, img := range imgs {
		if dy := img.Bounds().Dy(); dy > dyMax {
			dyMax = dy
		}
	}
	return dyMax
}

func (imgs Images) DyMin() int {
	dyMin := 0
	for i, img := range imgs {
		if i == 0 {
			dyMin = img.Bounds().Dy()
		} else if dy := img.Bounds().Dy(); dy < dyMin {
			dyMin = dy
		}
	}
	return dyMin
}

// DySum returns the sum of heights up to and including
// `maxIndexInclusive`. Use a negative value for `maxIndexInclusive`
// to include all elements.
func (imgs Images) DySum(maxIndexInclusive int) int {
	dySum := 0
	for i, img := range imgs {
		if maxIndexInclusive >= 0 && i > maxIndexInclusive {
			break
		}
		dySum += img.Bounds().Dy()
	}
	return dySum
}

// ConsistentSize resizes and crops the images so that they have all
// the same size. It prioritize resizing images to max Dx and then
// cropping Dy so they are consistent.
func (imgs Images) ConsistentSize(scale draw.Scaler, yAlign string) {
	if len(imgs) == 0 {
		return
	}
	dxMax := imgs.DxMax()
	if dxMax <= 0 {
		return
	}
	for i, img := range imgs {
		if img.Bounds().Dx() != dxMax {
			imgs[i] = Resize(dxMax, 0, img, scale)
		}
	}
	dyMin := imgs.DyMin()
	if dyMin <= 0 {
		return
	}
	for i, img := range imgs {
		if img.Bounds().Dy() != dyMin {
			imgs[i] = CropY(img, uint(dyMin), yAlign)
		}
	}
}

func (imgs Images) Stats() ImagesStats {
	return ImagesStats{
		Dxs:   imgs.Dxs(),
		DxMax: imgs.DxMax(),
		DxMin: imgs.DxMin(),
		DxSum: imgs.DxSum(-1),
		Dys:   imgs.Dys(),
		DyMax: imgs.DyMax(),
		DyMin: imgs.DyMin(),
		DySum: imgs.DySum(-1)}
}
