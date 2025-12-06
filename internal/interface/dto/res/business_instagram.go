package res

import "time"

type BusinessInstagram struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	BusinessName string    `json:"business_name"`
	InstagramID  string    `json:"instagram_id"`
	Memo         string    `json:"memo"`
	StartDate    time.Time `json:"start_date"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
