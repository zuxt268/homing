package res

import "time"

type GoogleBusinessList struct {
	GoogleBusinessList []GoogleBusiness `json:"google_business_list"`
	Paginate
}

type GoogleBusiness struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}