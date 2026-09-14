package healthz

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/grokify/mogo/net/http/httputilmore"
)

func TestLivenessHandler_NoFields(t *testing.T) {
	h := LivenessHandler(nil)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status field = %v, want ok", body["status"])
	}
}

func TestLivenessHandler_WithFields(t *testing.T) {
	h := LivenessHandler(func() map[string]any {
		return map[string]any{"clients": float64(3)}
	})
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status field = %v, want ok", body["status"])
	}
	if body["clients"] != float64(3) {
		t.Errorf("clients field = %v, want 3", body["clients"])
	}
}

func TestReadinessHandler_AllPass(t *testing.T) {
	h := ReadinessHandler(
		CheckerFunc{CheckerName: "a", Fn: func(context.Context) error { return nil }},
		CheckerFunc{CheckerName: "b", Fn: func(context.Context) error { return nil }},
	)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp readinessResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Status != "ok" || len(resp.Checks) != 2 {
		t.Errorf("resp = %+v, want status=ok with 2 checks", resp)
	}
	for _, c := range resp.Checks {
		if c.Status != "ok" {
			t.Errorf("check %s status = %s, want ok", c.Name, c.Status)
		}
	}
}

func TestReadinessHandler_OneFails(t *testing.T) {
	h := ReadinessHandler(
		CheckerFunc{CheckerName: "good", Fn: func(context.Context) error { return nil }},
		CheckerFunc{CheckerName: "bad", Fn: func(context.Context) error { return errors.New("boom") }},
	)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
	var resp readinessResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Status != "degraded" {
		t.Errorf("resp.Status = %s, want degraded", resp.Status)
	}
	var found bool
	for _, c := range resp.Checks {
		if c.Name == "bad" {
			found = true
			if c.Status != "fail" || c.Error != "boom" {
				t.Errorf("bad check = %+v, want status=fail error=boom", c)
			}
		}
	}
	if !found {
		t.Error("bad check missing from response")
	}
}

func TestReadinessHandler_NoChecks(t *testing.T) {
	h := ReadinessHandler()
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (vacuously ready)", rec.Code)
	}
}

func TestProbe_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if err := Probe(context.Background(), srv.URL+"/health", time.Second); err != nil {
		t.Errorf("Probe() = %v, want nil", err)
	}
}

func TestNormalizeLoopback(t *testing.T) {
	tests := []struct{ in, want string }{
		{"0.0.0.0:8080", "127.0.0.1:8080"},
		{"[::]:8080", "127.0.0.1:8080"},
		{":8080", "127.0.0.1:8080"},
		{"192.168.1.5:8080", "192.168.1.5:8080"},
		{"not-a-host-port", "not-a-host-port"},
	}
	for _, tt := range tests {
		if got := NormalizeLoopback(tt.in); got != tt.want {
			t.Errorf("NormalizeLoopback(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestProbe_BindAllNormalization(t *testing.T) {
	// G102: deliberately binding all interfaces is the exact case under
	// test — Probe must normalize this address to 127.0.0.1 to dial it.
	ln, err := net.Listen("tcp", "0.0.0.0:0") //nolint:gosec // G102: intentional bind-all, this test covers Probe's normalization of it
	if err != nil {
		t.Skipf("cannot bind 0.0.0.0: %v", err)
	}
	defer ln.Close()

	srv := httputilmore.NewServerTimeouts("", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), 5*time.Second)
	go func() { _ = srv.Serve(ln) }()
	defer srv.Close()

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("SplitHostPort: %v", err)
	}

	// A bind-all address is not itself dialable — Probe must normalize it
	// to 127.0.0.1 rather than trying (and failing) to dial 0.0.0.0.
	addr := net.JoinHostPort("0.0.0.0", port)
	url := "http://" + NormalizeLoopback(addr) + "/health"
	if err := Probe(context.Background(), url, time.Second); err != nil {
		t.Errorf("Probe(%q) = %v, want nil after loopback normalization", url, err)
	}
}

func TestProbe_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	if err := Probe(context.Background(), srv.URL+"/health", time.Second); err == nil {
		t.Error("Probe() = nil, want error for non-200 status")
	}
}

func TestProbe_ConnectionRefused(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close() // free the port immediately so the dial fails

	if err := Probe(context.Background(), "http://"+addr+"/health", 500*time.Millisecond); err == nil {
		t.Error("Probe() = nil, want error when nothing is listening")
	}
}

func TestProbe_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if err := Probe(context.Background(), srv.URL+"/health", 20*time.Millisecond); err == nil {
		t.Error("Probe() = nil, want timeout error")
	}
}
