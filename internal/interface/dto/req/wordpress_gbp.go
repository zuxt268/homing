package req

import "time"

type GetWordpressGbp struct {
	Limit  *int    `query:"limit"`
	Offset *int    `query:"offset"`
	Name   *string `query:"name"`
	Status *int    `query:"status"`
}

type GetWordpressGbpDetail struct {
	Limit  *int `query:"limit"`
	Offset *int `query:"offset"`
}

type WordpressGbp struct {
	Name            string    `json:"name"`
	WordpressDomain string    `json:"wordpress_domain"`
	BusinessName    string    `json:"business_name"`
	Memo            string    `json:"memo"`
	StartDate       time.Time `json:"start_date"`
	Status          int       `json:"status"`
}
