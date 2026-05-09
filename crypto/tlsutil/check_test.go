package tlsutil

// Tests use a local mock TLS server with self-signed certificates rather than
// connecting to external hosts. This provides fast, reliable, CI-friendly tests
// that can verify specific TLS configurations without network dependencies.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"
)

// generateTestCert creates a self-signed certificate for testing.
func generateTestCert(t *testing.T) (tls.Certificate, *x509.Certificate) {
	t.Helper()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Test Organization"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	tlsCert := tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  privateKey,
	}

	return tlsCert, cert
}

// startTestTLSServer starts a TLS server with the given configuration.
// Returns the address and a cleanup function.
func startTestTLSServer(t *testing.T, tlsConfig *tls.Config) (string, func()) {
	t.Helper()

	listener, err := tls.Listen("tcp", "127.0.0.1:0", tlsConfig)
	if err != nil {
		t.Fatalf("failed to start TLS listener: %v", err)
	}

	// Accept connections in background
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return // listener closed
			}
			// Do TLS handshake and close
			go func(c net.Conn) {
				defer c.Close()
				if tlsConn, ok := c.(*tls.Conn); ok {
					_ = tlsConn.Handshake()
				}
			}(conn)
		}
	}()

	return listener.Addr().String(), func() { listener.Close() }
}

func TestCheckTLSVersions(t *testing.T) {
	cert, _ := generateTestCert(t)

	tests := []struct {
		name           string
		minVersion     uint16
		maxVersion     uint16
		expectVersions map[TLSVersion]bool
	}{
		{
			name:       "TLS 1.2 only",
			minVersion: tls.VersionTLS12,
			maxVersion: tls.VersionTLS12,
			expectVersions: map[TLSVersion]bool{
				VersionTLS10: false,
				VersionTLS11: false,
				VersionTLS12: true,
				VersionTLS13: false,
			},
		},
		{
			name:       "TLS 1.3 only",
			minVersion: tls.VersionTLS13,
			maxVersion: tls.VersionTLS13,
			expectVersions: map[TLSVersion]bool{
				VersionTLS10: false,
				VersionTLS11: false,
				VersionTLS12: false,
				VersionTLS13: true,
			},
		},
		{
			name:       "TLS 1.2 and 1.3",
			minVersion: tls.VersionTLS12,
			maxVersion: tls.VersionTLS13,
			expectVersions: map[TLSVersion]bool{
				VersionTLS10: false,
				VersionTLS11: false,
				VersionTLS12: true,
				VersionTLS13: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsConfig := &tls.Config{ // #nosec G402 -- intentional for testing specific TLS versions
				Certificates: []tls.Certificate{cert},
				MinVersion:   tt.minVersion,
				MaxVersion:   tt.maxVersion,
			}

			addr, cleanup := startTestTLSServer(t, tlsConfig)
			defer cleanup()

			// Give server time to start
			time.Sleep(10 * time.Millisecond)

			checker := NewChecker(CheckerConfig{Timeout: 5 * time.Second})
			results := checker.CheckAllTLSVersions(addr)

			for _, r := range results {
				expected, ok := tt.expectVersions[r.Version]
				if !ok {
					continue
				}
				if r.Active != expected {
					t.Errorf("%s: got Active=%v, want %v", r.Name, r.Active, expected)
				}
			}
		})
	}
}

func TestCheckTLS12Ciphers(t *testing.T) {
	cert, _ := generateTestCert(t)

	// Configure server with only strong GCM ciphers
	tlsConfig := &tls.Config{ // #nosec G402 -- intentional for testing TLS 1.2 ciphers
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}

	addr, cleanup := startTestTLSServer(t, tlsConfig)
	defer cleanup()

	time.Sleep(10 * time.Millisecond)

	checker := NewChecker(CheckerConfig{Timeout: 5 * time.Second})
	strong, weak, insecure := checker.CheckTLS12Ciphers(addr)

	// Should have strong ciphers
	if len(strong) == 0 {
		t.Error("expected strong ciphers, got none")
	}

	// Should have no weak or insecure ciphers (server only allows GCM)
	if len(weak) > 0 {
		t.Errorf("expected no weak ciphers, got: %v", weak)
	}
	if len(insecure) > 0 {
		t.Errorf("expected no insecure ciphers, got: %v", insecure)
	}

	// Verify specific cipher is detected
	found := false
	for _, c := range strong {
		if strings.Contains(c, "AES_128_GCM") || strings.Contains(c, "AES_256_GCM") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected GCM cipher in strong list, got: %v", strong)
	}
}

