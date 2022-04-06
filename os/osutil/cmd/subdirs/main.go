package main

import (
	"fmt"
	"log"

	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/os/osutil"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Directory string `short:"d" long:"directory" description:"Directory" default:"."`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	err = osutil.VisitPath(opts.Directory, nil, true, false, false, func(dir string) error {
		fmt.Println(dir)
		return nil
	})
	logutil.FatalErr(err)
	fmt.Println("DONE")
}
