package main

import (
	"fmt"
	"log"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/os/osutil"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Dir     string   `short:"d" long:"dir" description:"Directory"`
	Columns []string `short:"c" long:"columns" description:"Columns"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	ems, err := osutil.ReadDirFiles(opts.Dir, true, true, true)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := ems.Rows(false, true, opts.Columns...)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(rows)

	fmt.Println("DONE")
}
