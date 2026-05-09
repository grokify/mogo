package tlsutil

import (
	"crypto/tls"
	"encoding/asn1"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/grokify/mogo/pointer"
)

// PQSignatureOIDs maps known post-quantum signature algorithm OIDs to names.
// These are NIST standardized and draft OIDs for PQ signature algorithms.
var PQSignatureOIDs = map[string]string{
	// ML-DSA (FIPS 204) - formerly Dilithium
	// NIST OID arc: 2.16.840.1.101.3.4.3.x
	"2.16.840.1.101.3.4.3.17": "ML-DSA-44",
	"2.16.840.1.101.3.4.3.18": "ML-DSA-65",
	"2.16.840.1.101.3.4.3.19": "ML-DSA-87",

	// SLH-DSA (FIPS 205) - formerly SPHINCS+
	// NIST OID arc for SLH-DSA
	"2.16.840.1.101.3.4.3.20": "SLH-DSA-SHA2-128s",
	"2.16.840.1.101.3.4.3.21": "SLH-DSA-SHA2-128f",
	"2.16.840.1.101.3.4.3.22": "SLH-DSA-SHA2-192s",
	"2.16.840.1.101.3.4.3.23": "SLH-DSA-SHA2-192f",
	"2.16.840.1.101.3.4.3.24": "SLH-DSA-SHA2-256s",
	"2.16.840.1.101.3.4.3.25": "SLH-DSA-SHA2-256f",
	"2.16.840.1.101.3.4.3.26": "SLH-DSA-SHAKE-128s",
	"2.16.840.1.101.3.4.3.27": "SLH-DSA-SHAKE-128f",
	"2.16.840.1.101.3.4.3.28": "SLH-DSA-SHAKE-192s",
	"2.16.840.1.101.3.4.3.29": "SLH-DSA-SHAKE-192f",
	"2.16.840.1.101.3.4.3.30": "SLH-DSA-SHAKE-256s",
	"2.16.840.1.101.3.4.3.31": "SLH-DSA-SHAKE-256f",

	// Experimental/draft OIDs that may be encountered
	// OQS (Open Quantum Safe) experimental OIDs
	"1.3.6.1.4.1.2.267.7.4.4":  "Dilithium2 (draft)",
	"1.3.6.1.4.1.2.267.7.6.5":  "Dilithium3 (draft)",
	"1.3.6.1.4.1.2.267.7.8.7":  "Dilithium5 (draft)",
	"1.3.6.1.4.1.2.267.12.4.4": "Falcon-512 (draft)",
	"1.3.6.1.4.1.2.267.12.9.9": "Falcon-1024 (draft)",

	// SPHINCS+ experimental OIDs
	"1.3.9999.6.4.1":  "SPHINCS+-SHA256-128f-robust (draft)",
	"1.3.9999.6.4.4":  "SPHINCS+-SHA256-128s-robust (draft)",
	"1.3.9999.6.4.10": "SPHINCS+-SHA256-192f-robust (draft)",
	"1.3.9999.6.4.13": "SPHINCS+-SHA256-192s-robust (draft)",
	"1.3.9999.6.5.1":  "SPHINCS+-SHA256-256f-robust (draft)",
	"1.3.9999.6.5.4":  "SPHINCS+-SHA256-256s-robust (draft)",
}

// IsPQSignatureOID returns true if the OID is a known post-quantum signature algorithm.
func IsPQSignatureOID(oid asn1.ObjectIdentifier) bool {
	_, ok := PQSignatureOIDs[oid.String()]
	return ok
}

// PQSignatureOIDName returns the name for a PQ signature OID, or empty string if not PQ.
func PQSignatureOIDName(oid asn1.ObjectIdentifier) string {
	return PQSignatureOIDs[oid.String()]
}

// TLSVersionStatus represents the support status of a TLS version.
type TLSVersionStatus struct {
	Version TLSVersion `json:"-"`
	Name    string     `json:"version"`
	Active  bool       `json:"active"`
}

