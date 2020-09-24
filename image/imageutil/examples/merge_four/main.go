package main

import (
	"fmt"
	"log"

	"github.com/grokify/gotilla/image/imageutil"
)

func main() {
	urls := []string{
		"https://example/img1.jpg",
		"https://example/img2.jpg",
		"https://example/img3.jpg",
		"https://example/img4.jpg"}

	img12, err := imageutil.MergeXSameYHttp(urls[0], urls[1], true)
	if err != nil {
		log.Fatal(err)
	}

	outfile12 := "_img12.jpg"
	err = imageutil.WriteFileJPEG(outfile12, img12, imageutil.JPEGQualityMax)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", outfile12)

	img34, err := imageutil.MergeXSameYHttp(urls[2], urls[3], true)
	if err != nil {
		log.Fatal(err)
	}

	outfile34 := "_img34.jpg"
	err = imageutil.WriteFileJPEG(outfile34, img34, imageutil.JPEGQualityMax)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", outfile34)

	img4 := imageutil.MergeYSameX(img12, img34, true)
	outfile := "_four.jpg"
	err = imageutil.WriteFileJPEG(outfile, img4, imageutil.JPEGQualityMax)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", outfile)

	outfile4 := "_img4.jpg"
	img4a, err := imageutil.Merge4Http(urls[0], urls[1], urls[2], urls[3], true)
	err = imageutil.WriteFileJPEG(outfile4, img4a, imageutil.JPEGQualityMax)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", outfile4)

	fmt.Println("DONE")
}
