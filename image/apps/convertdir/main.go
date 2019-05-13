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
	InDir      string `short:"i" long:"indir" description:"Input Directory" required:"true"`
	OutDir     string `short:"o" long:"outdir" description:"Output Directory" required:"true"`
	SrcFormat  string `short:"f" long:"format" description:"Image Format" required:"true"`
	SrcRewrite []bool `short:"r" long:"rewrite" description:"Rewrite Images"`
}

func (o *Options) Validate() error {
	o.SrcFormat = strings.ToLower(strings.TrimSpace(o.SrcFormat))
	if o.SrcFormat != "kindle" && o.SrcFormat != "pdf" {
		return fmt.Errorf("Invalid Output Format [%s]", o.SrcFormat)
	}
	return nil
}

func (o *Options) Format() convertutil.CopyType {
	format := convertutil.PDFFormat
	if o.SrcFormat == "kindle" {
		format = convertutil.KindleFormat
	}
	return format
}

func (o *Options) Rewrite() bool {
	if len(o.SrcRewrite) > 0 {
		return true
	}
	return false
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

	err = convertutil.ReformatImages(
		opts.InDir,
		opts.OutDir,
		opts.Format(),
		opts.Rewrite())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
