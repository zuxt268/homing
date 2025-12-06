package domain

import "time"

type GoogleBusinesses struct {
	ID        int
	Name      string
	Title     string
	CreatedAt time.Time
}
