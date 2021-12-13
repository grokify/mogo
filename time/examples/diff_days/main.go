package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tu "github.com/grokify/mogo/time/timeutil"
)

func main() {
	t1raw := os.Args[1]
	t2raw := os.Args[2]

	t1, err := time.Parse(tu.RFC3339FullDate, t1raw)
	if err != nil {
		log.Fatal(err)
	}
	t2, err := time.Parse(tu.RFC3339FullDate, t2raw)
	if err != nil {
		log.Fatal(err)
	}

	t1, t2 = tu.MinMax(t1, t2)
	dur := t2.Sub(t1)
	hours := dur.Hours()
	days := hours / 24.0
	fmt.Printf("DAYS: %v\n", days)
	weeks := days / 7.0
	fmt.Printf("WEEKS: %v\n", weeks)

	fmt.Println("DONE")
}
