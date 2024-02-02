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

func NewTLSConfig(certFilepath, keyFilepath string, rootCACertFilepaths, clientCACertFilepaths []string, requireAndVerifyClientCert bool) (*TLSConfig, error) {
	cfg := &tls.Config{
		Certificates: []tls.Certificate{},
		MinVersion:   tls.VersionTLS12}
	if requireAndVerifyClientCert {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}
	tc := &TLSConfig{Config: cfg}

	if certFilepath != "" || keyFilepath != "" {
		if err := tc.LoadX509KeyPair(certFilepath, keyFilepath); err != nil {
			return nil, err
		}
	}
	for _, rootCaCertFilepath := range rootCACertFilepaths {
		if err := tc.LoadRootCACert(rootCaCertFilepath); err != nil {
			return nil, err
		}
	}
	for _, clientCACertFilepath := range clientCACertFilepaths {
		if err := tc.LoadClientCACert(clientCACertFilepath); err != nil {
			return nil, err
		}
	}
	return &TLSConfig{Config: cfg}, nil
}

/*
tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},        // server certificate which is validated by the client
		ClientCAs:    caCertPool,                     // used to verify the client cert is signed by the CA and is therefore valid
		ClientAuth:   tls.RequireAndVerifyClientCert, // this requires a valid client certificate to be supplied during handshake
	}

*/

func (tc *TLSConfig) LoadX509KeyPair(certFilepath, keyFilepath string) error {
	if cert, err := tls.LoadX509KeyPair(certFilepath, keyFilepath); err != nil {
		return err
	} else {
		if tc.Config.Certificates == nil {
			tc.Config.Certificates = []tls.Certificate{}
		}
		tc.Config.Certificates = append(tc.Config.Certificates, cert)
		return nil
	}
}

func (tc *TLSConfig) LoadClientCACert(caCertFilepath string) error {
	cert, err := os.ReadFile(caCertFilepath)
	if err != nil {
		return err
	}
	if tc.Config.ClientCAs == nil {
		tc.Config.ClientCAs = x509.NewCertPool()
	}

	if ok := tc.Config.ClientCAs.AppendCertsFromPEM(cert); !ok {
		return fmt.Errorf("cannot add client CA cert (%s)", caCertFilepath)
	} else {
		return nil
	}
}

func (tc *TLSConfig) LoadRootCACert(caCertFilepath string) error {
	cert, err := os.ReadFile(caCertFilepath)
	if err != nil {
		return err
	}
	if tc.Config.RootCAs == nil {
		tc.Config.RootCAs = x509.NewCertPool()
	}

	if ok := tc.Config.RootCAs.AppendCertsFromPEM(cert); !ok {
		return fmt.Errorf("cannot add root CA cert (%s)", caCertFilepath)
	} else {
		return nil
	}
}
