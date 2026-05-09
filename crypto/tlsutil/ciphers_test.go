package tlsutil

import (
	"crypto/tls"
	"testing"
)

func TestTLS12CiphersWeakMap(t *testing.T) {
	weakMap := TLS12CiphersWeakMap()

	if len(weakMap) == 0 {
		t.Error("TLS12CiphersWeakMap should not be empty")
	}

	// Verify known weak ciphers are present
	expectedWeak := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_RC4_128_SHA,
	}

	for _, cipher := range expectedWeak {
		if _, ok := weakMap[cipher]; !ok {
			t.Errorf("expected cipher %d to be in weak map", cipher)
		}
	}

	// Verify strong ciphers are NOT in weak map
	strongCiphers := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	}

	for _, cipher := range strongCiphers {
		if _, ok := weakMap[cipher]; ok {
			t.Errorf("cipher %d should NOT be in weak map (it's strong)", cipher)
		}
	}
}

func TestTLS12CiphersStrongMap(t *testing.T) {
	strongMap := TLS12CiphersStrongMap()

	if len(strongMap) == 0 {
		t.Error("TLS12CiphersStrongMap should not be empty")
	}

	// Verify GCM ciphers are present (they are strong)
	expectedStrong := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	}

	for _, cipher := range expectedStrong {
		if _, ok := strongMap[cipher]; !ok {
			t.Errorf("expected cipher %d to be in strong map", cipher)
		}
	}
}

func TestTLS12CiphersWeak(t *testing.T) {
	weakList := TLS12CiphersWeak()

	if len(weakList) == 0 {
		t.Error("TLS12CiphersWeak should not be empty")
	}

	// Verify it returns strings, not IDs
	for _, name := range weakList {
		if name == "" {
			t.Error("cipher name should not be empty")
		}
	}
}

func TestWeakCipherCategories(t *testing.T) {
	weakMap := TLS12CiphersWeakMap()

	// Categorize weak ciphers by vulnerability type
	var cbcCiphers, rc4Ciphers, des3Ciphers, noFSCiphers int

	for id, name := range weakMap {
		_ = id // use id to silence lint
		switch {
		case contains(name, "RC4"):
			rc4Ciphers++
		case contains(name, "3DES"):
			des3Ciphers++
		case contains(name, "CBC") && contains(name, "RSA_WITH"):
			noFSCiphers++ // RSA key exchange without forward secrecy
		case contains(name, "CBC"):
			cbcCiphers++
		}
	}

	t.Logf("Weak cipher breakdown: CBC=%d, RC4=%d, 3DES=%d, No-FS=%d",
		cbcCiphers, rc4Ciphers, des3Ciphers, noFSCiphers)

	// We should have at least some of each category
	if rc4Ciphers == 0 {
		t.Log("warning: no RC4 ciphers in weak map")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
