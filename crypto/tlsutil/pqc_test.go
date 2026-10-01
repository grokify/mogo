package tlsutil

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCurveIDName(t *testing.T) {
	tests := []struct {
		name     string
		curveID  tls.CurveID
		expected string
	}{
		{"P-256", tls.CurveP256, "P-256"},
		{"P-384", tls.CurveP384, "P-384"},
		{"P-521", tls.CurveP521, "P-521"},
		{"X25519", tls.X25519, "X25519"},
		{"X25519MLKEM768", tls.X25519MLKEM768, "X25519MLKEM768"},
		{"SecP256r1MLKEM768", tls.SecP256r1MLKEM768, "SecP256r1MLKEM768"},
		{"SecP384r1MLKEM1024", tls.SecP384r1MLKEM1024, "SecP384r1MLKEM1024"},
		{"MLKEM1024", tls.CurveID(514), "MLKEM1024"},
		{"Unknown", tls.CurveID(9999), "CurveID(9999)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CurveIDName(tt.curveID)
			if got != tt.expected {
				t.Errorf("CurveIDName(%d) = %q, want %q", tt.curveID, got, tt.expected)
			}
		})
	}
}

func TestIsPQCCurve(t *testing.T) {
	tests := []struct {
		name     string
		curveID  tls.CurveID
		expected bool
	}{
		{"X25519MLKEM768 is PQC", tls.X25519MLKEM768, true},
		{"SecP256r1MLKEM768 is PQC", tls.SecP256r1MLKEM768, true},
		{"SecP384r1MLKEM1024 is PQC", tls.SecP384r1MLKEM1024, true},
		{"MLKEM1024 is PQC", tls.CurveID(514), true},
		{"X25519 is not PQC", tls.X25519, false},
		{"P-256 is not PQC", tls.CurveP256, false},
		{"P-384 is not PQC", tls.CurveP384, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPQCCurve(tt.curveID)
			if got != tt.expected {
				t.Errorf("IsPQCCurve(%d) = %v, want %v", tt.curveID, got, tt.expected)
			}
		})
	}
}

func TestCurveIDToPQCAlgorithm(t *testing.T) {
	tests := []struct {
		name          string
		curveID       tls.CurveID
		expectedAlgo  PQCAlgorithm
		expectedFound bool
	}{
		{"X25519MLKEM768", tls.X25519MLKEM768, PQCAlgorithmMLKEM768, true},
		{"SecP256r1MLKEM768", tls.SecP256r1MLKEM768, PQCAlgorithmMLKEM768, true},
		{"SecP384r1MLKEM1024", tls.SecP384r1MLKEM1024, PQCAlgorithmMLKEM1024, true},
		{"MLKEM1024", tls.CurveID(514), PQCAlgorithmMLKEM1024, true},
		{"X25519", tls.X25519, "", false},
		{"P-256", tls.CurveP256, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			algo, found := CurveIDToPQCAlgorithm(tt.curveID)
			if algo != tt.expectedAlgo || found != tt.expectedFound {
				t.Errorf("CurveIDToPQCAlgorithm(%d) = (%q, %v), want (%q, %v)",
					tt.curveID, algo, found, tt.expectedAlgo, tt.expectedFound)
			}
		})
	}
}

func TestPQCCurvePreferences(t *testing.T) {
	prefs := PQCCurvePreferences()
	if len(prefs) == 0 {
		t.Error("PQCCurvePreferences() returned empty slice")
	}
	if prefs[0] != tls.X25519MLKEM768 {
		t.Errorf("PQCCurvePreferences()[0] = %d, want %d (X25519MLKEM768)", prefs[0], tls.X25519MLKEM768)
	}
}

func TestPQCAlgorithms(t *testing.T) {
	algos := PQCAlgorithms()
	if len(algos) != 5 {
		t.Errorf("PQCAlgorithms() returned %d algorithms, want 5", len(algos))
	}

	// Check that ML-KEM-768 is first and has correct properties
	if algos[0].Algorithm != PQCAlgorithmMLKEM768 {
		t.Errorf("First algorithm = %q, want %q", algos[0].Algorithm, PQCAlgorithmMLKEM768)
	}
	if algos[0].Type != PQCAlgorithmTypeKEM {
		t.Errorf("ML-KEM-768 type = %q, want %q", algos[0].Type, PQCAlgorithmTypeKEM)
	}
	if !algos[0].StdlibCheck {
		t.Error("ML-KEM-768 should have StdlibCheck = true")
	}
}

