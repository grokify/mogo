/*
# httpsversioncheck

## Installation:

`go install github.com/grokify/mogo/crypto/tlsutil/cmd/httpsversioncheck“

## Usage

```
% httpsversioncheck https://example.com
```
*/
package main

import (
	"fmt"
	"os"

	"github.com/grokify/mogo/crypto/tlsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("usage: httpsversioncheck <url1> <url2> <url3>")
		os.Exit(1)
	}

	res := tlsutil.CheckURLs(os.Args[1:])
	fmtutil.PrintJSON(res)
}
