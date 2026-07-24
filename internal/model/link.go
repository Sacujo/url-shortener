package model

import "time"

type Link struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}
