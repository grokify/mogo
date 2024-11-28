package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/grokify/mogo/crypto/tlsutil"
	"github.com/grokify/mogo/log/logutil"
)

func main() {
	url := "https://pkg.go.dev/"

	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	fmt.Printf("Checking URL: [%s]\n", url)

	resp, err := http.Get(url)
	logutil.FatalErr(err)
	tlsv, err := tlsutil.HTTPResponseTLSVersion(resp)
	logutil.FatalErr(err)
	fmt.Printf("TLS_CONNECT_DEFAULT: [%s]\n", tlsv.String())

	// Test each TLS version
	tlsVersions := []tlsutil.TLSVersion{
		tls.VersionTLS10,
		tls.VersionTLS11,
		tls.VersionTLS12,
		tls.VersionTLS13}

	for _, tlsVersion := range tlsVersions {
		err := tlsutil.SupportsTLSVersion(context.Background(), tlsVersion, url)
		if err != nil {
			fmt.Printf("%s: Not Supported: (%s)\n", tlsVersion.String(), err.Error())
		} else {
			fmt.Printf("%s: Supported\n", tlsVersion.String())
		}
	}
}
