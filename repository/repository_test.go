package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"uptime-monitor/models"
)

func newMockRepo(t *testing.T) (*Repo, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	return NewRepo(db), mock, func() { db.Close() }
}

func TestNewRepo(t *testing.T) {
	if NewRepo(nil) == nil {
		t.Fatal("nil repo")
	}
}

func TestCreateSuccess(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	mock.ExpectQuery("INSERT INTO sites").
		WithArgs("u", "n").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(7))
	site := &models.Site{URL: "u", Name: "n"}
	if err := r.Create(site); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if site.ID != 7 {
		t.Errorf("expected id=7, got %d", site.ID)
	}
}

func TestCreateError(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	mock.ExpectQuery("INSERT INTO sites").WillReturnError(errors.New("boom"))
	if err := r.Create(&models.Site{}); err == nil {
		t.Error("expected error")
	}
}

func TestUpdateSuccess(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	mock.ExpectExec("UPDATE sites SET name").
		WithArgs("n", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := r.Update(&models.Site{ID: 1, Name: "n"}); err != nil {
		t.Fatalf("Update: %v", err)
	}
}

func TestUpdateError(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	mock.ExpectExec("UPDATE sites SET name").WillReturnError(errors.New("boom"))
	if err := r.Update(&models.Site{}); err == nil {
		t.Error("expected error")
	}
}

func TestReadSuccess(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	rows := sqlmock.NewRows([]string{"id", "url", "name", "status", "uptime", "checked_at", "created_at"}).
		AddRow(1, "u1", "n1", true, nil, nil, time.Now())
	mock.ExpectQuery("SELECT id, url, name").WillReturnRows(rows)
	sites, err := r.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(sites) != 1 || sites[0].ID != 1 {
		t.Errorf("bad result: %+v", sites)
	}
}

func TestReadQueryError(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	mock.ExpectQuery("SELECT id, url, name").WillReturnError(errors.New("boom"))
	if _, err := r.Read(); err == nil {
		t.Error("expected error")
	}
}

func TestReadScanError(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	rows := sqlmock.NewRows([]string{"id"}).AddRow("not-int")
	mock.ExpectQuery("SELECT id, url, name").WillReturnRows(rows)
	if _, err := r.Read(); err == nil {
		t.Error("expected scan error")
	}
}

func TestDeleteSuccess(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	mock.ExpectExec("DELETE FROM sites").WithArgs(3).WillReturnResult(sqlmock.NewResult(0, 1))
	if err := r.Delete(3); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestDeleteError(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	mock.ExpectExec("DELETE FROM sites").WillReturnError(errors.New("boom"))
	if err := r.Delete(1); err == nil {
		t.Error("expected error")
	}
}

func TestUpdateStatusUp(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	mock.ExpectExec("UPDATE sites SET").
		WithArgs(true, 1, 5).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := r.UpdateStatus(5, true); err != nil {
		t.Fatalf("UpdateStatus: %v", err)
	}
}

func TestUpdateStatusDown(t *testing.T) {
	r, mock, done := newMockRepo(t)
	defer done()
	mock.ExpectExec("UPDATE sites SET").
		WithArgs(false, 0, 5).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := r.UpdateStatus(5, false); err != nil {
		t.Fatalf("UpdateStatus: %v", err)
	}
}