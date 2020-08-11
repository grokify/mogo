package imageutil

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func ReadImageAny(s string) (image.Image, string, error) {
	if isHttpUri(s) {
		return ReadImageHttp(s)
	}
	return ReadImageFile(s)
}

func isHttpUri(s string) bool {
	try := strings.ToLower(strings.TrimSpace(s))
	if strings.Index(try, "http://") == 0 || strings.Index(try, "https://") == 0 {
		return true
	}
	return false
}

func ReadImageFile(filename string) (image.Image, string, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return image.NewRGBA(image.Rectangle{}), "", err
	}
	defer infile.Close()
	return image.Decode(infile)
}

func ReadImageHttp(imageUrl string) (image.Image, string, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return image.NewRGBA(image.Rectangle{}), "", err
	} else if resp.StatusCode >= 300 {
		return image.NewRGBA(image.Rectangle{}), "", fmt.Errorf("HTTP_STATUS_CODE_GTE_300 [%v]", resp.StatusCode)
	}
	return image.Decode(resp.Body)
}

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
