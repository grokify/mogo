package main

import (
	"fmt"
	"log"

	"github.com/grokify/gotilla/image/svgutil"
)

func main() {
	f := "ringcentral_developers_logo.svg"

	svg, err := svgutil.ReadFile(f, "", 100)
	if err != nil {
		log.Fatal(err)
	}
	ar, err := svgutil.AspectRatio(svg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SVG_ASPECT_RATIO [%s] [%v]\n", f, ar)
	fmt.Println("DONE")
}
