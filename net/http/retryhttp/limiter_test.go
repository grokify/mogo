package retryhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestIntervalLimiterSpacesRequests(t *testing.T) {
	const interval = 30 * time.Millisecond
	l := NewIntervalLimiter(interval)
	start := time.Now()
	for i := 0; i < 4; i++ {
		if err := l.Wait(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	// First call is immediate; three more are spaced by interval.
	if elapsed := time.Since(start); elapsed < 3*interval {
		t.Fatalf("4 waits took %v, want at least %v", elapsed, 3*interval)
	}
}

func TestIntervalLimiterConcurrentCallersAreSerialized(t *testing.T) {
	const interval = 20 * time.Millisecond
	l := NewIntervalLimiter(interval)
	var mu sync.Mutex
	var starts []time.Time
	var wg sync.WaitGroup
	begin := time.Now()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.Wait(context.Background()); err != nil {
				t.Error(err)
				return
			}
			mu.Lock()
			starts = append(starts, time.Now())
			mu.Unlock()
		}()
	}
	wg.Wait()
	if len(starts) != 5 {
		t.Fatalf("got %d starts", len(starts))
	}
	slices.SortFunc(starts, time.Time.Compare)
	// Slots are interval apart, so the k-th caller to return must not do so
	// before begin + k*interval. Measure from begin rather than between
	// consecutive returns: goroutine wake-up latency can delay any return
	// (shrinking the gap to the next one) but can never make one early.
	for k, s := range starts {
		if got, want := s.Sub(begin), time.Duration(k)*interval; got < want {
			t.Errorf("caller %d returned after %v, want at least %v", k, got, want)
		}
	}
}

func TestIntervalLimiterRespectsContext(t *testing.T) {
	l := NewIntervalLimiter(time.Second)
	if err := l.Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := l.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error")
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Fatal("wait did not return promptly on cancellation")
	}
}

func TestIntervalLimiterDisabled(t *testing.T) {
	var nilLimiter *IntervalLimiter
	if err := nilLimiter.Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	l := PerSecond(0)
	start := time.Now()
	for i := 0; i < 100; i++ {
		if err := l.Wait(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	if time.Since(start) > 100*time.Millisecond {
		t.Fatal("disabled limiter should not block")
	}
	if PerSecond(4).Interval() != 250*time.Millisecond {
		t.Fatalf("PerSecond(4) interval = %v", PerSecond(4).Interval())
	}
}

type countingLimiter struct{ n int32 }

func (c *countingLimiter) Wait(ctx context.Context) error {
	atomic.AddInt32(&c.n, 1)
	return ctx.Err()
}

func TestWithLimiterAppliesToEveryAttempt(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&attempts, 1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	lim := &countingLimiter{}
	rt := NewWithOptions(
		WithMaxRetries(5),
		WithInitialBackoff(5*time.Millisecond),
		WithJitter(0),
		WithLimiter(lim),
	)
	resp, err := rt.Client().Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if got := atomic.LoadInt32(&lim.n); got != 3 {
		t.Fatalf("limiter consulted %d times, want 3 (one per attempt)", got)
	}
}

func TestWithLimiterContextCancelledBeforeAttempt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	rt := NewWithOptions(WithLimiter(NewIntervalLimiter(time.Hour)))
	client := rt.Client()
	// Consume the immediate slot.
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Do(req); err == nil {
		t.Fatal("expected error when the limiter wait is cancelled")
	}
}
