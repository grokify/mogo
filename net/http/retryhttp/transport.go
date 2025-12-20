// Package retryhttp provides an HTTP RoundTripper with exponential backoff retry logic.
//
// This package implements retry functionality at the transport level, making it
// compatible with any HTTP client including ogen-generated clients.
package retryhttp

import (
	"bytes"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/grokify/mogo/log/slogutil"
)

// RetryTransport wraps an http.RoundTripper with retry logic.
type RetryTransport struct {
	// Transport is the underlying RoundTripper. If nil, http.DefaultTransport is used.
	Transport http.RoundTripper

	// MaxRetries is the maximum number of retry attempts. Default is 3.
	MaxRetries int

	// InitialBackoff is the initial backoff duration. Default is 1 second.
	InitialBackoff time.Duration

	// MaxBackoff is the maximum backoff duration. Default is 30 seconds.
	MaxBackoff time.Duration

	// BackoffMultiplier is the factor by which backoff increases. Default is 2.0.
	BackoffMultiplier float64

	// Jitter adds randomness to backoff to prevent thundering herd. Default is 0.1 (10%).
	Jitter float64

	// RetryableStatusCodes defines which HTTP status codes trigger a retry.
	// Default is 429, 500, 502, 503, 504.
	RetryableStatusCodes []int

	// ShouldRetry is an optional function for custom retry logic.
	// If set, it's called instead of checking RetryableStatusCodes.
	ShouldRetry func(resp *http.Response, err error) bool

	// OnRetry is an optional callback invoked before each retry attempt.
	OnRetry func(attempt int, req *http.Request, resp *http.Response, err error, backoff time.Duration)

	// Logger is used for logging errors. If nil, a null logger is used.
	Logger *slog.Logger
}

// DefaultRetryableStatusCodes are the status codes that trigger a retry by default.
var DefaultRetryableStatusCodes = []int{
	http.StatusTooManyRequests,     // 429
	http.StatusInternalServerError, // 500
	http.StatusBadGateway,          // 502
	http.StatusServiceUnavailable,  // 503
	http.StatusGatewayTimeout,      // 504
}

// New creates a new RetryTransport with default settings.
func New() *RetryTransport {
	return &RetryTransport{
		Transport:            http.DefaultTransport,
		MaxRetries:           3,
		InitialBackoff:       1 * time.Second,
		MaxBackoff:           30 * time.Second,
		BackoffMultiplier:    2.0,
		Jitter:               0.1,
		RetryableStatusCodes: DefaultRetryableStatusCodes,
	}
}

// Option is a functional option for configuring RetryTransport.
type Option func(*RetryTransport)

