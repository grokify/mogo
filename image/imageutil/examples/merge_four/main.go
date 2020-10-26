package main

import (
	"fmt"
	"log"

	"github.com/grokify/gotilla/image/imageutil"
)

func main() {
	urlsmatrix := [][]string{
		[]string{"https://raw.githubusercontent.com/grokify/gotilla/master/image/imageutil/read_testdata/gopher_appengine_color.jpg"},
		[]string{
			"https://raw.githubusercontent.com/grokify/gotilla/master/image/imageutil/read_testdata/gopher_color.jpg",
			"https://raw.githubusercontent.com/grokify/gotilla/master/image/imageutil/read_testdata/gopher_color.jpg"},
	}

	outfile := "_merged.jpg"

	matrix, err := imageutil.MatrixRead(urlsmatrix)
	if err != nil {
		log.Fatal(err)
	}
	merged := matrix.Merge(true, true)

	err = imageutil.WriteFileJPEG(outfile, merged, imageutil.JPEGQualityMax)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", outfile)

	fmt.Println("DONE")
}
