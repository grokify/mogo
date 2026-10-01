package httputilmore

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestCert(t *testing.T, cn string) *x509.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return cert
}

// testJWT builds an unsigned compact JWT with the given JSON payload.
func testJWT(payload string) string {
	enc := base64.RawURLEncoding
	return enc.EncodeToString([]byte(`{"alg":"none"}`)) + "." + enc.EncodeToString([]byte(payload)) + ".sig"
}

func newMTLSRequest(cert *x509.Certificate, authz string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
	if cert != nil {
		r.TLS = &tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}
	}
	if authz != "" {
		r.Header.Set(HeaderAuthorization, authz)
	}
	return r
}

func TestVerifyRFC8705CertificateBinding(t *testing.T) {
	cert := newTestCert(t, "client")
	other := newTestCert(t, "other")
	bound := testJWT(`{"sub":"a","cnf":{"x5t#S256":"` + CertificateThumbprintSHA256(cert) + `"}}`)

	tests := []struct {
		name    string
		cert    *x509.Certificate
		authz   string
		wantErr error
	}{
		{"bound token", cert, "Bearer " + bound, nil},
		{"lowercase scheme", cert, "bearer " + bound, nil},
		{"no client cert", nil, "Bearer " + bound, ErrClientCertificateRequired},
		{"no authorization", cert, "", ErrBearerTokenMissing},
		{"basic scheme", cert, "Basic dXNlcjpwYXNz", ErrBearerTokenMissing},
		{"empty token", cert, "Bearer  ", ErrBearerTokenMissing},
		{"not a jwt", cert, "Bearer opaque-token", ErrTokenMalformed},
		{"bad payload encoding", cert, "Bearer a.!!!.c", ErrTokenMalformed},
		{"payload not json", cert, "Bearer " + testJWT(`not json`), ErrTokenMalformed},
		{"no cnf claim", cert, "Bearer " + testJWT(`{"sub":"a"}`), ErrCnfThumbprintMissing},
		{"different cert", other, "Bearer " + bound, ErrCertificateBindingInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyRFC8705CertificateBinding(newMTLSRequest(tt.cert, tt.authz))
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("VerifyRFC8705CertificateBinding() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnforceRFC8705MTLSClientCertBinding(t *testing.T) {
	cert := newTestCert(t, "client")
	bound := testJWT(`{"cnf":{"x5t#S256":"` + CertificateThumbprintSHA256(cert) + `"}}`)
	handler := EnforceRFC8705MTLSClientCertBinding(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, newMTLSRequest(cert, "Bearer "+bound))
	if rec.Code != http.StatusNoContent {
		t.Errorf("bound request: status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, newMTLSRequest(newTestCert(t, "other"), "Bearer "+bound))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("mismatched request: status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if got, want := rec.Header().Get(HeaderWWWAuthenticate), `Bearer error="invalid_token"`; got != want {
		t.Errorf("WWW-Authenticate = %q, want %q", got, want)
	}
}

func TestCertificateThumbprintSHA256(t *testing.T) {
	cert := newTestCert(t, "client")
	got := CertificateThumbprintSHA256(cert)
	if len(got) != 43 { // 32-byte SHA-256, unpadded base64url
		t.Errorf("thumbprint length = %d, want 43", len(got))
	}
	if got != CertificateThumbprintSHA256(cert) {
		t.Error("thumbprint is not deterministic")
	}
}
