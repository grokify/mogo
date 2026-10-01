package retryhttp

import (
	"context"
	"sync"
	"time"
)

// IntervalLimiter paces requests so that at most one starts per interval,
// across all goroutines sharing it. It is a minimal, dependency-free Limiter
// for clients that must be polite to a shared public API. For burst-capable
// token buckets use golang.org/x/time/rate.Limiter, which also satisfies
// Limiter.
type IntervalLimiter struct {
	interval time.Duration
	mu       sync.Mutex
	next     time.Time
}

// NewIntervalLimiter returns a limiter allowing one request per interval. A
// non-positive interval disables pacing.
func NewIntervalLimiter(interval time.Duration) *IntervalLimiter {
	return &IntervalLimiter{interval: interval}
}

// PerSecond returns an IntervalLimiter allowing n requests per second. A
// non-positive n disables pacing.
func PerSecond(n float64) *IntervalLimiter {
	if n <= 0 {
		return NewIntervalLimiter(0)
	}
	return NewIntervalLimiter(time.Duration(float64(time.Second) / n))
}

// Wait blocks until this caller's slot arrives or ctx is done. Slots are
// handed out in call order; a waiter that gives up releases nothing, so the
// schedule is never compressed by cancelled callers.
func (l *IntervalLimiter) Wait(ctx context.Context) error {
	if l == nil || l.interval <= 0 {
		return ctx.Err()
	}
	l.mu.Lock()
	now := time.Now()
	slot := l.next
	if slot.Before(now) {
		slot = now
	}
	l.next = slot.Add(l.interval)
	l.mu.Unlock()

	delay := time.Until(slot)
	if delay <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// Interval reports the configured spacing.
func (l *IntervalLimiter) Interval() time.Duration { return l.interval }
