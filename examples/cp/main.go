package main

import (
	"fmt"
	"log"
	"os"

	iom "github.com/grokify/simplego/io/ioutilmore"
)

func main() {
	fmt.Println(len(os.Args))
	if len(os.Args) != 3 {
		log.Fatal("Needs 2 arguments")
	}
	src := os.Args[1]
	dst := os.Args[2]

	err := iom.CopyFile(src, dst)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
