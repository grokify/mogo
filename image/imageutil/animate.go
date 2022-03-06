package imageutil

import (
	"fmt"
	"image"
	"image/gif"
)

func BuildGifAnimationSimpleReadAny(src *gif.GIF, delay int, names []string, f ToPaletted) (*gif.GIF, error) {
	if src == nil {
		src = &gif.GIF{}
	}
	for i, name := range names {
		img, _, err := ReadImage(name)
		if err != nil {
			return src, fmt.Errorf("image index [%d], cannot be read at location [%s]", i, name)
		}
		pimg, err := imageToPalettedFuncWrap(img, f)
		if err != nil {
			return src, fmt.Errorf("image index [%d], cannot be converted to `*image.Paletted`", i)
		}
		src.Image = append(src.Image, pimg)
		src.Delay = append(src.Delay, delay)
	}
	return src, nil
}

func BuildGifAnimationSimple(src *gif.GIF, delay int, imgs []image.Image, f ToPaletted) (*gif.GIF, error) {
	if src == nil {
		src = &gif.GIF{}
	}
	for i, img := range imgs {
		pimg, err := imageToPalettedFuncWrap(img, f)
		if err != nil {
			return src, fmt.Errorf("image index [%d], cannot be converted to `*image.Paletted`", i)
		}
		src.Image = append(src.Image, pimg)
		src.Delay = append(src.Delay, delay)
	}
	return src, nil
}

func imageToPalettedFuncWrap(src image.Image, f func(s image.Image) *image.Paletted) (*image.Paletted, error) {
	if f != nil {
		return f(src), nil
	}
	if v, ok := src.(*image.Paletted); ok {
		return v, nil
	}
	return nil, fmt.Errorf("image cannot be converted to `*image.Paletted`, no function and is not `*image.Paletted`")
}
