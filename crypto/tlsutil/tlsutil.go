package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

type TLSConfig struct {
	Config *tls.Config
}

func NewTLSConfig() TLSConfig {
	return TLSConfig{
		Config: &tls.Config{
			Certificates: []tls.Certificate{},
		},
	}
}

func (tc *TLSConfig) LoadX509KeyPair(cert_filepath, key_filepath string) error {
	cert, err := tls.LoadX509KeyPair(cert_filepath, key_filepath)
	if err != nil {
		return err
	}
	tc.Config.Certificates = append(tc.Config.Certificates, cert)
	return nil
}

func (tc *TLSConfig) LoadCACert(ca_cert_filepath string) error {
	cert, err := ioutil.ReadFile(ca_cert_filepath)
	if err != nil {
		return err
	}
	if tc.Config.RootCAs == nil {
		tc.Config.RootCAs = x509.NewCertPool()
	}

	ok := tc.Config.RootCAs.AppendCertsFromPEM(cert)
	if !ok {
		return fmt.Errorf("Cannot add Root CA cert %v", ca_cert_filepath)
	}
	return nil
}

func (tc *TLSConfig) Inflate() {
	tc.Config.BuildNameToCertificate()
}
