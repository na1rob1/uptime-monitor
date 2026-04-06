package checker

import (
	"net/http"
	"time"
	"uptime-monitor/repository"
)

func Start(repo *repository.Repo) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		sites, err := repo.Read()
		if err != nil {
			continue
		}
		for _, site := range sites {
			resp, err := http.Get(site.URL)
			status := err == nil && resp.StatusCode == 200
			repo.UpdateStatus(site.ID, status)
		}
	}
}