package tlsutil

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/grokify/mogo/errors/errorsutil"
	"golang.org/x/net/context/ctxhttp"
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

// PQC Curve IDs for hybrid key exchange.
// X25519MLKEM768 is defined in Go 1.24+ as tls.X25519MLKEM768.
// For compatibility with Go 1.23, we define the constant here.
const (
	// X25519MLKEM768 is the hybrid X25519 + ML-KEM-768 key exchange.
	// This is the IANA-registered value (0x11ec = 4588).
	X25519MLKEM768 tls.CurveID = 0x11ec
)

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
	case X25519MLKEM768:
		return "X25519MLKEM768"
	default:
		return fmt.Sprintf("CurveID(%d)", id)
	}
}

// CurveIDToPQCAlgorithm returns the PQC algorithm for a curve ID, if any.
func CurveIDToPQCAlgorithm(id tls.CurveID) (PQCAlgorithm, bool) {
	switch id {
	case X25519MLKEM768:
		return PQCAlgorithmMLKEM768, true
	default:
		return "", false
	}
}

// IsPQCCurve returns true if the curve ID is a PQC or hybrid PQC curve.
func IsPQCCurve(id tls.CurveID) bool {
	_, ok := CurveIDToPQCAlgorithm(id)
	return ok
}

// PQCCurvePreferences returns curve preferences that prioritize PQC hybrid curves.
func PQCCurvePreferences() []tls.CurveID {
	return []tls.CurveID{
		X25519MLKEM768,
		tls.X25519,
		tls.CurveP256,
		tls.CurveP384,
	}
}

// CheckPQCSupport tests if a URL supports PQC key exchange.
// This requires TLS 1.3 on the server and Go 1.23+ on the client.
// The check attempts to negotiate a hybrid X25519+ML-KEM-768 key exchange.
func CheckPQCSupport(ctx context.Context, url string) PQCCheckResult {
	result := PQCCheckResult{URL: url}

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: PQCCurvePreferences(),
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := ctxhttp.Get(ctx, client, url)
	if err != nil {
		result.Error = errorsutil.Wrapf(err, "connection failed").Error()
		return result
	}
	defer resp.Body.Close()

	if resp.TLS == nil {
		result.Error = "no TLS connection state"
		return result
	}

	result.Supported = true
	result.TLSVersion = TLSVersion(resp.TLS.Version).String()

	// Note: resp.TLS.CurveID was added in Go 1.21 for TLS 1.3 connections.
	// It will be 0 for TLS 1.2 connections.
	if resp.TLS.Version >= tls.VersionTLS13 {
		// Access the negotiated curve via DidResume and other fields
		// For Go 1.23+, we can check the curve ID directly
		result.CurveID = getCurveID(resp.TLS)
		result.CurveName = CurveIDName(result.CurveID)
		result.PQCKeyExchange = IsPQCCurve(result.CurveID)
		if algo, ok := CurveIDToPQCAlgorithm(result.CurveID); ok {
			result.PQCAlgorithm = algo
		}
	}

	return result
}

// getCurveID extracts the curve ID from the TLS connection state.
// This is a helper to handle the field which was added in Go 1.21.
func getCurveID(state *tls.ConnectionState) tls.CurveID {
	// The CurveID field was added in Go 1.21 for TLS 1.3 connections.
	// We access it directly since we require Go 1.21+.
	return state.CurveID //nolint:govet // CurveID field exists since Go 1.21
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
