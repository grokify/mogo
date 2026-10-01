package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/grokify/mogo/fmt/fmtutil"
)

func main() {
	qryStr := "file=https%3A%2F%2Fuploads.example.com%2Fattachments%2F1234567890abcdefdeadbeefs.pdf"

	v, err := url.ParseQuery(qryStr)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.MustPrintJSON(v)
	fmt.Println("DONE")
}
