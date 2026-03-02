package sanitize

import (
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "clean string",
			input: "hello world",
			want:  "hello world",
		},
		{
			name:  "newline injection",
			input: "user123\nERROR: fake log entry",
			want:  "user123ERROR: fake log entry",
		},
		{
			name:  "carriage return injection",
			input: "session\r\nforged",
			want:  "sessionforged",
		},
		{
			name:  "tab character",
			input: "col1\tcol2",
			want:  "col1col2",
		},
		{
			name:  "null byte",
			input: "before\x00after",
			want:  "beforeafter",
		},
		{
			name:  "bell character",
			input: "alert\x07here",
			want:  "alerthere",
		},
		{
			name:  "delete character",
			input: "test\x7Fvalue",
			want:  "testvalue",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only control chars",
			input: "\n\r\t",
			want:  "",
		},
		{
			name:  "unicode preserved",
			input: "日本語テスト",
			want:  "日本語テスト",
		},
		{
			name:  "mixed unicode and control",
			input: "hello\n世界\rtest",
			want:  "hello世界test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := String(tt.input)
			if got != tt.want {
				t.Errorf("String(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStringReplace(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		replacement string
		want        string
	}{
		{
			name:        "replace with question mark",
			input:       "line1\nline2",
			replacement: "?",
			want:        "line1?line2",
		},
		{
			name:        "replace with unicode",
			input:       "bad\x00data",
			replacement: "�",
			want:        "bad�data",
		},
		{
			name:        "replace with empty",
			input:       "test\tvalue",
			replacement: "",
			want:        "testvalue",
		},
		{
			name:        "replace multiple",
			input:       "\n\r\t",
			replacement: "_",
			want:        "___",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringReplace(tt.input, tt.replacement)
			if got != tt.want {
				t.Errorf("StringReplace(%q, %q) = %q, want %q", tt.input, tt.replacement, got, tt.want)
			}
		})
	}
}

func TestStrings(t *testing.T) {
	input := []string{"clean", "with\nnewline", "also\ttab"}
	want := []string{"clean", "withnewline", "alsotab"}

	got := Strings(input...)

	if len(got) != len(want) {
		t.Fatalf("Strings() returned %d elements, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Strings()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestStringOrTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "no truncation needed",
			input:  "short",
			maxLen: 10,
			want:   "short",
		},
		{
			name:   "exact length",
			input:  "exact",
			maxLen: 5,
			want:   "exact",
		},
		{
			name:   "truncated with ellipsis",
			input:  "this is a long string",
			maxLen: 10,
			want:   "this is...",
		},
		{
			name:   "sanitize and truncate",
			input:  "bad\ndata is here",
			maxLen: 10,
			want:   "baddata...",
		},
		{
			name:   "very short maxLen",
			input:  "hello",
			maxLen: 2,
			want:   "he",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringOrTruncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("StringOrTruncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestIsClean(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"clean string", "hello world 123", true},
		{"with newline", "hello\nworld", false},
		{"with tab", "hello\tworld", false},
		{"with null", "hello\x00world", false},
		{"empty string", "", true},
		{"unicode clean", "日本語", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsClean(tt.input)
			if got != tt.want {
				t.Errorf("IsClean(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestHasControlChars(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"clean ASCII", "hello", false},
		{"with newline", "hello\n", true},
		{"clean unicode", "日本語", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasControlChars(tt.input)
			if got != tt.want {
				t.Errorf("HasControlChars(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStripAllControl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"clean", "hello", "hello"},
		{"newline", "hello\nworld", "helloworld"},
		{"unicode preserved", "日本語", "日本語"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripAllControl(tt.input)
			if got != tt.want {
				t.Errorf("StripAllControl(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// Benchmark to ensure sanitization is fast enough for high-volume logging
func BenchmarkString(b *testing.B) {
	input := "session_id_12345\ninjected_line\rmore_data"
	for i := 0; i < b.N; i++ {
		_ = String(input)
	}
}

func BenchmarkIsClean(b *testing.B) {
	clean := "session_id_12345_no_injection_here"
	for i := 0; i < b.N; i++ {
		_ = IsClean(clean)
	}
}
