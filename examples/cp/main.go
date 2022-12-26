package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grokify/mogo/os/osutil"
)

func main() {
	fmt.Println(len(os.Args))
	if len(os.Args) != 3 {
		log.Fatal("Needs 2 arguments")
	}
	src := os.Args[1]
	dst := os.Args[2]

	err := osutil.CopyFile(src, dst)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
