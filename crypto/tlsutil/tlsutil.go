package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type TLSConfig struct {
	Config *tls.Config
}

func NewTLSConfig() TLSConfig {
	return TLSConfig{
		Config: &tls.Config{
			Certificates: []tls.Certificate{},
			MinVersion:   tls.VersionTLS12,
		},
	}
}

func (tc *TLSConfig) LoadX509KeyPair(certFilepath, keyFilepath string) error {
	cert, err := tls.LoadX509KeyPair(certFilepath, keyFilepath)
	if err != nil {
		return err
	}
	tc.Config.Certificates = append(tc.Config.Certificates, cert)
	return nil
}

func (tc *TLSConfig) LoadCACert(caCertFilepath string) error {
	cert, err := os.ReadFile(caCertFilepath)
	if err != nil {
		return err
	}
	if tc.Config.RootCAs == nil {
		tc.Config.RootCAs = x509.NewCertPool()
	}

	ok := tc.Config.RootCAs.AppendCertsFromPEM(cert)
	if !ok {
		return fmt.Errorf("cannot add Root CA cert %v", caCertFilepath)
	}
	return nil
}

/*
func (tc *TLSConfig) Inflate() {
	// tc.Config.BuildNameToCertificate()
}
*/
