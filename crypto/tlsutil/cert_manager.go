package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"sync"

	"github.com/grokify/mogo/crypto/x509util"
)

type CertManager struct {
	sync.RWMutex
	clientCert                 tls.Certificate
	clientCACertPool           *x509.CertPool
	rootCACertPool             *x509.CertPool
	requireAndVerifyClientCert bool
}

func NewCertManager(certPath, keyPath string, rootCAPaths, clientCAPaths []string, requireAndVerifyClientCert bool) (*CertManager, error) {
	cm := &CertManager{
		requireAndVerifyClientCert: requireAndVerifyClientCert}
	if err := cm.reload(certPath, keyPath, rootCAPaths, clientCAPaths); err != nil {
		return nil, err
	} else {
		return cm, nil
	}
}

func (cm *CertManager) reload(certPath, keyPath string, rootCAPaths, clientCAPaths []string) error {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}

	rootCAPool, err := x509util.NewCertPoolWithFilepaths(rootCAPaths)
	if err != nil {
		return err
	}

	clientCAPool, err := x509util.NewCertPoolWithFilepaths(clientCAPaths)
	if err != nil {
		return err
	}

	cm.Lock()
	defer cm.Unlock()
	cm.clientCert = cert
	cm.clientCACertPool = clientCAPool
	cm.rootCACertPool = rootCAPool
	return nil
}

func (cm *CertManager) TLSConfig() *tls.Config {
	cm.RLock()
	defer cm.RUnlock()
	cfg := &tls.Config{
		ClientCAs:    cm.clientCACertPool,
		RootCAs:      cm.rootCACertPool,
		Certificates: []tls.Certificate{cm.clientCert},
	}
	if cm.requireAndVerifyClientCert {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	} else {
		cfg.ClientAuth = tls.NoClientCert
	}
	return cfg
}
