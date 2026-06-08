package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"uptime-monitor/repository"
)

func newTestHandler(t *testing.T) (*Handler, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	h := NewHandler(repository.NewRepo(db))
	return h, mock, func() { db.Close() }
}

func TestNewHandler(t *testing.T) {
	h := NewHandler(nil)
	if h == nil {
		t.Fatal("nil handler")
	}
}

func TestCreateSiteBadJSON(t *testing.T) {
	h := &Handler{repo: nil}
	req := httptest.NewRequest(http.MethodPost, "/sites", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	h.CreateSite(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateSiteRepoError(t *testing.T) {
	h, mock, done := newTestHandler(t)
	defer done()
	mock.ExpectQuery("INSERT INTO sites").
		WithArgs("http://example.com", "ex").
		WillReturnError(errors.New("db fail"))
	body, _ := json.Marshal(map[string]string{"url": "http://example.com", "name": "ex"})
	req := httptest.NewRequest(http.MethodPost, "/sites", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateSite(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestCreateSiteSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h, mock, done := newTestHandler(t)
	defer done()
	mock.ExpectQuery("INSERT INTO sites").
		WithArgs(srv.URL, "ex").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("UPDATE sites SET").
		WithArgs(true, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body, _ := json.Marshal(map[string]string{"url": srv.URL, "name": "ex"})
	req := httptest.NewRequest(http.MethodPost, "/sites", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateSite(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestGetSitesSuccess(t *testing.T) {
	h, mock, done := newTestHandler(t)
	defer done()
	mock.ExpectQuery("SELECT id, url, name").
		WillReturnRows(sqlmock.NewRows([]string{"id", "url", "name", "status", "uptime", "checked_at", "created_at"}))
	req := httptest.NewRequest(http.MethodGet, "/sites", nil)
	w := httptest.NewRecorder()
	h.GetSites(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetSitesRepoError(t *testing.T) {
	h, mock, done := newTestHandler(t)
	defer done()
	mock.ExpectQuery("SELECT id, url, name").WillReturnError(errors.New("db fail"))
	req := httptest.NewRequest(http.MethodGet, "/sites", nil)
	w := httptest.NewRecorder()
	h.GetSites(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestUpdateSiteBadJSON(t *testing.T) {
	h := &Handler{repo: nil}
	req := httptest.NewRequest(http.MethodPut, "/sites", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	h.UpdateSite(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateSiteSuccess(t *testing.T) {
	h, mock, done := newTestHandler(t)
	defer done()
	mock.ExpectExec("UPDATE sites SET name").
		WithArgs("newname", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	body, _ := json.Marshal(map[string]interface{}{"id": 1, "name": "newname"})
	req := httptest.NewRequest(http.MethodPut, "/sites", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.UpdateSite(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUpdateSiteRepoError(t *testing.T) {
	h, mock, done := newTestHandler(t)
	defer done()
	mock.ExpectExec("UPDATE sites SET name").WillReturnError(errors.New("db fail"))
	body, _ := json.Marshal(map[string]interface{}{"id": 1, "name": "x"})
	req := httptest.NewRequest(http.MethodPut, "/sites", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.UpdateSite(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
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

func TestDeleteSiteSuccess(t *testing.T) {
	h, mock, done := newTestHandler(t)
	defer done()
	mock.ExpectExec("DELETE FROM sites").WithArgs(5).WillReturnResult(sqlmock.NewResult(0, 1))
	req := httptest.NewRequest(http.MethodDelete, "/sites?id=5", nil)
	w := httptest.NewRecorder()
	h.DeleteSite(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDeleteSiteRepoError(t *testing.T) {
	h, mock, done := newTestHandler(t)
	defer done()
	mock.ExpectExec("DELETE FROM sites").WithArgs(5).WillReturnError(errors.New("db fail"))
	req := httptest.NewRequest(http.MethodDelete, "/sites?id=5", nil)
	w := httptest.NewRecorder()
	h.DeleteSite(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}