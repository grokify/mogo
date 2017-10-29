package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/gotilla/time/timeutil"
)

// http://choly.ca/post/go-json-marshalling/
// https://stackoverflow.com/questions/25087960/json-unmarshal-time-that-isnt-in-rfc-3339-format
// https://stackoverflow.com/questions/23695479/format-timestamp-in-outgoing-json-in-golang
// https://blog.charmes.net/post/json-dates-go/#custom-type

type Data struct {
	MyTime timeutil.RFC3339YMDTime
}

func main() {
	jsonStr := `{"MyTime":"2001-02-03"}`
	data := Data{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data.MyTime.String())
	fmtutil.PrintJSON(data)

	_ = json.NewEncoder(os.Stdout).Encode(data)
}
