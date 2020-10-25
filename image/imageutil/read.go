package imageutil

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/chai2010/webp"
)

const (
	FileExtensionWebp = ".webp"

	FormatNameJPG  = "jpeg"
	FormatNamePNG  = "png"
	FormatNameWEBP = "webp"
)

func ReadImageAny(location string) (image.Image, string, error) {
	if isHttpUri(location) {
		return ReadImageHttp(location)
	}
	return ReadImageFile(location)
}

func ReadImages(locations []string) ([]image.Image, error) {
	images := []image.Image{}
	for _, location := range locations {
		img, _, err := ReadImageAny(location)
		if err != nil {
			return images, err
		}
		images = append(images, img)
	}
	return images, nil
}

func isHttpUri(location string) bool {
	try := strings.ToLower(strings.TrimSpace(location))
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
	if strings.ToLower(strings.TrimSpace(filepath.Ext(filename))) == FileExtensionWebp {
		return DecodeWebpRGBA(infile)
	}
	return image.Decode(infile)
}

func DecodeWebpRGBA(r io.Reader) (image.Image, string, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, "", err
	}
	img, err := webp.DecodeRGBA(bytes)
	if err != nil {
		return img, "", err
	}
	return img, FormatNameWEBP, nil
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
