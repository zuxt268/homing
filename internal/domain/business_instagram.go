package domain

import "time"

type BusinessInstagram struct {
	ID            int
	Name          string
	Memo          string
	InstagramID   string
	InstagramName string
	BusinessName  string
	BusinessTitle string
	StartDate     time.Time
	Status        Status
	UpdatedAt     time.Time
	CreatedAt     time.Time
}
