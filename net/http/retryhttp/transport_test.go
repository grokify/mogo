package retryhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	rt := New()

	if rt.MaxRetries != 3 {
		t.Errorf("expected MaxRetries=3, got %d", rt.MaxRetries)
	}
	if rt.InitialBackoff != 1*time.Second {
		t.Errorf("expected InitialBackoff=1s, got %v", rt.InitialBackoff)
	}
	if rt.MaxBackoff != 30*time.Second {
		t.Errorf("expected MaxBackoff=30s, got %v", rt.MaxBackoff)
	}
	if rt.BackoffMultiplier != 2.0 {
		t.Errorf("expected BackoffMultiplier=2.0, got %f", rt.BackoffMultiplier)
	}
	if rt.Jitter != 0.1 {
		t.Errorf("expected Jitter=0.1, got %f", rt.Jitter)
	}
	if len(rt.RetryableStatusCodes) != 5 {
		t.Errorf("expected 5 retryable status codes, got %d", len(rt.RetryableStatusCodes))
	}
}

func TestNewWithOptions(t *testing.T) {
	rt := NewWithOptions(
		WithMaxRetries(5),
		WithInitialBackoff(500*time.Millisecond),
		WithMaxBackoff(10*time.Second),
		WithBackoffMultiplier(1.5),
		WithJitter(0.2),
	)

	if rt.MaxRetries != 5 {
		t.Errorf("expected MaxRetries=5, got %d", rt.MaxRetries)
	}
	if rt.InitialBackoff != 500*time.Millisecond {
		t.Errorf("expected InitialBackoff=500ms, got %v", rt.InitialBackoff)
	}
	if rt.MaxBackoff != 10*time.Second {
		t.Errorf("expected MaxBackoff=10s, got %v", rt.MaxBackoff)
	}
	if rt.BackoffMultiplier != 1.5 {
		t.Errorf("expected BackoffMultiplier=1.5, got %f", rt.BackoffMultiplier)
	}
	if rt.Jitter != 0.2 {
		t.Errorf("expected Jitter=0.2, got %f", rt.Jitter)
	}
}

