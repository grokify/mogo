package main

import (
	"image/jpeg"
	"log"

	flags "github.com/jessevdk/go-flags"

	"github.com/grokify/mogo/image/imageutil"
)

type cliOptions struct {
	Input   string `short:"i" long:"input dir/file" description:"A dir or file" value-name:"FILE" required:"true"`
	Output  string `short:"o" long:"output dir/filefile" description:"A dir or file" required:"true"`
	Height  uint32 `short:"h" long:"height" description:"Height"`
	Width   uint32 `short:"w" long:"width" description:"Width"`
	Quality int    `short:"q" long:"quality" description:"Quality"`
}

func main() {
	opts := cliOptions{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	err = imageutil.ResizePathJPEG(opts.Input, opts.Output, int(opts.Width), int(opts.Height),
		&imageutil.JPEGEncodeOptions{
			Options: &jpeg.Options{Quality: opts.Quality},
		})
	if err != nil {
		log.Fatal(err)
	}
}
