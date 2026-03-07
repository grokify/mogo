// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package httputilmore

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLimitRequestBody(t *testing.T) {
	body := strings.NewReader("test body content")
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	w := httptest.NewRecorder()

	LimitRequestBody(w, r, 1024)

	// Verify body was wrapped
	if r.Body == nil {
		t.Fatal("LimitRequestBody set r.Body to nil")
	}

	// Read should work within limit
	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("ReadAll error = %v", err)
	}
	if string(data) != "test body content" {
		t.Errorf("ReadAll = %q, want %q", string(data), "test body content")
	}
}

func TestLimitRequestBody_ExceedsLimit(t *testing.T) {
	// Create body larger than limit
	largeBody := strings.Repeat("x", 100)
	r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(largeBody))
	w := httptest.NewRecorder()

	LimitRequestBody(w, r, 10) // Only allow 10 bytes

	// Reading should fail when exceeding limit
	_, err := io.ReadAll(r.Body)
	if err == nil {
		t.Fatal("ReadAll should have failed for body exceeding limit")
	}
}

func TestLimitRequestBodyDefault(t *testing.T) {
	body := strings.NewReader("test")
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	w := httptest.NewRecorder()

	LimitRequestBodyDefault(w, r)

	// Just verify it doesn't panic and body is still readable
	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("ReadAll error = %v", err)
	}
	if string(data) != "test" {
		t.Errorf("ReadAll = %q, want %q", string(data), "test")
	}
}

func TestLimitRequestBodySmall(t *testing.T) {
	body := strings.NewReader("test")
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	w := httptest.NewRecorder()

	LimitRequestBodySmall(w, r)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("ReadAll error = %v", err)
	}
	if string(data) != "test" {
		t.Errorf("ReadAll = %q, want %q", string(data), "test")
	}
}

func TestLimitRequestBodyLarge(t *testing.T) {
	body := strings.NewReader("test")
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	w := httptest.NewRecorder()

	LimitRequestBodyLarge(w, r)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("ReadAll error = %v", err)
	}
	if string(data) != "test" {
		t.Errorf("ReadAll = %q, want %q", string(data), "test")
	}
}

func TestLimitRequestBodyMultipart(t *testing.T) {
	body := strings.NewReader("test")
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	w := httptest.NewRecorder()

	LimitRequestBodyMultipart(w, r)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("ReadAll error = %v", err)
	}
	if string(data) != "test" {
		t.Errorf("ReadAll = %q, want %q", string(data), "test")
	}
}

func TestParseFormLimited(t *testing.T) {
	body := strings.NewReader("key=value")
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	err := ParseFormLimited(w, r, DefaultMaxBodySize)
	if err != nil {
		t.Fatalf("ParseFormLimited error = %v", err)
	}

	if r.Form.Get("key") != "value" {
		t.Errorf("Form.Get(key) = %q, want %q", r.Form.Get("key"), "value")
	}
}

func TestParseFormLimitedDefault(t *testing.T) {
	body := strings.NewReader("key=value")
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	err := ParseFormLimitedDefault(w, r)
	if err != nil {
		t.Fatalf("ParseFormLimitedDefault error = %v", err)
	}

	if r.Form.Get("key") != "value" {
		t.Errorf("Form.Get(key) = %q, want %q", r.Form.Get("key"), "value")
	}
}

func TestBodySizeConstants(t *testing.T) {
	// Verify constants have expected values
	tests := []struct {
		name  string
		value int64
		want  int64
	}{
		{"DefaultMaxBodySize", DefaultMaxBodySize, 1 << 20},
		{"SmallMaxBodySize", SmallMaxBodySize, 64 << 10},
		{"LargeMaxBodySize", LargeMaxBodySize, 10 << 20},
		{"MultipartMaxBodySize", MultipartMaxBodySize, 32 << 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.want {
				t.Errorf("%s = %d, want %d", tt.name, tt.value, tt.want)
			}
		})
	}
}