// TestCheckPQCSupportLocal negotiates each curve against a local TLS server
// restricted to that curve, so detection is verified without network access.
func TestCheckPQCSupportLocal(t *testing.T) {
	tests := []struct {
		curve    tls.CurveID
		wantPQC  bool
		wantAlgo PQCAlgorithm
	}{
		{tls.X25519MLKEM768, true, PQCAlgorithmMLKEM768},
		{tls.SecP256r1MLKEM768, true, PQCAlgorithmMLKEM768},
		{tls.SecP384r1MLKEM1024, true, PQCAlgorithmMLKEM1024},
		{tls.X25519, false, ""},
	}
	for _, tt := range tests {
		t.Run(CurveIDName(tt.curve), func(t *testing.T) {
			srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			srv.TLS = &tls.Config{MinVersion: tls.VersionTLS13, CurvePreferences: []tls.CurveID{tt.curve}}
			srv.StartTLS()
			defer srv.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			result := checkPQCSupport(ctx, srv.URL, srv.Client().Transport.(*http.Transport).TLSClientConfig.RootCAs)

			if result.Error != "" {
				t.Fatalf("checkPQCSupport() error: %s", result.Error)
			}
			if !result.Supported || result.TLSVersion != "TLS 1.3" {
				t.Errorf("Supported = %v, TLSVersion = %q, want true, TLS 1.3", result.Supported, result.TLSVersion)
			}
			if result.CurveID != tt.curve {
				t.Errorf("CurveID = %s, want %s", CurveIDName(result.CurveID), CurveIDName(tt.curve))
			}
			if result.PQCKeyExchange != tt.wantPQC || result.PQCAlgorithm != tt.wantAlgo {
				t.Errorf("PQC = (%v, %q), want (%v, %q)", result.PQCKeyExchange, result.PQCAlgorithm, tt.wantPQC, tt.wantAlgo)
			}
		})
	}
}

// TestCheckPQCSupportProviders tests PQC support against known PQC-enabled providers.
// Cloudflare and Google have enabled PQC (X25519+ML-KEM-768) for their services.
// This is an integration test that requires network access.
func TestCheckPQCSupportProviders(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tests := []struct {
		name string
		url  string
	}{
		{"Cloudflare", "https://cloudflare.com"},
		{"Google", "https://www.google.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result := CheckPQCSupport(ctx, tt.url)

			if result.Error != "" {
				t.Fatalf("CheckPQCSupport(%s) failed: %s", tt.url, result.Error)
			}

			if !result.Supported {
				t.Errorf("Expected %s to support TLS connection", tt.name)
			}

			if result.TLSVersion != "TLS 1.3" {
				t.Errorf("Expected TLS 1.3, got %s", result.TLSVersion)
			}

			t.Logf("%s PQC check result:", tt.name)
			t.Logf("  TLS Version: %s", result.TLSVersion)
			t.Logf("  Curve ID: %d", result.CurveID)
			t.Logf("  Curve Name: %s", result.CurveName)
			t.Logf("  PQC Key Exchange: %v", result.PQCKeyExchange)
			t.Logf("  PQC Algorithm: %s", result.PQCAlgorithm)

			if result.PQCKeyExchange {
				t.Logf("%s supports PQC key exchange (X25519+ML-KEM-768)", tt.name)
				if result.PQCAlgorithm != PQCAlgorithmMLKEM768 {
					t.Errorf("Expected PQC algorithm %q, got %q", PQCAlgorithmMLKEM768, result.PQCAlgorithm)
				}
			} else {
				t.Logf("%s did not negotiate PQC key exchange (server may have disabled it or client Go version doesn't support it)", tt.name)
			}
		})
	}
}

// TestCheckPQCURLsWithSummary tests the summary function.
func TestCheckPQCURLsWithSummary(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	urls := []string{
		"https://cloudflare.com",
		"https://www.google.com",
	}

	summary := CheckPQCURLsWithSummary(ctx, urls)

	if summary.TotalChecked != 2 {
		t.Errorf("TotalChecked = %d, want 2", summary.TotalChecked)
	}

	t.Logf("PQC Support Summary:")
	t.Logf("  Total Checked: %d", summary.TotalChecked)
	t.Logf("  PQC Supported: %d", summary.PQCSupported)
	t.Logf("  TLS 1.3 Only: %d", summary.TLS13Only)
	t.Logf("  Failed: %d", summary.Failed)

	for _, r := range summary.Results {
		t.Logf("  %s: PQC=%v, Curve=%s", r.URL, r.PQCKeyExchange, r.CurveName)
	}
}

// TestCheckPQCSupportInvalidURL tests error handling for invalid URLs.
func TestCheckPQCSupportInvalidURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := CheckPQCSupport(ctx, "https://invalid.invalid.invalid")

	if result.Error == "" {
		t.Error("Expected error for invalid URL, got none")
	}

	if result.Supported {
		t.Error("Expected Supported=false for invalid URL")
	}
}
