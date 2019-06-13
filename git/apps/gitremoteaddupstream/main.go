package main

import (
	"fmt"
	"log"

	"github.com/grokify/gotilla/cmd/cmdutil"
	"github.com/grokify/gotilla/path/filepathutil"
	"github.com/jessevdk/go-flags"
)

type cliOptions struct {
	Parent string `short:"p" long:"parent" description:"GitHub parent user" required:"true"`
	Exec   []bool `short:"e" long:"exec" description:"execute" required:"false"`
}

func main() {
	opts := cliOptions{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	leafDir, err := filepathutil.CurLeafDir()
	if err != nil {
		log.Fatal(err)
	}

	gitCmd := fmt.Sprintf("remote add upstream https://github.com/%s/%s.git", opts.Parent, leafDir)
	fmt.Printf("CMD: %s\n", gitCmd)
	if len(opts.Exec) > 0 {
		_, _, err := cmdutil.ExecSimple(gitCmd)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Executed")
	}

	fmt.Println("DONE")
}
