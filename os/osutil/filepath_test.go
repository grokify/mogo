package osutil

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var existTests = []struct {
	filename string
	exists   bool
}{
	{"exist.txt", true},
	{"doesnotexist.txt", false},
}

func TestExist(t *testing.T) {
	for _, tt := range existTests {
		exists, err := Exists(filepath.Join("filepath_testdata", tt.filename))
		if err != nil {
			t.Errorf("osutil.Exists(\"%s\") Error [%s]", tt.filename, err.Error())
		}
		if exists != tt.exists {
			t.Errorf("osutil.Exists(\"%s\") Want [%v] Got [%v]", tt.filename, tt.exists, exists)
		}
	}
}

func TestSanitizePath(t *testing.T) {
	// Get current working directory for test paths
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	testFile := filepath.Join(cwd, "filepath_testdata", "exist.txt")
	testDir := filepath.Join(cwd, "filepath_testdata")

	tests := []struct {
		name    string
		path    string
		opts    *SanitizeOpts
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty path",
			path:    "",
			opts:    nil,
			wantErr: true,
			errMsg:  "path is empty",
		},
		{
			name:    "whitespace only path",
			path:    "   ",
			opts:    nil,
			wantErr: true,
			errMsg:  "path is empty",
		},
		{
			name:    "relative path without opts",
			path:    "filepath_testdata/exist.txt",
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "existing file with MustExist",
			path:    "filepath_testdata/exist.txt",
			opts:    &SanitizeOpts{MustExist: true},
			wantErr: false,
		},
		{
			name:    "non-existing file with MustExist",
			path:    "filepath_testdata/doesnotexist.txt",
			opts:    &SanitizeOpts{MustExist: true},
			wantErr: true,
			errMsg:  "path does not exist",
		},
		{
			name:    "file with MustBeFile",
			path:    "filepath_testdata/exist.txt",
			opts:    &SanitizeOpts{MustExist: true, MustBeFile: true},
			wantErr: false,
		},
		{
			name:    "directory with MustBeFile",
			path:    "filepath_testdata",
			opts:    &SanitizeOpts{MustExist: true, MustBeFile: true},
			wantErr: true,
			errMsg:  "path is not a file",
		},
		{
			name:    "directory with MustBeDir",
			path:    "filepath_testdata",
			opts:    &SanitizeOpts{MustExist: true, MustBeDir: true},
			wantErr: false,
		},
		{
			name:    "file with MustBeDir",
			path:    "filepath_testdata/exist.txt",
			opts:    &SanitizeOpts{MustExist: true, MustBeDir: true},
			wantErr: true,
			errMsg:  "path is not a directory",
		},
		{
			name:    "valid extension",
			path:    "filepath_testdata/exist.txt",
			opts:    &SanitizeOpts{AllowedExts: []string{".txt", ".md"}},
			wantErr: false,
		},
		{
			name:    "invalid extension",
			path:    "filepath_testdata/exist.txt",
			opts:    &SanitizeOpts{AllowedExts: []string{".mp4", ".srt"}},
			wantErr: true,
			errMsg:  "invalid file extension",
		},
		{
			name:    "case insensitive extension",
			path:    "filepath_testdata/exist.txt",
			opts:    &SanitizeOpts{AllowedExts: []string{".TXT"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizePath(tt.path, tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SanitizePath() expected error containing %q, got nil", tt.errMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("SanitizePath() error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("SanitizePath() unexpected error: %v", err)
				return
			}

			// Verify the result is an absolute path
			if !filepath.IsAbs(result) {
				t.Errorf("SanitizePath() result = %q, want absolute path", result)
			}

			// Verify path is cleaned (no . or ..)
			if strings.Contains(result, "..") || strings.HasSuffix(result, "/.") {
				t.Errorf("SanitizePath() result = %q, expected cleaned path", result)
			}
		})
	}

	// Test that result points to the correct absolute path
	t.Run("correct absolute path", func(t *testing.T) {
		result, err := SanitizePath("filepath_testdata/exist.txt", nil)
		if err != nil {
			t.Fatalf("SanitizePath() unexpected error: %v", err)
		}
		if result != testFile {
			t.Errorf("SanitizePath() = %q, want %q", result, testFile)
		}
	})

	t.Run("correct absolute path for directory", func(t *testing.T) {
		result, err := SanitizePath("filepath_testdata", nil)
		if err != nil {
			t.Fatalf("SanitizePath() unexpected error: %v", err)
		}
		if result != testDir {
			t.Errorf("SanitizePath() = %q, want %q", result, testDir)
		}
	})
}

