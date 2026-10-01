package tlsutil

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/grokify/mogo/errors/errorsutil"
)

// PQCAlgorithm represents a post-quantum cryptographic algorithm.
type PQCAlgorithm string

const (
	// PQCAlgorithmMLKEM768 is ML-KEM-768 (formerly CRYSTALS-Kyber-768).
	PQCAlgorithmMLKEM768 PQCAlgorithm = "ML-KEM-768"
	// PQCAlgorithmMLKEM1024 is ML-KEM-1024 (formerly CRYSTALS-Kyber-1024).
	PQCAlgorithmMLKEM1024 PQCAlgorithm = "ML-KEM-1024"
	// PQCAlgorithmMLDSA is ML-DSA (formerly CRYSTALS-Dilithium).
	PQCAlgorithmMLDSA PQCAlgorithm = "ML-DSA"
	// PQCAlgorithmFalcon is the Falcon signature algorithm.
	PQCAlgorithmFalcon PQCAlgorithm = "Falcon"
	// PQCAlgorithmSLHDSA is SLH-DSA (formerly SPHINCS+).
	PQCAlgorithmSLHDSA PQCAlgorithm = "SLH-DSA"
)

// PQCAlgorithmType represents the type of PQC algorithm.
type PQCAlgorithmType string

const (
	// PQCAlgorithmTypeKEM is a key encapsulation mechanism.
	PQCAlgorithmTypeKEM PQCAlgorithmType = "KEM"
	// PQCAlgorithmTypeSignature is a digital signature algorithm.
	PQCAlgorithmTypeSignature PQCAlgorithmType = "Signature"
)

// PQCAlgorithmInfo contains information about a PQC algorithm.
type PQCAlgorithmInfo struct {
	Algorithm    PQCAlgorithm     `json:"algorithm"`
	Type         PQCAlgorithmType `json:"type"`
	OriginalName string           `json:"originalName"`
	NISTLevel    int              `json:"nistLevel"`
	StdlibCheck  bool             `json:"stdlibCheck"`
}

// PQCAlgorithms returns information about known PQC algorithms.
func PQCAlgorithms() []PQCAlgorithmInfo {
	return []PQCAlgorithmInfo{
		{Algorithm: PQCAlgorithmMLKEM768, Type: PQCAlgorithmTypeKEM, OriginalName: "CRYSTALS-Kyber-768", NISTLevel: 3, StdlibCheck: true},
		{Algorithm: PQCAlgorithmMLKEM1024, Type: PQCAlgorithmTypeKEM, OriginalName: "CRYSTALS-Kyber-1024", NISTLevel: 5, StdlibCheck: true},
		{Algorithm: PQCAlgorithmMLDSA, Type: PQCAlgorithmTypeSignature, OriginalName: "CRYSTALS-Dilithium", NISTLevel: 3, StdlibCheck: false},
		{Algorithm: PQCAlgorithmFalcon, Type: PQCAlgorithmTypeSignature, OriginalName: "Falcon", NISTLevel: 5, StdlibCheck: false},
		{Algorithm: PQCAlgorithmSLHDSA, Type: PQCAlgorithmTypeSignature, OriginalName: "SPHINCS+", NISTLevel: 5, StdlibCheck: false},
	}
}

// X25519MLKEM768 is the hybrid X25519 + ML-KEM-768 key exchange.
//
// Deprecated: use tls.X25519MLKEM768.
const X25519MLKEM768 = tls.X25519MLKEM768

// curveMLKEM1024 is the pure ML-KEM-1024 key exchange (IANA 0x0202). It is
// tls.MLKEM1024 in Go 1.27+; it is defined here so results can name it while
// this module supports Go 1.26.
const curveMLKEM1024 tls.CurveID = 0x0202

// PQCCheckResult contains the result of a PQC support check.
type PQCCheckResult struct {
	URL            string       `json:"url"`
	TLSVersion     string       `json:"tlsVersion,omitempty"`
	CurveID        tls.CurveID  `json:"curveId,omitempty"`
	CurveName      string       `json:"curveName,omitempty"`
	PQCKeyExchange bool         `json:"pqcKeyExchange"`
	PQCAlgorithm   PQCAlgorithm `json:"pqcAlgorithm,omitempty"`
	Supported      bool         `json:"supported"`
	Error          string       `json:"error,omitempty"`
}

