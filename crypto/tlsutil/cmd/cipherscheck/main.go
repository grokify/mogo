package main

import (
	"fmt"
	"os"

	"github.com/grokify/mogo/crypto/tlsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <host:port>\n", os.Args[0])
		os.Exit(1)
	}

	host := os.Args[1]
	fmt.Printf("Checking TLS support for %s\n\n", host)

	result := tlsutil.CheckHost(host)

	// Print text report
	fmt.Print(result.String())

	// Print JSON output
	fmt.Println("\n=== JSON Output ===")
	fmtutil.MustPrintJSON(result)
}
