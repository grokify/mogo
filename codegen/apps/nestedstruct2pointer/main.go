package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/grokify/mogo/codegen"
	"github.com/grokify/mogo/fmt/fmtutil"
	flags "github.com/jessevdk/go-flags"
)

/*

go get github.com/grokify/mogo/codegen/apps/nestedstruct2pointer

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
	err = fmtutil.PrintJSON(opts)
	if err != nil {
		log.Fatal(err)
	}
	if len(opts.Pattern) > 0 {
		files, err := codegen.ConvertFilesInPlaceNestedstructsToPointers(
			opts.Dir, regexp.MustCompile(opts.Pattern))
		if err != nil {
			log.Fatal(err)
		}
		err = fmtutil.PrintJSON(files)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files, err := codegen.ConvertFilesInPlaceNestedstructsToPointers(
			opts.Dir, nil)
		if err != nil {
			log.Fatal(err)
		}
		err = fmtutil.PrintJSON(files)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("DONE")
}
