// Package sshutil provides helpers for working with SSH public keys, such as
// parsing `authorized_keys`-format `.pub` files and matching them against a
// SHA256 fingerprint (e.g. `SHA256:AbCdEfGhIjKlMnOpQrStUvWxYz0123456789ABCDEF`).
package sshutil

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/crypto/ssh"
)

// KeyInfo describes a single parsed SSH public key file.
type KeyInfo struct {
	Path        string // filesystem path to the .pub file
	Type        string // e.g. "ssh-ed25519", "ecdsa-sha2-nistp521"
	Comment     string // trailing comment, typically "user@host"
	Fingerprint string // SHA256 fingerprint, e.g. "SHA256:eLpJ..."
}

// NormalizeFingerprint adds the "SHA256:" prefix if the caller omitted it.
func NormalizeFingerprint(fingerprint string) string {
	fingerprint = strings.TrimSpace(fingerprint)
	if !strings.HasPrefix(fingerprint, "SHA256:") {
		fingerprint = "SHA256:" + fingerprint
	}
	return fingerprint
}

// ParseAuthorizedKeyFile parses a single `authorized_keys`-format file, such
// as a `~/.ssh/id_ed25519.pub` file, and returns its KeyInfo.
func ParseAuthorizedKeyFile(path string) (KeyInfo, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return KeyInfo{}, fmt.Errorf("sshutil: read %s: %w", path, err)
	}
	pub, comment, _, _, err := ssh.ParseAuthorizedKey(b)
	if err != nil {
		return KeyInfo{}, fmt.Errorf("sshutil: parse %s: %w", path, err)
	}
	return KeyInfo{
		Path:        path,
		Type:        pub.Type(),
		Comment:     comment,
		Fingerprint: ssh.FingerprintSHA256(pub),
	}, nil
}

// ScanDir parses every `*.pub` file directly inside dir. It does not
// recurse into subdirectories. Files that fail to parse are skipped for the
// keys result but reported in errs so callers can surface them rather than
// having them silently dropped.
func ScanDir(dir string) (keys []KeyInfo, errs []error, err error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.pub"))
	if err != nil {
		return nil, nil, fmt.Errorf("sshutil: glob %s: %w", dir, err)
	}
	sort.Strings(paths)

	for _, path := range paths {
		key, parseErr := ParseAuthorizedKeyFile(path)
		if parseErr != nil {
			errs = append(errs, parseErr)
			continue
		}
		keys = append(keys, key)
	}
	return keys, errs, nil
}

// FindByFingerprint scans dir for `*.pub` files whose SHA256 fingerprint
// matches target. The "SHA256:" prefix on target is optional. Files that
// fail to parse do not abort the search; they are returned via errs.
func FindByFingerprint(dir, target string) (matches []KeyInfo, errs []error, err error) {
	keys, errs, err := ScanDir(dir)
	if err != nil {
		return nil, errs, err
	}

	target = NormalizeFingerprint(target)
	for _, key := range keys {
		if key.Fingerprint == target {
			matches = append(matches, key)
		}
	}
	return matches, errs, nil
}
