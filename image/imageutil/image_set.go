package imageutil

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func NewImageSetFiles(filenames []string) (ImageMetaSet, error) {
	imSet := ImageMetaSet{ImageMetas: []ImageMeta{}}
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return imSet, err
		}
		img, formatName, err := image.Decode(file)
		if err != nil {
			return imSet, err
		}
		imSet.ImageMetas = append(imSet.ImageMetas, ImageMeta{
			File:       file,
			Image:      img,
			FormatName: formatName,
		})
	}
	err := imSet.CloseFilesAll()
	return imSet, err
}

type ImageMetaSet struct {
	ImageMetas []ImageMeta
}

func (ims *ImageMetaSet) CloseFilesAll() error {
	for _, im := range ims.ImageMetas {
		if im.File != nil {
			err := im.File.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ims *ImageMetaSet) MaxX() int {
	maxX := 0
	for _, im := range ims.ImageMetas {
		if im.Image.Bounds().Dx() > maxX {
			maxX = im.Image.Bounds().Dx()
		}
	}
	return maxX
}

func (ims *ImageMetaSet) MaxY() int {
	maxY := 0
	for _, im := range ims.ImageMetas {
		if im.Image.Bounds().Dy() > maxY {
			maxY = im.Image.Bounds().Dy()
		}
	}
	return maxY
}

func (ims *ImageMetaSet) SumX(maxIndexInclusive int) int {
	sumX := 0
	for i, im := range ims.ImageMetas {
		if maxIndexInclusive >= 0 && i > maxIndexInclusive {
			break
		}
		sumX += im.Image.Bounds().Dx()
	}
	return sumX
}

func (ims *ImageMetaSet) SumY(maxIndexInclusive int) int {
	sumY := 0
	for i, im := range ims.ImageMetas {
		if maxIndexInclusive >= 0 && i > maxIndexInclusive {
			break
		}
		sumY += im.Image.Bounds().Dy()
	}
	return sumY
}

func (ims *ImageMetaSet) Stats() ImageStatsMulti {
	return ImageStatsMulti{
		MaxX: ims.MaxX(),
		MaxY: ims.MaxY(),
		SumX: ims.SumX(-1),
		SumY: ims.SumY(-1)}
}

type ImageStatsMulti struct {
	MaxX int
	SumX int
	MaxY int
	SumY int
}
