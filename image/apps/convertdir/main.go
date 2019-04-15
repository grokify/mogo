package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/grokify/gotilla/image/convertutil"
	"github.com/jessevdk/go-flags"
)

// go get github.com/grokify/gotilla/image/apps/convertdir

type Options struct {
	InDir  string `short:"i" long:"indir" description:"Input Directory" required:"true"`
	OutDir string `short:"o" long:"outdir" description:"Output Directory" required:"true"`
	Format string `short:"f" long:"format" description:"Image Format" required:"true"`
}

func (o *Options) Validate() error {
	o.Format = strings.ToLower(strings.TrimSpace(o.Format))
	if o.Format != "kindle" && o.Format != "pdf" {
		return fmt.Errorf("Invalid Output Format [%s]", o.Format)
	}
	return nil
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	err = opts.Validate()
	if err != nil {
		log.Fatal(err)
	}
	convertType := convertutil.PDFFormat
	if opts.Format == "kindle" {
		convertType = convertutil.KindleFormat
	}

	err = convertutil.ReformatImages(
		opts.InDir,
		opts.OutDir,
		convertType)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
