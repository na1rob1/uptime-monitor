package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"uptime-monitor/models"
	"uptime-monitor/repository"
)

type Handler struct {
	repo *repository.Repo
}

func NewHandler(repo *repository.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateSite(w http.ResponseWriter, r *http.Request) {
	var site models.Site
	err := json.NewDecoder(r.Body).Decode(&site)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	err = h.repo.Create(&site)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(site.URL)
	status := err == nil && resp.StatusCode == 200
	err = h.repo.UpdateStatus(site.ID, status)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(site)
	if err != nil {
		return
	}
}

func (h *Handler) GetSites(w http.ResponseWriter, r *http.Request) {
	sites, err := h.repo.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(sites)
	if err != nil {
		return
	}
}

func (h *Handler) UpdateSite(w http.ResponseWriter, r *http.Request) {
	var site models.Site
	err := json.NewDecoder(r.Body).Decode(&site)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	err = h.repo.Update(&site)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteSite(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
