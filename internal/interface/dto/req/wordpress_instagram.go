package req

import "time"

type GetWordpressInstagram struct {
	Limit              *int    `query:"limit"`
	Offset             *int    `query:"offset"`
	Name               *string `query:"name"`
	WordpressDomain    *string `query:"wordpress_domain"`
	WordpressSiteTitle *string `query:"wordpress_site_title"`
	InstagramID        *string `query:"instagram_id"`
	InstagramName      *string `query:"instagram_name"`
	Status             *int    `query:"status"`
	DeleteHash         *bool   `query:"delete_hash"`
	CustomerType       *int    `query:"customer_type"`
}

type CreateWordpressInstagram struct {
	Name            string    `json:"name"`
	WordpressDomain string    `json:"wordpress_domain"`
	InstagramID     string    `json:"instagram_id"`
	Memo            string    `json:"memo"`
	StartDate       time.Time `json:"start_date"`
	Status          int       `json:"status"`
	DeleteHash      bool      `json:"delete_hash"`
	CustomerType    int       `json:"customer_type"`
}

type UpdateWordpressInstagram struct {
	ID           *int       `json:"id"`
	Name         *string    `json:"name"`
	Wordpress    *string    `json:"wordpress_domain"`
	InstagramID  *string    `json:"instagram_id"`
	Memo         *string    `json:"memo"`
	StartDate    *time.Time `json:"start_date"`
	Status       *int       `json:"status"`
	DeleteHash   *bool      `json:"delete_hash"`
	CustomerType *int       `json:"customer_type"`
}
