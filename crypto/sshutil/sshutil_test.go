package sshutil

import (
	"crypto/ed25519"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/ssh"
)

// writeTestPubKey generates an ed25519 keypair, writes its authorized_keys
// line to dir/name, and returns the expected KeyInfo for it.
func writeTestPubKey(t *testing.T, dir, name, comment string) KeyInfo {
	t.Helper()

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("ed25519.GenerateKey() error = %v", err)
	}
	sshPub, err := ssh.NewPublicKey(pub)
	if err != nil {
		t.Fatalf("ssh.NewPublicKey() error = %v", err)
	}

	line := ssh.MarshalAuthorizedKey(sshPub)
	line = append(line[:len(line)-1], []byte(" "+comment+"\n")...) // replace trailing newline with " comment\n"

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, line, 0o600); err != nil {
		t.Fatalf("os.WriteFile(%s) error = %v", path, err)
	}

	return KeyInfo{
		Path:        path,
		Type:        sshPub.Type(),
		Comment:     comment,
		Fingerprint: ssh.FingerprintSHA256(sshPub),
	}
}

func TestNormalizeFingerprint(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"AbCdEfGhIjKlMnOpQrStUvWxYz0123456789ABCDEF", "SHA256:AbCdEfGhIjKlMnOpQrStUvWxYz0123456789ABCDEF"},
		{"SHA256:AbCdEfGhIjKlMnOpQrStUvWxYz0123456789ABCDEF", "SHA256:AbCdEfGhIjKlMnOpQrStUvWxYz0123456789ABCDEF"},
		{"  AbCdEfGhIjKlMnOpQrStUvWxYz0123456789ABCDEF  ", "SHA256:AbCdEfGhIjKlMnOpQrStUvWxYz0123456789ABCDEF"},
	}
	for _, tt := range tests {
		if got := NormalizeFingerprint(tt.input); got != tt.want {
			t.Errorf("NormalizeFingerprint(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseAuthorizedKeyFile(t *testing.T) {
	dir := t.TempDir()
	want := writeTestPubKey(t, dir, "id_ed25519.pub", "user@host")

	got, err := ParseAuthorizedKeyFile(want.Path)
	if err != nil {
		t.Fatalf("ParseAuthorizedKeyFile() error = %v", err)
	}
	if got != want {
		t.Errorf("ParseAuthorizedKeyFile() = %+v, want %+v", got, want)
	}
}

func TestParseAuthorizedKeyFile_Errors(t *testing.T) {
	dir := t.TempDir()

	if _, err := ParseAuthorizedKeyFile(filepath.Join(dir, "missing.pub")); err == nil {
		t.Error("ParseAuthorizedKeyFile() on missing file: want error, got nil")
	}

	badPath := filepath.Join(dir, "bad.pub")
	if err := os.WriteFile(badPath, []byte("not a key\n"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	if _, err := ParseAuthorizedKeyFile(badPath); err == nil {
		t.Error("ParseAuthorizedKeyFile() on malformed file: want error, got nil")
	}
}

func TestScanDir(t *testing.T) {
	dir := t.TempDir()
	k1 := writeTestPubKey(t, dir, "a.pub", "a@host")
	k2 := writeTestPubKey(t, dir, "b.pub", "b@host")
	if err := os.WriteFile(filepath.Join(dir, "bad.pub"), []byte("garbage\n"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "id_ed25519"), []byte("private key material\n"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	keys, errs, err := ScanDir(dir)
	if err != nil {
		t.Fatalf("ScanDir() error = %v", err)
	}
	if len(errs) != 1 {
		t.Errorf("ScanDir() errs = %v, want 1 entry for bad.pub", errs)
	}
	if len(keys) != 2 {
		t.Fatalf("ScanDir() returned %d keys, want 2", len(keys))
	}

	byPath := map[string]KeyInfo{keys[0].Path: keys[0], keys[1].Path: keys[1]}
	if got, ok := byPath[k1.Path]; !ok || got != k1 {
		t.Errorf("ScanDir() missing/mismatched entry for %s: got %+v", k1.Path, got)
	}
	if got, ok := byPath[k2.Path]; !ok || got != k2 {
		t.Errorf("ScanDir() missing/mismatched entry for %s: got %+v", k2.Path, got)
	}
}

func TestFindByFingerprint(t *testing.T) {
	dir := t.TempDir()
	target := writeTestPubKey(t, dir, "target.pub", "target@host")
	writeTestPubKey(t, dir, "other.pub", "other@host")

	t.Run("with prefix", func(t *testing.T) {
		matches, _, err := FindByFingerprint(dir, target.Fingerprint)
		if err != nil {
			t.Fatalf("FindByFingerprint() error = %v", err)
		}
		if len(matches) != 1 || matches[0] != target {
			t.Errorf("FindByFingerprint() = %+v, want [%+v]", matches, target)
		}
	})

	t.Run("without prefix", func(t *testing.T) {
		bare := target.Fingerprint[len("SHA256:"):]
		matches, _, err := FindByFingerprint(dir, bare)
		if err != nil {
			t.Fatalf("FindByFingerprint() error = %v", err)
		}
		if len(matches) != 1 || matches[0] != target {
			t.Errorf("FindByFingerprint() = %+v, want [%+v]", matches, target)
		}
	})

	t.Run("no match", func(t *testing.T) {
		matches, _, err := FindByFingerprint(dir, "SHA256:doesnotexist000000000000000000000000000")
		if err != nil {
			t.Fatalf("FindByFingerprint() error = %v", err)
		}
		if len(matches) != 0 {
			t.Errorf("FindByFingerprint() = %+v, want no matches", matches)
		}
	})
}
