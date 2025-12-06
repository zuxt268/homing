package domain

import "time"

type GooglePost struct {
	ID                int
	GoogleBusinessURL string
	InstagramURL      string
	MediaID           string
	CustomerID        int
	Name              string
	MediaFormat       string
	GoogleURL         string
	CreateTime        string
	CreatedAt         time.Time
}
