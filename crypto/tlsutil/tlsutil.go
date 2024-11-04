package tlsutil

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
	"golang.org/x/net/context/ctxhttp"
)

type TLSConfig struct {
	Config *tls.Config
}

func NewTLSConfig() *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{},
		MinVersion:   tls.VersionTLS12}
}

func NewTLSConfigCerts(certFilepath, keyFilepath string, rootCACertFilepaths, clientCACertFilepaths []string, requireAndVerifyClientCert bool) (*TLSConfig, error) {
	cfg := NewTLSConfig()
	if requireAndVerifyClientCert {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	} else {
		cfg.ClientAuth = tls.NoClientCert
	}
	tc := &TLSConfig{Config: cfg}

	if certFilepath != "" || keyFilepath != "" {
		if err := tc.LoadX509KeyPair(certFilepath, keyFilepath); err != nil {
			return nil, errorsutil.Wrapf(err, "err on LoadX509KeyPair (%s,%s)", certFilepath, keyFilepath)
		}
	}
	for _, rootCACertFilepath := range rootCACertFilepaths {
		if strings.TrimSpace(rootCACertFilepath) != "" {
			if err := tc.LoadRootCACert(rootCACertFilepath); err != nil {
				return nil, errorsutil.Wrap(err, "err on LoadRootCACert")
			}
		}
	}
	for _, clientCACertFilepath := range clientCACertFilepaths {
		if strings.TrimSpace(clientCACertFilepath) != "" {
			if err := tc.LoadClientCACert(clientCACertFilepath); err != nil {
				fmt.Printf("ERR 3")
				return nil, errorsutil.Wrap(err, "err on LoadClientCACert")
			}
		}
	}
	return &TLSConfig{Config: cfg}, nil
}

func (tc *TLSConfig) LoadX509KeyPair(certFilepath, keyFilepath string) error {
	if cert, err := tls.LoadX509KeyPair(certFilepath, keyFilepath); err != nil {
		return errorsutil.Wrap(err, "err in LoadServerKeyPair")
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
		return errorsutil.Wrap(err, "err in LoadClientCACert")
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
		return errorsutil.Wrap(err, "err in LoadRootCACert")
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

// SupportsTLSVersion returns an error if a connection cannot be made and a nil
// if the connection is successful.
func SupportsTLSVersion(ctx context.Context, tlsVersion TLSVersion, url string) error {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: uint16(tlsVersion),
			MaxVersion: uint16(tlsVersion),
		},
	}}

	if resp, err := ctxhttp.Head(ctx, client, url); err != nil {
		return errorsutil.Wrapf(err, "tls version not supported (%s)", tlsVersion.String())
	} else {
		defer resp.Body.Close()
		return nil
	}
}
