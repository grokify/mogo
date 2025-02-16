package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/jessevdk/go-flags"
)

func main() {
	cli := httpsimple.CLI{}
	_, err := flags.Parse(&cli)
	if err != nil {
		log.Fatal(err)
	}

	fmtutil.MustPrintJSON(cli)
	err = cli.Do(context.Background(), os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DONE")
}
