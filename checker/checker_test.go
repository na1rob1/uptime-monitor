package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"uptime-monitor/repository"
)

func TestCheckOnceReadError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT id, url, name").WillReturnError(http.ErrServerClosed)
	checkOnce(repository.NewRepo(db))
}

func TestCheckOnceSiteUp(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "url", "name", "status", "uptime", "checked_at", "created_at"}).
		AddRow(1, srv.URL, "n", true, nil, nil, time.Now())
	mock.ExpectQuery("SELECT id, url, name").WillReturnRows(rows)
	mock.ExpectExec("UPDATE sites SET").
		WithArgs(true, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	checkOnce(repository.NewRepo(db))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCheckOnceSiteDown(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "url", "name", "status", "uptime", "checked_at", "created_at"}).
		AddRow(1, "http://127.0.0.1:1", "n", true, nil, nil, time.Now())
	mock.ExpectQuery("SELECT id, url, name").WillReturnRows(rows)
	mock.ExpectExec("UPDATE sites SET").
		WithArgs(false, 0, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	checkOnce(repository.NewRepo(db))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCheckOnceSiteNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "url", "name", "status", "uptime", "checked_at", "created_at"}).
		AddRow(2, srv.URL, "n2", false, nil, nil, time.Now())
	mock.ExpectQuery("SELECT id, url, name").WillReturnRows(rows)
	mock.ExpectExec("UPDATE sites SET").
		WithArgs(false, 0, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	checkOnce(repository.NewRepo(db))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}