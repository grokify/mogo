package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/grokify/mogo/time/timeutil"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Date   int `short:"d" long:"date" description:"Date in 20060102 format" required:"true"`
	Offset int `short:"o" long:"offset" description:"Show verbose debug information" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	t, err := time.Parse(timeutil.DT8, strconv.Itoa(opts.Date))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}

	t2 := t.Add(time.Duration(int64(opts.Offset)) * time.Hour * 24)
	fmt.Printf("Original Date: %s\nOffset: %d days\nNew Date: %s\n",
		t.Format(time.DateOnly), opts.Offset, t2.Format(time.DateOnly))

	fmt.Println("DONE")
	os.Exit(0)
}