// WithTransport sets the underlying transport.
func WithTransport(t http.RoundTripper) Option {
	return func(rt *RetryTransport) {
		rt.Transport = t
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(n int) Option {
	return func(rt *RetryTransport) {
		rt.MaxRetries = n
	}
}

// WithInitialBackoff sets the initial backoff duration.
func WithInitialBackoff(d time.Duration) Option {
	return func(rt *RetryTransport) {
		rt.InitialBackoff = d
	}
}

// WithMaxBackoff sets the maximum backoff duration.
func WithMaxBackoff(d time.Duration) Option {
	return func(rt *RetryTransport) {
		rt.MaxBackoff = d
	}
}

// WithBackoffMultiplier sets the backoff multiplier.
func WithBackoffMultiplier(m float64) Option {
	return func(rt *RetryTransport) {
		rt.BackoffMultiplier = m
	}
}

// WithJitter sets the jitter factor (0.0 to 1.0).
func WithJitter(j float64) Option {
	return func(rt *RetryTransport) {
		rt.Jitter = j
	}
}

// WithRetryableStatusCodes sets the status codes that trigger a retry.
func WithRetryableStatusCodes(codes []int) Option {
	return func(rt *RetryTransport) {
		rt.RetryableStatusCodes = codes
	}
}

// WithShouldRetry sets a custom retry decision function.
func WithShouldRetry(fn func(resp *http.Response, err error) bool) Option {
	return func(rt *RetryTransport) {
		rt.ShouldRetry = fn
	}
}

// WithOnRetry sets a callback invoked before each retry.
func WithOnRetry(fn func(attempt int, req *http.Request, resp *http.Response, err error, backoff time.Duration)) Option {
	return func(rt *RetryTransport) {
		rt.OnRetry = fn
	}
}

// WithLogger sets the logger for error logging.
func WithLogger(l *slog.Logger) Option {
	return func(rt *RetryTransport) {
		rt.Logger = l
	}
}

// NewWithOptions creates a new RetryTransport with the given options.
func NewWithOptions(opts ...Option) *RetryTransport {
	rt := New()
	for _, opt := range opts {
		opt(rt)
	}
	return rt
}

// logger returns the configured logger or a null logger if not set.
func (rt *RetryTransport) logger() *slog.Logger {
	if rt.Logger != nil {
		return rt.Logger
	}
	return slogutil.Null()
}

// RoundTrip implements http.RoundTripper with retry logic.
func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := rt.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// Buffer the request body if present so we can replay it
	var bodyBytes []byte
	if req.Body != nil && req.Body != http.NoBody {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= rt.MaxRetries; attempt++ {
		// Reset body for retry attempts
		if bodyBytes != nil && attempt > 0 {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		// Clone the request to avoid modifying the original
		reqCopy := req.Clone(req.Context())
		if bodyBytes != nil {
			reqCopy.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		resp, err = transport.RoundTrip(reqCopy)

		// Check if we should retry
		if !rt.shouldRetry(resp, err) {
			return resp, err
		}

		// Don't retry if this was the last attempt
		if attempt >= rt.MaxRetries {
			return resp, err
		}

		// Calculate backoff duration
		backoff := rt.calculateBackoff(attempt, resp)

		// Invoke callback if set
		if rt.OnRetry != nil {
			rt.OnRetry(attempt+1, req, resp, err, backoff)
		}

		// Drain and close the response body to allow connection reuse.
		if resp != nil && resp.Body != nil {
			if _, err := io.Copy(io.Discard, resp.Body); err != nil {
				rt.logger().Warn("failed to drain response body", slog.String("error", err.Error()))
			}
			if err := resp.Body.Close(); err != nil {
				rt.logger().Warn("failed to close response body", slog.String("error", err.Error()))
			}
		}

		// Wait for backoff duration or context cancellation
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(backoff):
		}
	}

	return resp, err
}

// shouldRetry determines if a request should be retried.
func (rt *RetryTransport) shouldRetry(resp *http.Response, err error) bool {
	// Use custom function if provided
	if rt.ShouldRetry != nil {
		return rt.ShouldRetry(resp, err)
	}

	// Retry on connection errors
	if err != nil {
		return true
	}

	// Check status code
	if resp != nil {
		for _, code := range rt.RetryableStatusCodes {
			if resp.StatusCode == code {
				return true
			}
		}
	}

	return false
}

// calculateBackoff calculates the backoff duration for a retry attempt.
func (rt *RetryTransport) calculateBackoff(attempt int, resp *http.Response) time.Duration {
	// Check for Retry-After header
	if resp != nil {
		if retryAfter := rt.parseRetryAfter(resp); retryAfter > 0 {
			return retryAfter
		}
	}

	// Calculate exponential backoff
	backoff := float64(rt.InitialBackoff) * math.Pow(rt.BackoffMultiplier, float64(attempt))

	// Apply jitter (math/rand is fine for jitter - no crypto needed)
	if rt.Jitter > 0 {
		jitterRange := backoff * rt.Jitter
		backoff += (rand.Float64()*2 - 1) * jitterRange //nolint:gosec // jitter doesn't need crypto rand
	}

	// Cap at max backoff
	if backoff > float64(rt.MaxBackoff) {
		backoff = float64(rt.MaxBackoff)
	}

	return time.Duration(backoff)
}

// parseRetryAfter parses the Retry-After header value.
// It handles both delay-seconds and HTTP-date formats.
func (rt *RetryTransport) parseRetryAfter(resp *http.Response) time.Duration {
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter == "" {
		return 0
	}

	// Try parsing as seconds
	if seconds, err := strconv.Atoi(retryAfter); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Try parsing as HTTP-date
	if t, err := http.ParseTime(retryAfter); err == nil {
		return time.Until(t)
	}

	return 0
}

// Client returns an *http.Client configured with this RetryTransport.
func (rt *RetryTransport) Client() *http.Client {
	return &http.Client{
		Transport: rt,
	}
}

// WrapClient wraps an existing http.Client with retry logic.
// It returns a new client that uses the original client's transport.
func WrapClient(client *http.Client, opts ...Option) *http.Client {
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	rt := NewWithOptions(append([]Option{WithTransport(transport)}, opts...)...)

	return &http.Client{
		Transport:     rt,
		CheckRedirect: client.CheckRedirect,
		Jar:           client.Jar,
		Timeout:       client.Timeout,
	}
}
