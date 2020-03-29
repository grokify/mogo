package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/grokify/gotilla/codegen"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/jessevdk/go-flags"
)

/*

go get github.com/grokify/gotilla/codegen/apps/nestedstruct2pointer

*/

type Options struct {
	Dir     string `short:"d" long:"dir" description:"Directory"`
	Pattern string `short:"p" long:"pattern" description:"Pattern"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	if opts.Dir == "" {
		opts.Dir = "."
	}
	fmtutil.PrintJSON(opts)
	if len(opts.Pattern) > 0 {
		files, err := codegen.ConvertFilesInPlaceNestedstructsToPointers(
			opts.Dir, regexp.MustCompile(opts.Pattern))
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(files)
	} else {
		files, err := codegen.ConvertFilesInPlaceNestedstructsToPointers(
			opts.Dir, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(files)
	}
	fmt.Println("DONE")
}