// CurveIDName returns the name of a tls.CurveID.
func CurveIDName(id tls.CurveID) string {
	switch id {
	case tls.CurveP256:
		return "P-256"
	case tls.CurveP384:
		return "P-384"
	case tls.CurveP521:
		return "P-521"
	case tls.X25519:
		return "X25519"
	case tls.X25519MLKEM768:
		return "X25519MLKEM768"
	case tls.SecP256r1MLKEM768:
		return "SecP256r1MLKEM768"
	case tls.SecP384r1MLKEM1024:
		return "SecP384r1MLKEM1024"
	case curveMLKEM1024:
		return "MLKEM1024"
	default:
		return fmt.Sprintf("CurveID(%d)", id)
	}
}

// CurveIDToPQCAlgorithm returns the PQC algorithm for a curve ID, if any.
// Hybrid curves map to their ML-KEM component.
func CurveIDToPQCAlgorithm(id tls.CurveID) (PQCAlgorithm, bool) {
	switch id {
	case tls.X25519MLKEM768, tls.SecP256r1MLKEM768:
		return PQCAlgorithmMLKEM768, true
	case tls.SecP384r1MLKEM1024, curveMLKEM1024:
		return PQCAlgorithmMLKEM1024, true
	default:
		return "", false
	}
}

// IsPQCCurve returns true if the curve ID is a PQC or hybrid PQC curve.
func IsPQCCurve(id tls.CurveID) bool {
	_, ok := CurveIDToPQCAlgorithm(id)
	return ok
}

// PQCCurvePreferences returns curve preferences that offer the hybrid PQC
// curves before classical ones.
func PQCCurvePreferences() []tls.CurveID {
	return []tls.CurveID{
		tls.X25519MLKEM768,
		tls.SecP256r1MLKEM768,
		tls.SecP384r1MLKEM1024,
		tls.X25519,
		tls.CurveP256,
		tls.CurveP384,
	}
}

// CheckPQCSupport tests if a URL supports PQC key exchange by connecting with
// TLS 1.3 and offering the hybrid curves from PQCCurvePreferences.
func CheckPQCSupport(ctx context.Context, url string) PQCCheckResult {
	return checkPQCSupport(ctx, url, nil)
}

// checkPQCSupport implements CheckPQCSupport. rootCAs overrides the system
// roots when non-nil, which lets tests use a local TLS server.
func checkPQCSupport(ctx context.Context, url string, rootCAs *x509.CertPool) (result PQCCheckResult) {
	result.URL = url

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion:       tls.VersionTLS13,
				CurvePreferences: PQCCurvePreferences(),
				RootCAs:          rootCAs,
			},
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		result.Error = errorsutil.Wrapf(err, "invalid request").Error()
		return result
	}
	resp, err := client.Do(req)
	if err != nil {
		result.Error = errorsutil.Wrapf(err, "connection failed").Error()
		return result
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			msg := errorsutil.Wrapf(err, "close response body").Error()
			if result.Error != "" {
				msg = result.Error + "; " + msg
			}
			result.Error = msg
		}
	}()

	if resp.TLS == nil {
		result.Error = "no TLS connection state"
		return result
	}

	result.Supported = true
	result.TLSVersion = TLSVersion(resp.TLS.Version).String()
	result.CurveID = resp.TLS.CurveID
	result.CurveName = CurveIDName(result.CurveID)
	result.PQCAlgorithm, result.PQCKeyExchange = CurveIDToPQCAlgorithm(result.CurveID)
	return result
}

// CheckPQCURLs checks multiple URLs for PQC support.
func CheckPQCURLs(ctx context.Context, urls []string) []PQCCheckResult {
	results := make([]PQCCheckResult, 0, len(urls))
	for _, url := range urls {
		results = append(results, CheckPQCSupport(ctx, url))
	}
	return results
}

// PQCSupportSummary provides a summary of PQC support checks.
type PQCSupportSummary struct {
	TotalChecked int              `json:"totalChecked"`
	PQCSupported int              `json:"pqcSupported"`
	TLS13Only    int              `json:"tls13Only"`
	Failed       int              `json:"failed"`
	Results      []PQCCheckResult `json:"results"`
}

// CheckPQCURLsWithSummary checks multiple URLs and returns a summary.
func CheckPQCURLsWithSummary(ctx context.Context, urls []string) PQCSupportSummary {
	results := CheckPQCURLs(ctx, urls)
	summary := PQCSupportSummary{
		TotalChecked: len(results),
		Results:      results,
	}
	for _, r := range results {
		switch {
		case r.Error != "":
			summary.Failed++
		case r.PQCKeyExchange:
			summary.PQCSupported++
		case r.TLSVersion == "TLS 1.3":
			summary.TLS13Only++
		}
	}
	return summary
}
