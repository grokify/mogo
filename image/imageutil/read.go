package imageutil

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
	"golang.org/x/image/webp" // "github.com/chai2010/webp"
)

const (
	FileExtensionWebp = ".webp"

	FormatNameJPG  = "jpeg"
	FormatNamePNG  = "png"
	FormatNameWEBP = "webp"
)

var RxFileExtensionJPG = regexp.MustCompile(`(?i)^.*\.*jpe?g$`)

func ReadImage(location string) (image.Image, string, error) {
	if urlutil.IsHTTP(location, true, true) {
		return ReadImageHTTP(location)
	}
	return ReadImageFile(location)
}

func ReadImages(locations []string) ([]image.Image, error) {
	images := []image.Image{}
	for _, location := range locations {
		img, _, err := ReadImage(location)
		if err != nil {
			return images, err
		}
		images = append(images, img)
	}
	return images, nil
}

func ReadImageFile(filename string) (image.Image, string, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer infile.Close()
	if strings.ToLower(strings.TrimSpace(filepath.Ext(filename))) == FileExtensionWebp {
		return DecodeWebP(infile)
	} else {
		return image.Decode(infile)
	}
}

func DecodeWebP(r io.Reader) (image.Image, string, error) {
	if img, err := webp.Decode(r); err != nil {
		return nil, "", err
	} else {
		return img, FormatNameWEBP, nil
	}
}

func ReadImageHTTP(imageURL string) (image.Image, string, error) {
	imageURL = strings.TrimSpace(imageURL)
	if !urlutil.IsHTTP(imageURL, true, true) {
		return nil, "", errors.New("url is not valid")
	}
	resp, err := http.Get(imageURL) // #nosec G107
	if err != nil {
		return nil, "", err
	} else if resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("HTTP_STATUS_CODE_GTE_300 [%v]", resp.StatusCode)
	} else if httputilmore.ResponseIsContentType(httputilmore.ContentTypeImageWebP, resp) {
		return DecodeWebP(resp.Body)
	} else {
		return image.Decode(resp.Body)
	}
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

// DecodeBytes wraps Decode which decodes an image that has been encoded in a registered format. The string returned is the format name used during format registration. Format registration is typically done by an init function in the codec- specific package.
func DecodeBytes(data []byte) (image.Image, string, error) {
	return image.Decode(bytes.NewReader(data))
}
