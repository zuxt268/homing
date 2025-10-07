package domain

import "time"

type Post struct {
	ID           int
	WordpressURL string
	InstagramURL string
	CreatedAt    time.Time
}
