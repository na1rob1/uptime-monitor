package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"uptime-monitor/checker"

	_ "github.com/lib/pq"
	"uptime-monitor/handlers"
	"uptime-monitor/repository"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://nikpopo@localhost:5432/uptime?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to database")

	repo := repository.NewRepo(db)
	h := handlers.NewHandler(repo)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "db down", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	http.HandleFunc("/sites", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.GetSites(w, r)
		case http.MethodPost:
			h.CreateSite(w, r)
		case http.MethodPut:
			h.UpdateSite(w, r)
		case http.MethodDelete:
			h.DeleteSite(w, r)
		}
	})

	go checker.Start(repo)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}