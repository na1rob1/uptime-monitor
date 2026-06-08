package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusRecorderWriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
	rec.WriteHeader(http.StatusCreated)
	if rec.status != http.StatusCreated {
		t.Errorf("status not captured, got %d", rec.status)
	}
	if w.Code != http.StatusCreated {
		t.Errorf("underlying writer not updated, got %d", w.Code)
	}
}

func TestMetricsMiddlewareSkipsMetricsPath(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	mw := MetricsMiddleware(next)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req)
	if !called {
		t.Error("next handler not called")
	}
}

func TestMetricsMiddlewareCountsRequest(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	mw := MetricsMiddleware(next)
	req := httptest.NewRequest(http.MethodPost, "/sites", nil)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req)
	if w.Code != http.StatusTeapot {
		t.Errorf("status not propagated, got %d", w.Code)
	}
}

func TestMetricsMiddlewareDefaultStatus(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mw := MetricsMiddleware(next)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("default 200 expected, got %d", w.Code)
	}
}