package res

import "github.com/zuxt268/homing/internal/domain"

type Customer struct {
	ID                         int      `json:"id"`
	Name                       string   `json:"name"`
	WordpressUrl               string   `json:"wordpress_url"`
	AccessToken                string   `json:"access_token"`
	InstagramBusinessAccountID []string `json:"instagram_account_id"`
}

func GetCustomer(entity *domain.Customer) *Customer {
	return &Customer{
		ID:                         entity.ID,
		Name:                       entity.Name,
		WordpressUrl:               entity.WordpressUrl,
		AccessToken:                entity.FacebookToken,
		InstagramBusinessAccountID: entity.InstagramBusinessAccountID,
	}
}
