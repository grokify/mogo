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
