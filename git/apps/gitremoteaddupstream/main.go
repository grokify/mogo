package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/grokify/simplego/cmd/cmdutil"
	"github.com/grokify/simplego/path/filepathutil"
	"github.com/jessevdk/go-flags"
)

// $ go get github.com/grokify/simplego/git/apps/gitremoteaddupstream

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

	gitCmd := fmt.Sprintf("git remote add upstream https://github.com/%s/%s.git", strings.TrimSpace(opts.Parent), leafDir)
	fmt.Printf("CMD: %s\n", gitCmd)
	if len(opts.Exec) > 0 {
		_, _, err := cmdutil.ExecSimple(gitCmd)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Executed...")
		fmt.Println("Common next steps:\n$ git fetch upstream\n$ git merge upstream/master\n$ git push origin master")
	}

	fmt.Println("DONE")
}
