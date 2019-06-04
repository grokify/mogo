package imageutil

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"regexp"
)

const (
	rxImagesExtFormat = `\.(gif|jpg|jpeg|png)$`
)

var rxImagesExt = regexp.MustCompile(rxImagesExtFormat)

func IsImageExt(imagePath string) bool {
	return rxImagesExt.MatchString(imagePath)
}

/*
https://gist.github.com/sergiotapia/7882944
If you already have loaded an image with image.Decode(), you can also

b := img.Bounds()
imgWidth := b.Max.X
imgHeight := b.Max.Y
*/

func ReadImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return -1, -1, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return -1, -1, err
	}
	return img.Width, img.Height, nil
}
