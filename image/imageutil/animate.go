package imageutil

import (
	"image"
	"image/gif"
)

func BuildGifAnimationSimpleReadAny(src *gif.GIF, delay int, names []string, f ToPalettedFunc) (*gif.GIF, error) {
	imgs, err := ReadImages(names)
	if err != nil {
		return nil, err
	}
	return BuildGifAnimationSimple(src, delay, imgs, f), nil
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
