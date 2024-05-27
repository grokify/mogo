package httputilmore

import (
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/grokify/mogo/crypto/tlsutil"
)

func TransportTLSVersions(tr *http.Transport) (tlsutil.TLSVersion, tlsutil.TLSVersion, error) {
	if tr == nil {
		return tls.VersionSSL30, tls.VersionSSL30, errors.New("transport not set")
	} else if tr.TLSClientConfig == nil {
		return tls.VersionSSL30, tls.VersionSSL30, errors.New("tls config not set")
	} else {
		return tlsutil.TLSVersion(tr.TLSClientConfig.MinVersion), tlsutil.TLSVersion(tr.TLSClientConfig.MaxVersion), nil
	}
}
