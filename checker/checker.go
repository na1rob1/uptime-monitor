package checker

import (
	"net/http"
	"time"
	"uptime-monitor/repository"
)

var client = &http.Client{Timeout: 5 * time.Second}

func checkOnce(repo *repository.Repo) {
	sites, err := repo.Read()
	if err != nil {
		return
	}
	for _, site := range sites {
		resp, err := client.Get(site.URL)
		if err != nil {
			repo.UpdateStatus(site.ID, false)
			continue
		}
		status := resp.StatusCode == 200
		resp.Body.Close()
		repo.UpdateStatus(site.ID, status)
	}
}

func Start(repo *repository.Repo) {
	for range time.NewTicker(5 * time.Second).C { checkOnce(repo) }
}