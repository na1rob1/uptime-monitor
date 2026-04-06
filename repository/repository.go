package repository

import (
	"database/sql"
	"uptime-monitor/models"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(site *models.Site) error {
	return r.db.QueryRow(
		"INSERT INTO sites (url, name) VALUES ($1, $2) RETURNING id",
		site.URL, site.Name).Scan(&site.ID)
}

func (r *Repo) Update(site *models.Site) error {
	_, err := r.db.Exec("UPDATE sites SET url = $1, name = $2 WHERE id = $3",
	site.URL, site.Name, site.ID)
	return err
}

func (r *Repo) Read() ([]models.Site, error) {
	rows, err := r.db.Query("SELECT id, url, name, status, uptime, checked_at, created_at FROM sites")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []models.Site
	for rows.Next() {
		var s models.Site
		err := rows.Scan(&s.ID, &s.URL, &s.Name, &s.Status, &s.Uptime, &s.CheckedAt, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, nil
}

func (r *Repo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM sites WHERE id = $1", id)
	return err
}

//dlya checker
func (r *Repo) UpdateStatus(id int, status bool) error {
	upInc := 0
	if status {
		upInc = 1
	}
	_, err := r.db.Exec(`
          UPDATE sites SET
              status = $1,
              checked_at = NOW(),
              total_checks = total_checks + 1,
              up_checks = up_checks + $2,
              uptime = (up_checks + $2)::REAL / (total_checks + 1) * 100
          WHERE id = $3`,
		status, upInc, id,
	)
	return err
}
