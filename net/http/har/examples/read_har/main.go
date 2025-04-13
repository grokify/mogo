package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/net/http/har"
)

func main() {
	filename := "path/to/my.har"

	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	h := har.Log{}

	err = json.Unmarshal(bytes, &h)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.MustPrintJSON(h)

	for _, entry := range h.Log.Entries {
		method := entry.Request.Method
		url := entry.Request.URL
		endpoint := method + " " + url
		fmt.Printf("ENDPOINT [%v]\n", endpoint)
	}
	fmt.Println("DONE")
}
