package main

import (
	"fmt"
	"log"

	"github.com/grokify/mogo/image/imageutil"
	"github.com/grokify/mogo/log/logutil"
	flags "github.com/jessevdk/go-flags"
)

type cliOptions struct {
	Input    string `short:"i" long:"input dir/file" description:"A dir or file" value-name:"FILE" required:"true"`
	Output   string `short:"o" long:"output dir/filefile" description:"A dir or file" required:"true"`
	Quality  int    `short:"q" long:"quality" description:"Quality"`
	WidthMin int    `short:"2" long:"width-minimum" description:"Resize to minimum width"`
}

func main() {
	opts := cliOptions{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	img, format, err := imageutil.ReadImageHTTP(opts.Input)
	logutil.FatalErr(err)
	fmt.Printf("GOT: [%s]\n", format)
	if opts.WidthMin > 0 {
		img = imageutil.ResizeMin(uint(opts.WidthMin), 0, img, imageutil.ScalerBest())
	}

	im := imageutil.Image{Image: img}
	err = im.WriteJPEGFileSimple(opts.Output, opts.Quality)
	logutil.FatalErr(err)
	fmt.Printf("Wrote [%s]\n", opts.Output)
	fmt.Println("DONE")
}
