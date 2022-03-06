package imageutil

import (
	"fmt"
	"image"
	"image/gif"
)

func BuildGifAnimationSimpleReadAny(src *gif.GIF, delay int, names []string, f ToPalettedFunc) (*gif.GIF, error) {
	if src == nil {
		src = &gif.GIF{}
	}
	for i, name := range names {
		img, _, err := ReadImage(name)
		if err != nil {
			return src, fmt.Errorf("image index [%d], cannot be read at location [%s]", i, name)
		}
		pimg := imageToPalettedFuncWrap(img, f)
		src.Image = append(src.Image, pimg)
		src.Delay = append(src.Delay, delay)
	}
	return src, nil
}

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
