package imageutil

import (
	"image"
	"image/gif"
)

func BuildGifAnimationSimpleRead(src *gif.GIF, delay int, names []string, f ToPalettedFunc, consistentSize bool) (*gif.GIF, error) {
	imgs, err := ReadImages(names)
	if err != nil {
		return nil, err
	}
	if consistentSize {
		imgs2 := Images(imgs)
		imgs2.ConsistentSize(ScalerBest(), AlignCenter)
		imgs = []image.Image(imgs2)
	}
	return BuildGifAnimationSimple(src, delay, imgs, f), nil
}

// BuildGifAnimationSimple assembles a set of images in an animated GIF file.
// Set `delay` to `0`.
func BuildGifAnimationSimple(src *gif.GIF, delay int, imgs []image.Image, f ToPalettedFunc) *gif.GIF {
	if src == nil {
		src = &gif.GIF{}
	}
	for _, img := range imgs {
		pimg := imageToPalettedFuncWrap(img, f)
		src.Image = append(src.Image, pimg)
		src.Delay = append(src.Delay, delay)
	}
	return src
}

func imageToPalettedFuncWrap(src image.Image, f ToPalettedFunc) *image.Paletted {
	if v, ok := src.(*image.Paletted); ok {
		return v
	}
	if f != nil {
		return f(src)
	}
	return ImageToPalettedPlan9(src)
}
