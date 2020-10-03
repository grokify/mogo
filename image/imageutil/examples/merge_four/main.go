package main

import (
	"fmt"
	"image"
	"log"

	"github.com/grokify/gotilla/image/imageutil"
)

func main() {
	urls := []string{
		"https://example.com/img1.jpg",
		"https://example.com/img2.jpg",
		"https://example.com/img3.jpg",
		"https://example.com/img4.jpg"}

	writeManual := true
	writeManualIntermediate := true

	img12, err := imageutil.MergeXSameYRead([]string{urls[0], urls[1]}, true)
	if err != nil {
		log.Fatal(err)
	}
	if writeManual && writeManualIntermediate {
		outfile12 := "_img12.jpg"
		err = imageutil.WriteFileJPEG(outfile12, img12, imageutil.JPEGQualityMax)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("WROTE [%v]\n", outfile12)
	}
	img34, err := imageutil.MergeXSameYRead([]string{urls[2], urls[3]}, true)
	if err != nil {
		log.Fatal(err)
	}

	if writeManual && writeManualIntermediate {
		outfile34 := "_img34.jpg"
		err = imageutil.WriteFileJPEG(outfile34, img34, imageutil.JPEGQualityMax)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("WROTE [%v]\n", outfile34)
	}
	if writeManual {
		img4 := imageutil.MergeYSameX([]image.Image{img12, img34}, true)
		outfile := "_img4_manual.jpg"
		err = imageutil.WriteFileJPEG(outfile, img4, imageutil.JPEGQualityMax)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("WROTE [%v]\n", outfile)
	}

	outfile4 := "_img4_auto.jpg"

	img4a, err := imageutil.MatrixMergeRead(
		[][]string{
			[]string{urls[0]},
			[]string{urls[0], urls[1]},
			[]string{urls[2], urls[3], urls[2]},
		},
		true, true)

	err = imageutil.WriteFileJPEG(outfile4, img4a, imageutil.JPEGQualityMax)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", outfile4)

	fmt.Println("DONE")
}
