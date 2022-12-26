package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grokify/mogo/text/markdown"
)

func main() {
	slides := markdown.PresentationData{
		Slides: []markdown.RemarkSlideData{
			{
				Layout:   "middle, center, inverse",
				Class:    "false",
				Markdown: "# Test Slide\n\nTest Remark Slide",
			},
			{
				Markdown: "# Test Slide\n\nTest Remark Slide",
			},
		},
	}
	html := markdown.RemarkHTML(slides)
	fmt.Println(html)

	err := os.WriteFile("test_slides.html", []byte(html), 0600)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
