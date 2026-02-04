package req

import (
	"time"
)

type GetBusinessInstagram struct {
	Limit       *int    `query:"limit"`
	Offset      *int    `query:"offset"`
	Name        *string `query:"name"`
	InstagramID *string `query:"instagram_id"`
	Status      *int    `query:"status"`
}

type GetBusinessInstagramDetail struct {
	Limit  *int `query:"limit"`
	Offset *int `query:"offset"`
}

type BusinessInstagram struct {
	Name         string    `json:"name"`
	BusinessName string    `json:"business_name"`
	InstagramID  string    `json:"instagram_id"`
	Memo         string    `json:"memo"`
	StartDate    time.Time `json:"start_date"`
	Status       int       `json:"status"`
}
