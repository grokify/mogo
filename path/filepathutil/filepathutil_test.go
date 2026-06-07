package filepathutil

import (
	"testing"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "unix path unchanged",
			input:    "/home/user/file.go",
			expected: "/home/user/file.go",
		},
		{
			name:     "windows backslashes to forward",
			input:    `C:\Users\user\file.go`,
			expected: "C:/Users/user/file.go",
		},
		{
			name:     "mixed separators",
			input:    `/home/user\subdir/file.go`,
			expected: "/home/user/subdir/file.go",
		},
		{
			name:     "redundant separators cleaned",
			input:    "/home//user///file.go",
			expected: "/home/user/file.go",
		},
		{
			name:     "dot segments cleaned",
			input:    "/home/user/../other/file.go",
			expected: "/home/other/file.go",
		},
		{
			name:     "empty string",
			input:    "",
			expected: ".",
		},
		{
			name:     "relative path",
			input:    "dir/subdir/file.go",
			expected: "dir/subdir/file.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePath(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePath(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsAbsAny(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		// Unix absolute paths
		{"unix root", "/", true},
		{"unix absolute", "/Users/test/file.txt", true},
		{"unix home", "/home/user/file.txt", true},

		// Windows absolute paths
		{"windows drive backslash", `C:\Users\test\file.txt`, true},
		{"windows drive forward", "C:/Users/test/file.txt", true},
		{"windows lowercase drive", "c:\\Windows\\System32", true},
		{"windows drive only", "D:", true},

		// Relative paths
		{"relative simple", "relative/path.txt", false},
		{"relative dot", "./relative/path.txt", false},
		{"relative dotdot", "../parent/path.txt", false},
		{"filename only", "file.txt", false},

		// Edge cases
		{"empty string", "", false},
		{"single char", "a", false},
		{"colon only", ":", false},
		{"colon in path", "foo:bar", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAbsAny(tt.path)
			if result != tt.expected {
				t.Errorf("IsAbsAny(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}
