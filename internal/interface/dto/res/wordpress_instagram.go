package res

import "time"

type WordpressInstagramList struct {
	WordpressInstagramList []WordpressInstagram `json:"instagram_list"`
}

type WordpressInstagram struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	WordpressDomain    string    `json:"wordpress_domain"`
	WordpressSiteTitle string    `json:"wordpress_site_title"`
	InstagramID        string    `json:"instagram_id"`
	InstagramName      string    `json:"instagram_name"`
	Memo               string    `json:"memo"`
	StartDate          time.Time `json:"start_date"`
	Status             int       `json:"status"`
	DeleteHash         bool      `json:"delete_hash"`
	CustomerType       int       `json:"customer_type"`
	Posts              []Post    `json:"posts"`
}

type Post struct {
	WordpressUrl string    `json:"wordpress_url"`
	InstagramUrl string    `json:"instagram_url"`
	CreatedAt    time.Time `json:"created_at"`
}
