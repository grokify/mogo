// Package healthz provides standard liveness/readiness HTTP handlers and a
// dependency-free client-side probe, for services that need a consistent
// health contract without pulling in a framework.
//
// Liveness answers "is the process up and serving HTTP" — it should never
// depend on anything external and should stay cheap, since orchestrators
// poll it every few seconds. Readiness answers "can this instance actually
// do its job right now" by running a set of Checkers; a failing checker
// takes the instance out of a load balancer's rotation without killing the
// process. See LivenessHandler and ReadinessHandler.
//
// Probe is the client half: a single dependency-free HTTP GET, intended for
// container HEALTHCHECK exec-form on shell-less runtime images that have no
// wget/curl to shell out to (e.g. Chainguard/distroless static images).
package healthz

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// Checker is a named readiness check. Implementations should return
// promptly and respect ctx's deadline — ReadinessHandler runs checkers
// concurrently, bounded by the incoming request's context.
type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

// CheckerFunc adapts a plain function into a Checker, for callers that
// don't want to define a named type per check.
type CheckerFunc struct {
	CheckerName string
	Fn          func(ctx context.Context) error
}

// Name implements Checker.
func (c CheckerFunc) Name() string { return c.CheckerName }

// Check implements Checker.
func (c CheckerFunc) Check(ctx context.Context) error { return c.Fn(ctx) }

// checkResult is one checker's outcome in the readiness JSON response.
type checkResult struct {
	Name   string `json:"name"`
	Status string `json:"status"` // "ok" or "fail"
	Error  string `json:"error,omitempty"`
}

// readinessResponse is the readiness endpoint's JSON body.
type readinessResponse struct {
	Status string        `json:"status"` // "ok" or "degraded"
	Checks []checkResult `json:"checks"`
}

// LivenessHandler returns 200 with a JSON body while the process is
// serving. fields, if non-nil, is called on every request to add extra
// top-level fields (e.g. a connected-client count) — it must return a map
// safe to marshal. Pass nil for a bare {"status":"ok"} body.
func LivenessHandler(fields func() map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		body := map[string]any{"status": "ok"}
		if fields != nil {
			for k, v := range fields() {
				body[k] = v
			}
		}
		writeJSON(w, http.StatusOK, body)
	}
}

// ReadinessHandler runs every checker concurrently, each bounded by the
// incoming request's context, and returns 200 with per-check detail if all
// pass, 503 if any fail. Checks run concurrently so one slow dependency
// doesn't serialize total probe latency. Zero checkers is vacuously ready
// (200, empty checks list) — useful before any real checker exists yet.
func ReadinessHandler(checks ...Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results := make([]checkResult, len(checks))
		var wg sync.WaitGroup
		for i, c := range checks {
			wg.Add(1)
			go func(i int, c Checker) {
				defer wg.Done()
				if err := c.Check(r.Context()); err != nil {
					results[i] = checkResult{Name: c.Name(), Status: "fail", Error: err.Error()}
				} else {
					results[i] = checkResult{Name: c.Name(), Status: "ok"}
				}
			}(i, c)
		}
		wg.Wait()

		resp := readinessResponse{Status: "ok", Checks: results}
		status := http.StatusOK
		for _, res := range results {
			if res.Status != "ok" {
				resp.Status = "degraded"
				status = http.StatusServiceUnavailable
				break
			}
		}
		writeJSON(w, status, resp)
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// NormalizeLoopback rewrites a bind-all host (0.0.0.0, ::, or empty) in
// addr (a host:port pair) to 127.0.0.1, since such an address is valid to
// listen on but not to dial. addr is returned unchanged if it already
// names a real host, or if it doesn't parse as host:port at all.
func NormalizeLoopback(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	if host == "0.0.0.0" || host == "::" || host == "" {
		return net.JoinHostPort("127.0.0.1", port)
	}
	return addr
}

// Probe performs a single GET against rawURL and returns nil only on a
// 200 response. timeout bounds both connection and response. rawURL must
// be a complete URL with scheme — for a container's own loopback address
// build it with NormalizeLoopback first, e.g.
// "http://" + NormalizeLoopback(addr) + "/health".
//
// This is the client half of the liveness/readiness contract: usable both
// as a dependency-free container HEALTHCHECK exec-form probe on
// shell-less runtime images (no wget/curl/shell to run an HTTP probe
// with), and as a post-deploy verification step against a real deployed
// endpoint.
func Probe(ctx context.Context, rawURL string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf("healthz: build request for %s: %w", rawURL, err)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("healthz: probe %s: %w", rawURL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("healthz: probe %s: %s", rawURL, resp.Status)
	}
	return nil
}
