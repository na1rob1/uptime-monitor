package checker

import (
	"net/http"
	"time"
	"uptime-monitor/repository"
)

func Start(repo *repository.Repo) {
	client := &http.Client{Timeout: 5 * time.Second}
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		sites, err := repo.Read()
		if err != nil {
			continue
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
}