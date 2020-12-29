package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/io/ioutilmore"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		log.Fatal("Please enter directory")
	}

	rx := regexp.MustCompile(`^[a-z]+_(\d+)_`)

	vars, err := ioutilmore.DirEntriesNameRxVarFirsts(args[1], rx)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(vars)

	fmt.Println("DONE")
}
