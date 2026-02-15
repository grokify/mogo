package osutil

import (
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
