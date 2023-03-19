package main

import (
	"image/jpeg"

	"github.com/grokify/mogo/image/imageutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/jessevdk/go-flags"
)

type cliOptions struct {
	Input  string `short:"i" long:"input file" description:"A file" value-name:"FILE" required:"true"`
	Output string `short:"o" long:"output file" description:"A file" required:"true"`
	Height uint   `short:"h" long:"height" description:"Height"`
	Width  uint   `short:"w" long:"width" description:"width"`
}

func main() {
	opts := cliOptions{}
	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)

	err = imageutil.ResizeFileJPEG(opts.Input, opts.Output, opts.Width, opts.Height,
		&imageutil.JPEGEncodeOptions{
			Options: &jpeg.Options{Quality: imageutil.JPEGQualityMax},
		})
	logutil.FatalErr(err)
}