func TestRetryOn429(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("success")); err != nil {
			t.Logf("write error: %v", err)
		}
	}))
	defer server.Close()

	rt := NewWithOptions(
		WithMaxRetries(5),
		WithInitialBackoff(10*time.Millisecond),
		WithJitter(0),
	)

	client := rt.Client()
	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryOn500(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	rt := NewWithOptions(
		WithMaxRetries(3),
		WithInitialBackoff(10*time.Millisecond),
		WithJitter(0),
	)

	client := rt.Client()
	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if atomic.LoadInt32(&attempts) != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestNoRetryOnSuccess(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	rt := NewWithOptions(
		WithMaxRetries(3),
		WithInitialBackoff(10*time.Millisecond),
	)

	client := rt.Client()
	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestNoRetryOn400(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	rt := NewWithOptions(
		WithMaxRetries(3),
		WithInitialBackoff(10*time.Millisecond),
	)

	client := rt.Client()
	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("expected 1 attempt (no retry on 400), got %d", attempts)
	}
}

func TestMaxRetriesExhausted(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	rt := NewWithOptions(
		WithMaxRetries(2),
		WithInitialBackoff(10*time.Millisecond),
		WithJitter(0),
	)

	client := rt.Client()
	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should return the last response after exhausting retries
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", resp.StatusCode)
	}

	// Initial attempt + 2 retries = 3 total
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryAfterHeader(t *testing.T) {
	var attempts int32
	var lastBackoff time.Duration

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count < 2 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	rt := NewWithOptions(
		WithMaxRetries(3),
		WithInitialBackoff(10*time.Millisecond),
		WithJitter(0),
		WithOnRetry(func(attempt int, req *http.Request, resp *http.Response, err error, backoff time.Duration) {
			lastBackoff = backoff
		}),
	)

	client := rt.Client()
	start := time.Now()
	resp, err := client.Get(server.URL)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should have respected Retry-After: 1 second
	if lastBackoff != 1*time.Second {
		t.Errorf("expected backoff of 1s from Retry-After, got %v", lastBackoff)
	}

	// Elapsed time should be at least 1 second
	if elapsed < 900*time.Millisecond {
		t.Errorf("expected at least ~1s delay, got %v", elapsed)
	}
}

func TestContextCancellation(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	rt := NewWithOptions(
		WithMaxRetries(10),
		WithInitialBackoff(1*time.Second),
		WithJitter(0),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	client := rt.Client()
	_, err := client.Do(req)

	if err == nil {
		t.Error("expected context error, got nil")
	}

	// Should have only made 1-2 attempts before context cancelled
	if atomic.LoadInt32(&attempts) > 2 {
		t.Errorf("expected 1-2 attempts before cancellation, got %d", attempts)
	}
}

func TestOnRetryCallback(t *testing.T) {
	var callbackCount int
	var lastAttempt int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if callbackCount < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	rt := NewWithOptions(
		WithMaxRetries(5),
		WithInitialBackoff(10*time.Millisecond),
		WithJitter(0),
		WithOnRetry(func(attempt int, req *http.Request, resp *http.Response, err error, backoff time.Duration) {
			callbackCount++
			lastAttempt = attempt
		}),
	)

	client := rt.Client()
	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if callbackCount != 2 {
		t.Errorf("expected 2 callback invocations, got %d", callbackCount)
	}

	if lastAttempt != 2 {
		t.Errorf("expected last attempt=2, got %d", lastAttempt)
	}
}

func TestCustomShouldRetry(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count < 3 {
			w.WriteHeader(http.StatusBadRequest) // Normally not retried
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	rt := NewWithOptions(
		WithMaxRetries(5),
		WithInitialBackoff(10*time.Millisecond),
		WithJitter(0),
		WithShouldRetry(func(resp *http.Response, err error) bool {
			// Custom: retry on 400
			return resp != nil && resp.StatusCode == http.StatusBadRequest
		}),
	)

	client := rt.Client()
	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestWrapClient(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	originalClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	wrappedClient := WrapClient(originalClient,
		WithMaxRetries(3),
		WithInitialBackoff(10*time.Millisecond),
		WithJitter(0),
	)

	// Should preserve timeout
	if wrappedClient.Timeout != originalClient.Timeout {
		t.Errorf("expected timeout to be preserved")
	}

	resp, err := wrappedClient.Get(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if atomic.LoadInt32(&attempts) != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestExponentialBackoff(t *testing.T) {
	rt := NewWithOptions(
		WithInitialBackoff(100*time.Millisecond),
		WithMaxBackoff(10*time.Second),
		WithBackoffMultiplier(2.0),
		WithJitter(0),
	)

	// Test backoff calculation
	backoff0 := rt.calculateBackoff(0, nil)
	backoff1 := rt.calculateBackoff(1, nil)
	backoff2 := rt.calculateBackoff(2, nil)

	if backoff0 != 100*time.Millisecond {
		t.Errorf("expected backoff0=100ms, got %v", backoff0)
	}
	if backoff1 != 200*time.Millisecond {
		t.Errorf("expected backoff1=200ms, got %v", backoff1)
	}
	if backoff2 != 400*time.Millisecond {
		t.Errorf("expected backoff2=400ms, got %v", backoff2)
	}
}

func TestMaxBackoffCap(t *testing.T) {
	rt := NewWithOptions(
		WithInitialBackoff(1*time.Second),
		WithMaxBackoff(5*time.Second),
		WithBackoffMultiplier(10.0),
		WithJitter(0),
	)

	// Attempt 2 would be 1s * 10^2 = 100s, but should be capped at 5s
	backoff := rt.calculateBackoff(2, nil)

	if backoff != 5*time.Second {
		t.Errorf("expected backoff to be capped at 5s, got %v", backoff)
	}
}
