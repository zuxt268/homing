package res

import "time"

type WordpressGbp struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	WordpressDomain string    `json:"wordpress_domain"`
	BusinessName    string    `json:"business_name"`
	BusinessTitle   string    `json:"business_title"`
	Memo            string    `json:"memo"`
	MapsURL         string    `json:"maps_url"`
	StartDate       time.Time `json:"start_date"`
	Status          int       `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type WordpressGbpList struct {
	WordpressGbpList []WordpressGbp `json:"wordpress_gbp_list"`
	Paginate
}

type WordpressGbpDetail struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	WordpressDomain   string    `json:"wordpress_domain"`
	BusinessName      string    `json:"business_name"`
	BusinessTitle     string    `json:"business_title"`
	Memo              string    `json:"memo"`
	MapsURL           string    `json:"maps_url"`
	StartDate         time.Time `json:"start_date"`
	Status            int       `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	GooglePhotosCount int64     `json:"google_photos_count"`
	GooglePostsCount  int64     `json:"google_posts_count"`
}
