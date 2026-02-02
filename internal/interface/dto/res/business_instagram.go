package res

import "time"

type BusinessInstagram struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	BusinessName  string    `json:"business_name"`
	InstagramID   string    `json:"instagram_id"`
	InstagramName string    `json:"instagram_name"`
	Memo          string    `json:"memo"`
	StartDate     time.Time `json:"start_date"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BusinessInstagramList struct {
	BusinessInstagramList []BusinessInstagram `json:"business_instagram_list"`
	Paginate
}