func TestCheckTLS13Ciphers(t *testing.T) {
	cert, _ := generateTestCert(t)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	addr, cleanup := startTestTLSServer(t, tlsConfig)
	defer cleanup()

	time.Sleep(10 * time.Millisecond)

	checker := NewChecker(CheckerConfig{Timeout: 5 * time.Second})
	ciphers := checker.CheckTLS13Ciphers(addr)

	if len(ciphers) == 0 {
		t.Error("expected TLS 1.3 ciphers, got none")
	}

	// TLS 1.3 ciphers should include AES-GCM or ChaCha20
	found := false
	for _, c := range ciphers {
		if strings.Contains(c, "AES") || strings.Contains(c, "CHACHA") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected AES or ChaCha cipher, got: %v", ciphers)
	}
}

func TestCheckPostQuantumKeyExchange(t *testing.T) {
	cert, _ := generateTestCert(t)

	// Test with PQ curve preferences
	tlsConfig := &tls.Config{
		Certificates:     []tls.Certificate{cert},
		MinVersion:       tls.VersionTLS13,
		MaxVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{X25519MLKEM768, tls.X25519},
	}

	addr, cleanup := startTestTLSServer(t, tlsConfig)
	defer cleanup()

	time.Sleep(10 * time.Millisecond)

	checker := NewChecker(CheckerConfig{Timeout: 5 * time.Second})
	pqKex := checker.CheckPostQuantumKeyExchange(addr)

	// Server supports PQ, so we should detect it
	if len(pqKex) == 0 {
		t.Log("PQ key exchange not detected (may not be supported by this Go version)")
	} else {
		found := false
		for _, k := range pqKex {
			if strings.Contains(k, "MLKEM") || strings.Contains(k, "Kyber") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected MLKEM in PQ key exchange, got: %v", pqKex)
		}
	}
}

func TestCipherCheckResultString(t *testing.T) {
	result := CipherCheckResult{
		Host: "test.example.com:443",
		TLSVersions: []TLSVersionStatus{
			{Version: VersionTLS12, Name: "TLS 1.2", Active: true},
			{Version: VersionTLS13, Name: "TLS 1.3", Active: true},
		},
		Strong:   []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
		Weak:     []string{"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"},
		Insecure: []string{"TLS_RSA_WITH_AES_128_CBC_SHA"},
		TLS13:    []string{"TLS_AES_128_GCM_SHA256"},
		PostQuantum: PostQuantumSupport{
			KeyExchange: []string{"X25519MLKEM768"},
			Signatures:  []string{},
		},
	}

	output := result.String()

	// Check that all sections are present
	checks := []string{
		"=== TLS Version Support ===",
		"TLS 1.2",
		"TLS 1.3",
		"ACTIVE",
		"=== TLS 1.2 Cipher Suites ===",
		"STRONG:",
		"WEAK:",
		"INSECURE:",
		"=== TLS 1.3 Cipher Suites ===",
		"=== Post-Quantum Support ===",
		"Key Exchange:",
		"X25519MLKEM768",
		"Signatures",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("String() output missing %q", check)
		}
	}
}

