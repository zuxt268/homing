package req

import "time"

type GetWordpressInstagram struct {
	Limit        *int
	Offset       *int
	Name         *string
	Wordpress    *string
	InstagramID  *string
	Status       *int
	DeleteHash   *bool
	CustomerType *int
}

type CreateWordpressInstagram struct {
	Name         string    `json:"name"`
	Wordpress    string    `json:"wordpress"`
	InstagramID  string    `json:"instagram_id"`
	Memo         string    `json:"memo"`
	StartDate    time.Time `json:"start_date"`
	Status       int       `json:"status"`
	DeleteHash   bool      `json:"delete_hash"`
	CustomerType int       `json:"customer_type"`
}

type UpdateWordpressInstagram struct {
	ID           *int       `json:"id"`
	Name         *string    `json:"name"`
	Wordpress    *string    `json:"wordpress"`
	InstagramID  *string    `json:"instagram_id"`
	Memo         *string    `json:"memo"`
	StartDate    *time.Time `json:"start_date"`
	Status       *int       `json:"status"`
	DeleteHash   *bool      `json:"delete_hash"`
	CustomerType *int       `json:"customer_type"`
}
