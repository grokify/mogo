package aesutil

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func testKey(t *testing.T, size int) []byte {
	t.Helper()
	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	return key
}

// overhead is the 12-byte nonce plus the 16-byte GCM tag.
const overhead = 28

func TestEncryptDecryptAES(t *testing.T) {
	for _, size := range []int{16, 24, 32} {
		key := testKey(t, size)
		for _, plaintext := range [][]byte{nil, []byte("x"), []byte("hello, world"), bytes.Repeat([]byte{0xff}, 10000)} {
			ct, err := EncryptAES(plaintext, key)
			if err != nil {
				t.Fatalf("EncryptAES(%d-byte key, %d bytes) error: %v", size, len(plaintext), err)
			}
			if len(ct) != len(plaintext)+overhead {
				t.Errorf("ciphertext length = %d, want %d", len(ct), len(plaintext)+overhead)
			}
			pt, err := DecryptAES(ct, key)
			if err != nil {
				t.Fatalf("DecryptAES(%d-byte key, %d bytes) error: %v", size, len(plaintext), err)
			}
			if !bytes.Equal(pt, plaintext) {
				t.Errorf("round trip with %d-byte key = %q, want %q", size, pt, plaintext)
			}
		}
	}
}

func TestEncryptAESFreshNonce(t *testing.T) {
	key := testKey(t, 32)
	a, err := EncryptAES([]byte("same"), key)
	if err != nil {
		t.Fatal(err)
	}
	b, err := EncryptAES([]byte("same"), key)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(a, b) {
		t.Error("encrypting the same plaintext twice produced identical ciphertext")
	}
}

func TestDecryptAESRejects(t *testing.T) {
	key := testKey(t, 32)
	ct, err := EncryptAES([]byte("attack at dawn"), key)
	if err != nil {
		t.Fatal(err)
	}
	for i := range ct { // flipping any bit of nonce, ciphertext, or tag must fail
		tampered := bytes.Clone(ct)
		tampered[i] ^= 0x01
		if _, err := DecryptAES(tampered, key); err == nil {
			t.Fatalf("DecryptAES accepted ciphertext with byte %d modified", i)
		}
	}
	if _, err := DecryptAES(ct, testKey(t, 32)); err == nil {
		t.Error("DecryptAES accepted the wrong key")
	}
	if _, err := DecryptAES(ct[:overhead-1], key); err == nil {
		t.Error("DecryptAES accepted truncated ciphertext")
	}
	if _, err := DecryptAES(nil, key); err == nil {
		t.Error("DecryptAES accepted empty ciphertext")
	}
}

// TestDecryptAESRejectsLegacyCFB checks that ciphertext from the pre-v0.75.0
// AES-CFB format fails with an error rather than decrypting to garbage. The
// fixture is "legacy secret" encrypted by the v0.74.10 EncryptAES.
func TestDecryptAESRejectsLegacyCFB(t *testing.T) {
	legacy, err := hex.DecodeString("421019c03f84ece33094a8d7b0c51d7c878b9e1c6daeb554b9bbe807cfc79087debd3db6")
	if err != nil {
		t.Fatal(err)
	}
	if pt, err := DecryptAES(legacy, []byte("0123456789abcdef0123456789abcdef")); err == nil {
		t.Errorf("DecryptAES(legacy CFB ciphertext) = %q, want error", pt)
	}
}

func TestDecryptAESDoesNotModifyInput(t *testing.T) {
	key := testKey(t, 16)
	ct, err := EncryptAES([]byte("keep me"), key)
	if err != nil {
		t.Fatal(err)
	}
	orig := bytes.Clone(ct)
	if _, err := DecryptAES(ct, key); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ct, orig) {
		t.Error("DecryptAES modified its ciphertext argument")
	}
}

func TestInvalidKeySize(t *testing.T) {
	key := make([]byte, 15)
	if _, err := EncryptAES([]byte("x"), key); err == nil {
		t.Error("EncryptAES accepted a 15-byte key")
	}
	if _, err := DecryptAES(make([]byte, 64), key); err == nil {
		t.Error("DecryptAES accepted a 15-byte key")
	}
}

func TestBase58AndJSON(t *testing.T) {
	key := testKey(t, 32)
	ct, err := EncryptAESBase58([]byte("base58 text"), key)
	if err != nil {
		t.Fatal(err)
	}
	pt, err := DecryptAESBase58(ct, key)
	if err != nil || string(pt) != "base58 text" {
		t.Errorf("Base58 round trip = (%q, %v), want %q", pt, err, "base58 text")
	}
	if _, err := DecryptAESBase58([]byte("not-base58-0OIl"), key); err == nil {
		t.Error("DecryptAESBase58 accepted invalid input")
	}

	type item struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	in := item{Name: "widget", Count: 3}
	ctJSON, err := EncryptAESBase58JSON(in, key)
	if err != nil {
		t.Fatal(err)
	}
	var out item
	if err := DecryptAESBase58JSON(ctJSON, key, &out); err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Errorf("JSON round trip = %+v, want %+v", out, in)
	}
}

func TestDecryptAESBase64String(t *testing.T) {
	key := testKey(t, 32)
	ct, err := EncryptAES([]byte("b64"), key)
	if err != nil {
		t.Fatal(err)
	}
	pt, err := DecryptAESBase64String(base64.StdEncoding.EncodeToString(ct), key)
	if err != nil || string(pt) != "b64" {
		t.Errorf("DecryptAESBase64String = (%q, %v), want %q", pt, err, "b64")
	}
	if _, err := DecryptAESBase64String("!!!", key); err == nil {
		t.Error("DecryptAESBase64String accepted invalid base64")
	}
}

func TestFileHelpers(t *testing.T) {
	key := testKey(t, 32)
	dir := t.TempDir()

	enc := filepath.Join(dir, "secret.enc")
	if err := WriteFileAES(enc, []byte("file data"), 0o600, key); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(enc)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte("file data")) {
		t.Error("WriteFileAES wrote plaintext to disk")
	}
	pt, err := ReadFileAES(enc, key)
	if err != nil || string(pt) != "file data" {
		t.Errorf("ReadFileAES = (%q, %v), want %q", pt, err, "file data")
	}

	plain := filepath.Join(dir, "plain.txt")
	if err := os.WriteFile(plain, []byte("in place"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := EncryptFileAES(plain, "", 0o600, key); err != nil { // empty target encrypts in place
		t.Fatal(err)
	}
	if pt, err := ReadFileAES(plain, key); err != nil || string(pt) != "in place" {
		t.Errorf("EncryptFileAES in place: ReadFileAES = (%q, %v), want %q", pt, err, "in place")
	}
}

func TestEncryptDirectoryFilesAES(t *testing.T) {
	key := testKey(t, 32)
	src, dst := t.TempDir(), t.TempDir()
	files := map[string]string{"a.txt": "alpha", "b.json": `{"b":1}`}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(src, name), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(src, "subdir"), 0o700); err != nil { // skipped
		t.Fatal(err)
	}
	if err := EncryptDirectoryFilesAES(src, dst, 0o600, key); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != len(files) {
		t.Errorf("wrote %d files, want %d", len(entries), len(files))
	}
	for name, body := range files {
		pt, err := ReadFileAES(filepath.Join(dst, name+".enc"), key)
		if err != nil || string(pt) != body {
			t.Errorf("%s.enc decrypts to (%q, %v), want %q", name, pt, err, body)
		}
	}
}
