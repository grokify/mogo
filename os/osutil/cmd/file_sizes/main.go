package main

import (
	"fmt"
	"log"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/os/osutil"
	flags "github.com/jessevdk/go-flags"
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
	err = fmtutil.PrintJSON(rows)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
