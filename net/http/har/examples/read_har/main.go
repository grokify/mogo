package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/net/http/har"
)

func main() {
	filename := "path/to/my.har"

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	h := har.Log{}

	err = json.Unmarshal(bytes, &h)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(h)

	for _, entry := range h.Log.Entries {
		method := entry.Request.Method
		url := entry.Request.Url
		endpoint := method + " " + url
		fmt.Printf("ENDPOINT [%v]\n", endpoint)
	}
	fmt.Println("DONE")
}