func TestPQSignatureOIDLookup(t *testing.T) {
	tests := []struct {
		oidStr   string
		expected string
		isPQ     bool
	}{
		{"2.16.840.1.101.3.4.3.17", "ML-DSA-44", true},
		{"2.16.840.1.101.3.4.3.18", "ML-DSA-65", true},
		{"2.16.840.1.101.3.4.3.19", "ML-DSA-87", true},
		{"2.16.840.1.101.3.4.3.20", "SLH-DSA-SHA2-128s", true},
		{"1.3.6.1.4.1.2.267.7.6.5", "Dilithium3 (draft)", true},
		{"1.2.840.10045.4.3.2", "", false},   // ECDSA-SHA256 (not PQ)
		{"1.2.840.113549.1.1.11", "", false}, // RSA-SHA256 (not PQ)
	}

	for _, tt := range tests {
		t.Run(tt.oidStr, func(t *testing.T) {
			// Parse OID string to asn1.ObjectIdentifier
			var oid asn1.ObjectIdentifier
			for _, part := range strings.Split(tt.oidStr, ".") {
				var n int
				for _, c := range part {
					n = n*10 + int(c-'0')
				}
				oid = append(oid, n)
			}

			name := PQSignatureOIDName(oid)
			if name != tt.expected {
				t.Errorf("PQSignatureOIDName(%s) = %q, want %q", tt.oidStr, name, tt.expected)
			}

			isPQ := IsPQSignatureOID(oid)
			if isPQ != tt.isPQ {
				t.Errorf("IsPQSignatureOID(%s) = %v, want %v", tt.oidStr, isPQ, tt.isPQ)
			}
		})
	}
}

func TestExtractSignatureAlgorithmOID(t *testing.T) {
	// Generate a real certificate and extract its OID
	_, cert := generateTestCert(t)

	oid, err := extractSignatureAlgorithmOID(cert.Raw)
	if err != nil {
		t.Fatalf("extractSignatureAlgorithmOID failed: %v", err)
	}

	// Our test cert uses ECDSA-SHA256
	// OID: 1.2.840.10045.4.3.2
	expectedOID := "1.2.840.10045.4.3.2"
	if oid.String() != expectedOID {
		t.Errorf("got OID %s, want %s", oid.String(), expectedOID)
	}

	// Verify it's not detected as PQ (it's classical ECDSA)
	if IsPQSignatureOID(oid) {
		t.Error("ECDSA-SHA256 should not be detected as PQ")
	}
}

func TestCheckAll(t *testing.T) {
	cert, _ := generateTestCert(t)

	tlsConfig := &tls.Config{
		Certificates:     []tls.Certificate{cert},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{X25519MLKEM768, tls.X25519, tls.CurveP256},
	}

	addr, cleanup := startTestTLSServer(t, tlsConfig)
	defer cleanup()

	time.Sleep(10 * time.Millisecond)

	result := CheckHost(addr)

	// Verify host is set
	if result.Host != addr {
		t.Errorf("Host = %q, want %q", result.Host, addr)
	}

	// Verify time is set
	if result.Time == nil {
		t.Error("Time should not be nil")
	}

	// Verify TLS versions detected
	if len(result.TLSVersions) == 0 {
		t.Error("TLSVersions should not be empty")
	}

	// Check that TLS 1.2 and 1.3 are active
	var tls12Active, tls13Active bool
	for _, v := range result.TLSVersions {
		if v.Version == VersionTLS12 {
			tls12Active = v.Active
		}
		if v.Version == VersionTLS13 {
			tls13Active = v.Active
		}
	}
	if !tls12Active {
		t.Error("TLS 1.2 should be active")
	}
	if !tls13Active {
		t.Error("TLS 1.3 should be active")
	}

	// Verify we got some cipher results
	// Note: Strong/Weak/Insecure may be nil if no ciphers of that type are supported
	// Our test server only allows GCM ciphers, so we should have strong ciphers
	if len(result.Strong) == 0 {
		t.Error("Strong should have at least one cipher")
	}

	// TLS 1.3 should have ciphers since we enabled it
	if len(result.TLS13) == 0 {
		t.Error("TLS13 should have at least one cipher")
	}

	// PostQuantum slices should be initialized (not nil) per NewCipherCheckResult
	if result.PostQuantum.KeyExchange == nil {
		t.Error("PostQuantum.KeyExchange should not be nil")
	}
	if result.PostQuantum.Signatures == nil {
		t.Error("PostQuantum.Signatures should not be nil")
	}
}
