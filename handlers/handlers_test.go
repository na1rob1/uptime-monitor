package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSiteBadJSON(t *testing.T) {
	h := &Handler{repo: nil}
	req := httptest.NewRequest(http.MethodPost, "/sites", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	h.CreateSite(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDeleteSiteBadID(t *testing.T) {
	h := &Handler{repo: nil}
	req := httptest.NewRequest(http.MethodDelete, "/sites?id=abc", nil)
	w := httptest.NewRecorder()
	h.DeleteSite(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}