package models

import "time"

type Site struct {
	ID 		  int 		`json:"id"`
	URL 	  string 	`json:"url"`
	Name	  string 	`json:"name"`
	Status    bool 		`json:"status"`
	Uptime    *float64   `json:"uptime"`
	CheckedAt *time.Time `json:"checked_at"`
	CreatedAt time.Time `json:"created_at"`
}