func TestValidatePathComponent(t *testing.T) {
	tests := []struct {
		name      string
		component string
		wantErr   bool
	}{
		{"valid alphanumeric", "abc123", false},
		{"valid with hyphen", "INIT-VISIONSTUDIO-001", false},
		{"valid with underscore", "period_2026_W30", false},
		{"empty", "", true},
		{"dot", ".", true},
		{"dot-dot", "..", true},
		{"path traversal", "../../etc/passwd", true},
		{"forward slash", "foo/bar", true},
		{"backslash", "foo\\bar", true},
		{"leading traversal no separator", "..secret", true},
		{"space", "foo bar", true},
		{"null byte", "foo\x00bar", true},
		{"unicode lookalike slash", "foo∕bar", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePathComponent(tt.component)
			if tt.wantErr && err == nil {
				t.Errorf("ValidatePathComponent(%q) = nil, want error", tt.component)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ValidatePathComponent(%q) = %v, want nil", tt.component, err)
			}
			if tt.wantErr && err != nil && !errors.Is(err, ErrInvalidPathComponent) {
				t.Errorf("ValidatePathComponent(%q) error = %v, want wrapped ErrInvalidPathComponent", tt.component, err)
			}
		})
	}
}

func TestJoinSecure(t *testing.T) {
	root := filepath.FromSlash("/a/b")
	tests := []struct {
		name    string
		root    string
		elem    []string
		want    string
		wantErr bool
	}{
		{"direct child", root, []string{"c"}, filepath.FromSlash("/a/b/c"), false},
		{"nested child", root, []string{"c", "d.md"}, filepath.FromSlash("/a/b/c/d.md"), false},
		{"id plus suffix", root, []string{"myid.json"}, filepath.FromSlash("/a/b/myid.json"), false},
		{"escapes via ..", root, []string{"../../etc/passwd"}, "", true},
		{"escapes via embedded ..", root, []string{"c/../../etc/passwd"}, "", true},
		{"sibling that shares a prefix", root, []string{"../bc"}, "", true},
		// filepath.Join treats every element as a segment to clean, not an
		// override — an absolute-looking elem still lands under root.
		{"absolute-looking element still joins under root", root, []string{"/etc/passwd"}, filepath.FromSlash("/a/b/etc/passwd"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JoinSecure(tt.root, tt.elem...)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("JoinSecure(%q, %v) = %q, nil; want error", tt.root, tt.elem, got)
				}
				if !errors.Is(err, ErrPathEscapesRoot) {
					t.Errorf("JoinSecure(%q, %v) error = %v, want wrapped ErrPathEscapesRoot", tt.root, tt.elem, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("JoinSecure(%q, %v) unexpected error: %v", tt.root, tt.elem, err)
			}
			if got != tt.want {
				t.Errorf("JoinSecure(%q, %v) = %q, want %q", tt.root, tt.elem, got, tt.want)
			}
		})
	}

	t.Run("exact root with no elem", func(t *testing.T) {
		dir := t.TempDir()
		got, err := JoinSecure(dir)
		if err != nil {
			t.Fatalf("JoinSecure(root) unexpected error: %v", err)
		}
		if filepath.Clean(got) != filepath.Clean(dir) {
			t.Errorf("JoinSecure(root) = %q, want %q", got, dir)
		}
	})
}

func TestFindFirstExistingSecure(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "real.json"), []byte(`{"ok":true}`), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("finds the first existing candidate", func(t *testing.T) {
		path, data, err := FindFirstExistingSecure(dir, "missing.json", "real.json", "also-missing.json")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if filepath.Base(path) != "real.json" {
			t.Errorf("path = %q, want basename real.json", path)
		}
		if string(data) != `{"ok":true}` {
			t.Errorf("data = %q, want %q", data, `{"ok":true}`)
		}
	})

	t.Run("returns ErrNotExist when nothing matches", func(t *testing.T) {
		_, _, err := FindFirstExistingSecure(dir, "missing.json", "also-missing.json")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("err = %v, want os.ErrNotExist", err)
		}
	})

	t.Run("skips candidates that escape root instead of erroring", func(t *testing.T) {
		path, data, err := FindFirstExistingSecure(dir, "../../etc/passwd", "real.json")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if filepath.Base(path) != "real.json" {
			t.Errorf("path = %q, want basename real.json (traversal candidate should be skipped)", path)
		}
		if string(data) != `{"ok":true}` {
			t.Errorf("data = %q, want %q", data, `{"ok":true}`)
		}
	})

	t.Run("never escapes root even if every safe candidate is missing", func(t *testing.T) {
		_, _, err := FindFirstExistingSecure(dir, "../../etc/passwd")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("err = %v, want os.ErrNotExist", err)
		}
	})
}
