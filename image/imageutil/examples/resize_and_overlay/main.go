package main

import (
	"fmt"
	"log"

	"github.com/grokify/mogo/image/imageutil"
	"github.com/pkg/errors"
)

// https://www.golangprograms.com/how-to-add-watermark-or-merge-two-image.html

func main() {
	fileTop := "_top.png"
	fileBkg := "_background.png"
	fileOut := "_output.png"

	imgTop, _, err := imageutil.ReadImageFile(fileTop)
	if err != nil {
		log.Fatal(errors.Wrap(err, fileTop))
	}
	imgBkg, _, err := imageutil.ReadImageFile(fileBkg)
	if err != nil {
		log.Fatal(errors.Wrap(err, fileBkg))
	}

	imgTop = imageutil.AddBackgroundWhite(imgTop)
	imgTop = imageutil.Resize(120, 0, imgTop, imageutil.ScalerBest())
	imgOut := imageutil.OverlayCenterYLeftAlign(imgBkg, imgTop)
	err = imageutil.WriteFilePNG(fileOut, imgOut)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DONE")
}
