package tlsutil

import (
	"crypto/tls"
)

type TLSVersion uint16

const (
	VersionTLS13 TLSVersion = tls.VersionTLS13
	VersionTLS12 TLSVersion = tls.VersionTLS12
	VersionTLS11 TLSVersion = tls.VersionTLS11
	VersionTLS10 TLSVersion = tls.VersionTLS10
	VersionSSL30 TLSVersion = tls.VersionSSL30
)

func (t TLSVersion) String() string {
	switch t {
	case VersionTLS13:
		return "TLS 1.3"
	case VersionTLS12:
		return "TLS 1.2"
	case VersionTLS11:
		return "TLS 1.1"
	case VersionTLS10:
		return "TLS 1.0"
	case VersionSSL30:
		return "SSL 3.0"
	default:
		return "Unknown"
	}
}
