package res

import "time"

type WordpressInstagramList struct {
	WordpressInstagramList []WordpressInstagram `json:"instagram_list"`
}

type WordpressInstagram struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Wordpress    string    `json:"wordpress"`
	InstagramID  string    `json:"instagram_id"`
	Memo         string    `json:"memo"`
	StartDate    time.Time `json:"start_date"`
	Status       int       `json:"status"`
	DeleteHash   bool      `json:"delete_hash"`
	CustomerType int       `json:"customer_type"`
}
