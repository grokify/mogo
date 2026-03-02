package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/pointer"
)

// List of weak TLS 1.2 ciphers
var weakCiphers = map[uint16]string{
	// tls.TLS_RSA_WITH_RC4_128_MD5:                "TLS_RSA_WITH_RC4_128_MD5",
	// tls.TLS_RSA_WITH_AES_256_CBC_SHA256:         "TLS_RSA_WITH_AES_256_CBC_SHA256",
	// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384:   "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384",
	// tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384: "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384",
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:    "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:    "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
	tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:        "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",    // RC4-based Ciphers (Insecure due to biases in RC4)
	tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:     "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA", // 3DES-based Ciphers (Vulnerable to Sweet32 attack)
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:      "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:      "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:   "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
	tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:          "TLS_ECDHE_RSA_WITH_RC4_128_SHA",  // RC4-based Ciphers (Insecure due to biases in RC4)
	tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:           "TLS_RSA_WITH_3DES_EDE_CBC_SHA",   // 3DES-based Ciphers (Vulnerable to Sweet32 attack)
	tls.TLS_RSA_WITH_AES_128_CBC_SHA:            "TLS_RSA_WITH_AES_128_CBC_SHA",    // RSA Key Exchange without Forward Secrecy
	tls.TLS_RSA_WITH_AES_256_CBC_SHA:            "TLS_RSA_WITH_AES_256_CBC_SHA",    // RSA Key Exchange without Forward Secrecy
	tls.TLS_RSA_WITH_AES_128_CBC_SHA256:         "TLS_RSA_WITH_AES_128_CBC_SHA256", // RSA Key Exchange without Forward Secrecy
	tls.TLS_RSA_WITH_RC4_128_SHA:                "TLS_RSA_WITH_RC4_128_SHA",        // RC4-based Ciphers (Insecure due to biases in RC4)
}

/* strong ciphers
TLS_RSA_WITH_AES_128_GCM_SHA256 (0x009c)
TLS_RSA_WITH_AES_256_GCM_SHA384 (0x009d)
TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 (0xc02f)
TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 (0xc02b)
TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 (0xc030)
TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384 (0xc02c)
TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256 (0xcca8)
TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256 (0xcca9)
*/

type CipherSuites struct {
	Host   string
	Time   *time.Time
	Strong []string
	Weak   []string
}

func checkAllCiphers(host string) {
	fmt.Printf("Checking supported ciphers for %s\n", host)

	res := CipherSuites{
		Host: host,
		Time: pointer.Pointer(time.Now().UTC()),
	}

	for _, cipher := range tls.CipherSuites() {
		config := &tls.Config{
			InsecureSkipVerify: true, // #nosec G402 -- intentional for cipher testing tool
			CipherSuites:       []uint16{cipher.ID},
		}
		conn, err := tls.Dial("tcp", host, config)
		if err == nil {
			conn.Close()
			if _, found := weakCiphers[cipher.ID]; found {
				fmt.Printf("WEAK: %s\n", cipher.Name)
				res.Weak = append(res.Weak, cipher.Name)
			} else {
				fmt.Printf("STRONG: %s\n", cipher.Name)
				res.Strong = append(res.Strong, cipher.Name)
			}
		}
	}
	if len(res.Strong) > 0 {
		sort.Strings(res.Strong)
	}
	if len(res.Weak) > 0 {
		sort.Strings(res.Weak)
	}
	fmtutil.MustPrintJSON(res)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <host:port>")
		os.Exit(1)
	}

	host := os.Args[1]
	checkAllCiphers(host)
}
