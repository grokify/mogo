package httputilmore

import (
	"crypto/sha256"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Errors returned by VerifyRFC8705CertificateBinding.
var (
	ErrClientCertificateRequired = errors.New("client certificate required")
	ErrBearerTokenMissing        = errors.New("missing bearer token")
	ErrTokenMalformed            = errors.New("malformed access token")
	ErrCnfThumbprintMissing      = errors.New("missing cnf x5t#S256 claim")
	ErrCertificateBindingInvalid = errors.New("client certificate does not match cnf x5t#S256 claim")
)

// CertificateThumbprintSHA256 returns the RFC 8705 `x5t#S256` value for a
// certificate: the base64url-encoded (unpadded) SHA-256 hash of its DER bytes.
func CertificateThumbprintSHA256(cert *x509.Certificate) string {
	sum := sha256.Sum256(cert.Raw)
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// VerifyRFC8705CertificateBinding checks that the request's mutual-TLS client
// certificate matches the `cnf.x5t#S256` confirmation claim of its JWT bearer
// access token, as defined by RFC 8705 section 3 (certificate-bound access
// tokens).
//
// The token's signature is NOT verified. Call this only after the access
// token has been authenticated, otherwise the claim can be forged. Opaque
// tokens that require introspection are not supported.
func VerifyRFC8705CertificateBinding(r *http.Request) error {
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		return ErrClientCertificateRequired
	}
	token, err := bearerToken(r)
	if err != nil {
		return err
	}
	want, err := jwtCnfThumbprintSHA256(token)
	if err != nil {
		return err
	}
	got := CertificateThumbprintSHA256(r.TLS.PeerCertificates[0])
	if subtle.ConstantTimeCompare([]byte(got), []byte(want)) != 1 {
		return ErrCertificateBindingInvalid
	}
	return nil
}

// EnforceRFC8705MTLSClientCertBinding is middleware that rejects requests
// whose mutual-TLS client certificate does not match the access token's
// `cnf.x5t#S256` claim, responding `401 Unauthorized` with an RFC 6750
// `invalid_token` challenge. It must run after access token signature
// verification; see VerifyRFC8705CertificateBinding.
func EnforceRFC8705MTLSClientCertBinding(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := VerifyRFC8705CertificateBinding(r); err != nil {
			w.Header().Set(HeaderWWWAuthenticate, `Bearer error="invalid_token"`)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// bearerToken extracts the token from an `Authorization: Bearer` header. The
// scheme is matched case-insensitively per RFC 6750 and RFC 9110.
func bearerToken(r *http.Request) (string, error) {
	scheme, token, ok := strings.Cut(r.Header.Get(HeaderAuthorization), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", ErrBearerTokenMissing
	}
	if token = strings.TrimSpace(token); token == "" {
		return "", ErrBearerTokenMissing
	}
	return token, nil
}

// jwtCnfThumbprintSHA256 returns the `cnf.x5t#S256` claim from an unverified
// compact-serialized JWT.
func jwtCnfThumbprintSHA256(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrTokenMalformed
	}
	payload, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(parts[1], "="))
	if err != nil {
		return "", errors.Join(ErrTokenMalformed, err)
	}
	var claims struct {
		Cnf struct {
			X5tS256 string `json:"x5t#S256"`
		} `json:"cnf"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", errors.Join(ErrTokenMalformed, err)
	}
	if claims.Cnf.X5tS256 == "" {
		return "", ErrCnfThumbprintMissing
	}
	return claims.Cnf.X5tS256, nil
}
