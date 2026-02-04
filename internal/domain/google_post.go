package domain

import "time"

type GooglePost struct {
	ID           int
	InstagramURL string
	MediaID      string
	CustomerID   int
	Name         string
	GoogleURL    string
	CreateTime   string
	PostType     string
	CreatedAt    time.Time
}

const (
	PostTypePhoto = "photo"
	PostTypePost  = "post"
)
