package main

import (
	"fmt"

	"golang.org/x/text/language"
)

func PrintTag(tag language.Tag) {
	fmt.Printf("%v\n", tag)
}

func main() {
	lang := language.English
	PrintTag(lang)

	fmt.Println("DONE")
}
