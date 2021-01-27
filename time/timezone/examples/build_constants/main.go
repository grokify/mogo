// main this code was sourced from Stack Overflow here:
// https://stackoverflow.com/a/40130882/1908967
package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var zoneDirs = []string{
	// Update path according to your OS
	"/usr/share/zoneinfo/",
	"/usr/share/lib/zoneinfo/",
	"/usr/lib/locale/TZ/",
}

var zoneDir string

func main() {
	for _, zoneDir = range zoneDirs {
		ReadFile("")
	}
}

func ReadFile(path string) {
	files, _ := ioutil.ReadDir(zoneDir + path)
	for _, f := range files {
		if f.Name() != strings.ToUpper(f.Name()[:1])+f.Name()[1:] {
			continue
		}
		if f.IsDir() {
			ReadFile(path + "/" + f.Name())
		} else {
			fmt.Println((path + "/" + f.Name())[1:])
		}
	}
}
