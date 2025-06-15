package x509util

import (
	"crypto/x509"
	"fmt"
	"os"
)

func NewCertPoolWithFilepaths(certPaths []string) (*x509.CertPool, error) {
	out := x509.NewCertPool()
	for _, certPath := range certPaths {
		if certBytes, err := os.ReadFile(certPath); err != nil {
			return nil, err
		} else if !out.AppendCertsFromPEM(certBytes) {
			return nil, fmt.Errorf("cannot append certs (%s)", string(certBytes))
		}
	}
	return out, nil
}
