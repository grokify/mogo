// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package httputilmore

import "net/http"

// Body size constants for common use cases.
const (
	// DefaultMaxBodySize is the default maximum body size (1MB).
	// Suitable for most form submissions and JSON payloads.
	DefaultMaxBodySize int64 = 1 << 20 // 1MB

	// SmallMaxBodySize is for endpoints expecting small payloads (64KB).
	// Suitable for login forms, simple API calls.
	SmallMaxBodySize int64 = 64 << 10 // 64KB

	// LargeMaxBodySize is for endpoints expecting larger payloads (10MB).
	// Suitable for file uploads with known size limits.
	LargeMaxBodySize int64 = 10 << 20 // 10MB

	// MultipartMaxBodySize is for multipart form uploads (32MB).
	// Matches http.defaultMaxMemory.
	MultipartMaxBodySize int64 = 32 << 20 // 32MB
)

// LimitRequestBody wraps r.Body with http.MaxBytesReader to prevent memory
// exhaustion attacks. Call this before ParseForm, ParseMultipartForm, or
// any operation that reads the request body.
//
// This function addresses gosec rule G120:
// "Parsing form data without limiting request body size can allow memory exhaustion"
//
// Example:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    httputilmore.LimitRequestBody(w, r, httputilmore.DefaultMaxBodySize)
//	    if err := r.ParseForm(); err != nil {
//	        http.Error(w, "Bad request", http.StatusBadRequest)
//	        return
//	    }
//	    // ... handle form data
//	}
//
// Note: This modifies r.Body in place. The original body is replaced with
// a limited reader that returns an error if the limit is exceeded.
func LimitRequestBody(w http.ResponseWriter, r *http.Request, maxBytes int64) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
}

// LimitRequestBodyDefault calls LimitRequestBody with DefaultMaxBodySize (1MB).
//
// This is a convenience function for the common case where 1MB is sufficient.
func LimitRequestBodyDefault(w http.ResponseWriter, r *http.Request) {
	LimitRequestBody(w, r, DefaultMaxBodySize)
}

// LimitRequestBodySmall calls LimitRequestBody with SmallMaxBodySize (64KB).
//
// Use this for endpoints that expect small payloads like login forms.
func LimitRequestBodySmall(w http.ResponseWriter, r *http.Request) {
	LimitRequestBody(w, r, SmallMaxBodySize)
}

// LimitRequestBodyLarge calls LimitRequestBody with LargeMaxBodySize (10MB).
//
// Use this for endpoints that accept file uploads or larger payloads.
func LimitRequestBodyLarge(w http.ResponseWriter, r *http.Request) {
	LimitRequestBody(w, r, LargeMaxBodySize)
}

// LimitRequestBodyMultipart calls LimitRequestBody with MultipartMaxBodySize (32MB).
//
// Use this for multipart form uploads. This matches the default memory limit
// used by http.Request.ParseMultipartForm when called with 0.
func LimitRequestBodyMultipart(w http.ResponseWriter, r *http.Request) {
	LimitRequestBody(w, r, MultipartMaxBodySize)
}

// ParseFormLimited is a convenience function that limits the body size and
// parses the form in one call. It returns any error from ParseForm.
//
// Example:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    if err := httputilmore.ParseFormLimited(w, r, httputilmore.DefaultMaxBodySize); err != nil {
//	        http.Error(w, "Bad request", http.StatusBadRequest)
//	        return
//	    }
//	    // ... handle form data
//	}
func ParseFormLimited(w http.ResponseWriter, r *http.Request, maxBytes int64) error {
	LimitRequestBody(w, r, maxBytes)
	return r.ParseForm()
}

// ParseFormLimitedDefault calls ParseFormLimited with DefaultMaxBodySize.
func ParseFormLimitedDefault(w http.ResponseWriter, r *http.Request) error {
	return ParseFormLimited(w, r, DefaultMaxBodySize)
}
