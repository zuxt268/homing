package domain

import "time"

type WordpressGbp struct {
	ID              int
	Name            string
	Memo            string
	WordpressDomain string
	BusinessName    string
	BusinessTitle   string
	MapsURL         string
	StartDate       time.Time
	Status          Status
	UpdatedAt       time.Time
	CreatedAt       time.Time
}