// PostQuantumSupport represents post-quantum cryptography support.
type PostQuantumSupport struct {
	KeyExchange []string `json:"keyExchange"`
	Signatures  []string `json:"signatures"`
}

// CipherCheckResult contains comprehensive TLS cipher check results.
type CipherCheckResult struct {
	Host        string             `json:"host"`
	Time        *time.Time         `json:"time"`
	TLSVersions []TLSVersionStatus `json:"tlsVersions"`
	Strong      []string           `json:"strong"`
	Weak        []string           `json:"weak"`
	Insecure    []string           `json:"insecure"`
	TLS13       []string           `json:"tls13"`
	PostQuantum PostQuantumSupport `json:"postQuantum"`
}

// String returns a formatted text report of the check results.
func (r CipherCheckResult) String() string {
	var sb strings.Builder

	// TLS version support
	sb.WriteString("=== TLS Version Support ===\n")
	for _, v := range r.TLSVersions {
		status := "INACTIVE"
		if v.Active {
			status = "ACTIVE"
		}
		sb.WriteString(fmt.Sprintf("%-8s: %s\n", v.Name, status))
	}

	// TLS 1.2 cipher suites
	sb.WriteString("\n=== TLS 1.2 Cipher Suites ===\n")
	for _, c := range r.Strong {
		sb.WriteString(fmt.Sprintf("STRONG: %s\n", c))
	}
	for _, c := range r.Weak {
		sb.WriteString(fmt.Sprintf("WEAK: %s\n", c))
	}
	for _, c := range r.Insecure {
		sb.WriteString(fmt.Sprintf("INSECURE: %s\n", c))
	}

	// TLS 1.3 cipher suites
	sb.WriteString("\n=== TLS 1.3 Cipher Suites ===\n")
	for _, c := range r.TLS13 {
		sb.WriteString(fmt.Sprintf("TLS 1.3: %s\n", c))
	}

	// Post-quantum support
	sb.WriteString("\n=== Post-Quantum Support ===\n")
	sb.WriteString("Key Exchange:\n")
	if len(r.PostQuantum.KeyExchange) > 0 {
		for _, k := range r.PostQuantum.KeyExchange {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	} else {
		sb.WriteString("  none detected\n")
	}
	sb.WriteString("Signatures (certificate chain):\n")
	if len(r.PostQuantum.Signatures) > 0 {
		for _, s := range r.PostQuantum.Signatures {
			sb.WriteString(fmt.Sprintf("  %s\n", s))
		}
	} else {
		sb.WriteString("  none detected (classical signatures)\n")
	}

	return sb.String()
}

// NewCipherCheckResult creates a new CipherCheckResult with initialized slices.
func NewCipherCheckResult(host string) CipherCheckResult {
	return CipherCheckResult{
		Host:        host,
		Time:        pointer.Pointer(time.Now().UTC()),
		TLSVersions: []TLSVersionStatus{},
		Strong:      []string{},
		Weak:        []string{},
		Insecure:    []string{},
		TLS13:       []string{},
		PostQuantum: PostQuantumSupport{
			KeyExchange: []string{},
			Signatures:  []string{},
		},
	}
}

// CheckerConfig contains configuration for the TLS checker.
type CheckerConfig struct {
	Timeout time.Duration
}

// DefaultCheckerConfig returns the default checker configuration.
func DefaultCheckerConfig() CheckerConfig {
	return CheckerConfig{
		Timeout: 10 * time.Second,
	}
}

// Checker performs TLS checks against a host.
type Checker struct {
	config CheckerConfig
	dialer *net.Dialer
}

// NewChecker creates a new TLS checker with the given configuration.
func NewChecker(config CheckerConfig) *Checker {
	return &Checker{
		config: config,
		dialer: &net.Dialer{Timeout: config.Timeout},
	}
}

// CheckTLSVersion tests if a host supports a specific TLS version.
func (c *Checker) CheckTLSVersion(host string, version TLSVersion) bool {
	config := &tls.Config{
		InsecureSkipVerify: true, // #nosec G402 -- intentional for TLS version testing
		MinVersion:         uint16(version),
		MaxVersion:         uint16(version),
	}
	conn, err := tls.DialWithDialer(c.dialer, "tcp", host, config)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// CheckAllTLSVersions tests all TLS versions (1.0, 1.1, 1.2, 1.3).
func (c *Checker) CheckAllTLSVersions(host string) []TLSVersionStatus {
	versions := TLSVersions()
	results := make([]TLSVersionStatus, 0, len(versions))
	for _, v := range versions {
		results = append(results, TLSVersionStatus{
			Version: v,
			Name:    v.String(),
			Active:  c.CheckTLSVersion(host, v),
		})
	}
	return results
}

// CheckTLS12Cipher tests if a host supports a specific TLS 1.2 cipher suite.
func (c *Checker) CheckTLS12Cipher(host string, cipherID uint16) bool {
	config := &tls.Config{
		InsecureSkipVerify: true, // #nosec G402 -- intentional for cipher testing
		CipherSuites:       []uint16{cipherID},
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
	}
	conn, err := tls.DialWithDialer(c.dialer, "tcp", host, config)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// CheckTLS12Ciphers tests all TLS 1.2 cipher suites and categorizes them.
func (c *Checker) CheckTLS12Ciphers(host string) (strong, weak, insecure []string) {
	weakMap := TLS12CiphersWeakMap()

	// Check standard cipher suites
	for _, cipher := range tls.CipherSuites() {
		if c.CheckTLS12Cipher(host, cipher.ID) {
			if _, isWeak := weakMap[cipher.ID]; isWeak {
				weak = append(weak, cipher.Name)
			} else {
				strong = append(strong, cipher.Name)
			}
		}
	}

	// Check insecure cipher suites
	for _, cipher := range tls.InsecureCipherSuites() {
		config := &tls.Config{
			InsecureSkipVerify: true, // #nosec G402 -- intentional for cipher testing
			CipherSuites:       []uint16{cipher.ID},
			MinVersion:         tls.VersionTLS10,
			MaxVersion:         tls.VersionTLS12,
		}
		conn, err := tls.DialWithDialer(c.dialer, "tcp", host, config)
		if err == nil {
			conn.Close()
			insecure = append(insecure, cipher.Name)
		}
	}

	sort.Strings(strong)
	sort.Strings(weak)
	sort.Strings(insecure)
	return strong, weak, insecure
}

// CheckTLS13Ciphers tests TLS 1.3 cipher suites by connecting and checking what's negotiated.
func (c *Checker) CheckTLS13Ciphers(host string) []string {
	var ciphers []string

	// TLS 1.3 cipher suites are automatically negotiated; we can't force specific ones.
	// Connect and see what cipher suite is negotiated.
	config := &tls.Config{
		InsecureSkipVerify: true, // #nosec G402 -- intentional for cipher testing
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
	}
	conn, err := tls.DialWithDialer(c.dialer, "tcp", host, config)
	if err != nil {
		return ciphers
	}
	state := conn.ConnectionState()
	conn.Close()

	// Report the negotiated cipher suite
	cipherName := tls.CipherSuiteName(state.CipherSuite)
	ciphers = append(ciphers, cipherName)

	// Also list other TLS 1.3 cipher suites available in Go
	// (the server likely supports these too if it supports TLS 1.3)
	for _, cipher := range tls.CipherSuites() {
		if isTLS13Cipher(cipher) && cipher.Name != cipherName {
			ciphers = append(ciphers, cipher.Name+" (available)")
		}
	}

	sort.Strings(ciphers)
	return ciphers
}

// isTLS13Cipher returns true if the cipher suite supports TLS 1.3.
func isTLS13Cipher(cipher *tls.CipherSuite) bool {
	for _, v := range cipher.SupportedVersions {
		if v == tls.VersionTLS13 {
			return true
		}
	}
	return false
}

// CheckPostQuantumKeyExchange tests post-quantum key exchange support.
func (c *Checker) CheckPostQuantumKeyExchange(host string) []string {
	algorithms := []string{} // Initialize to empty slice, not nil

	// Test X25519MLKEM768 (hybrid X25519 + ML-KEM-768)
	config := &tls.Config{
		InsecureSkipVerify: true, // #nosec G402 -- intentional for PQ testing
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
		CurvePreferences:   []tls.CurveID{X25519MLKEM768},
	}
	conn, err := tls.DialWithDialer(c.dialer, "tcp", host, config)
	if err == nil {
		state := conn.ConnectionState()
		conn.Close()
		if state.CurveID == X25519MLKEM768 {
			algorithms = append(algorithms, CurveIDName(X25519MLKEM768))
		}
	}

	return algorithms
}

// CheckPostQuantumSignatures tests post-quantum signature algorithms in server certificates.
// It connects to the server, retrieves the certificate chain, and checks signature algorithm OIDs.
func (c *Checker) CheckPostQuantumSignatures(host string) []string {
	signatures := []string{} // Initialize to empty slice, not nil
	seen := make(map[string]bool)

	// Extract hostname for SNI
	hostname, _, err := extractHostPort(host)
	if err != nil {
		hostname = host
	}

	config := &tls.Config{
		InsecureSkipVerify: true, // #nosec G402 -- intentional for cert inspection
		ServerName:         hostname,
	}

	conn, err := tls.DialWithDialer(c.dialer, "tcp", host, config)
	if err != nil {
		return signatures
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return signatures
	}

	// Check each certificate in the chain
	for _, cert := range state.PeerCertificates {
		// Extract signature algorithm OID from raw certificate ASN.1
		oid, err := extractSignatureAlgorithmOID(cert.Raw)
		if err != nil {
			continue
		}

		// Check if this is a PQ signature algorithm
		oidStr := oid.String()
		if name, ok := PQSignatureOIDs[oidStr]; ok && !seen[name] {
			seen[name] = true
			signatures = append(signatures, name)
		}
	}

	sort.Strings(signatures)
	return signatures
}

// certificateASN1 represents the top-level X.509 certificate structure.
type certificateASN1 struct {
	TBSCertificate     asn1.RawValue
	SignatureAlgorithm algorithmIdentifier
	SignatureValue     asn1.BitString
}

// algorithmIdentifier represents an X.509 AlgorithmIdentifier.
type algorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}

// extractSignatureAlgorithmOID extracts the signature algorithm OID from a DER-encoded certificate.
func extractSignatureAlgorithmOID(der []byte) (asn1.ObjectIdentifier, error) {
	var cert certificateASN1
	_, err := asn1.Unmarshal(der, &cert)
	if err != nil {
		return nil, err
	}
	return cert.SignatureAlgorithm.Algorithm, nil
}

// extractHostPort splits a host:port string.
func extractHostPort(addr string) (host, port string, err error) {
	host, port, err = net.SplitHostPort(addr)
	return
}

// CheckAll performs a comprehensive TLS check on the given host.
func (c *Checker) CheckAll(host string) CipherCheckResult {
	result := NewCipherCheckResult(host)

	// Check TLS versions
	result.TLSVersions = c.CheckAllTLSVersions(host)

	// Check TLS 1.2 ciphers
	result.Strong, result.Weak, result.Insecure = c.CheckTLS12Ciphers(host)

	// Check TLS 1.3 ciphers
	result.TLS13 = c.CheckTLS13Ciphers(host)

	// Check post-quantum support
	result.PostQuantum.KeyExchange = c.CheckPostQuantumKeyExchange(host)
	result.PostQuantum.Signatures = c.CheckPostQuantumSignatures(host)

	return result
}

// CheckHost performs a comprehensive TLS check using default configuration.
func CheckHost(host string) CipherCheckResult {
	checker := NewChecker(DefaultCheckerConfig())
	return checker.CheckAll(host)
}
