package res

import "github.com/zuxt268/homing/internal/domain"

type Customer struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	WordpressUrl  string      `json:"wordpress_url"`
	AccessToken   string      `json:"access_token"`
	InstagramList []Instagram `json:"instagram_list"`
}

type Instagram struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetCustomer(entity *domain.Customer) *Customer {
	instagramList := make([]Instagram, 0)
	for i := range entity.InstagramBusinessAccountID {
		instagramList = append(instagramList, Instagram{
			ID:   entity.InstagramBusinessAccountID[i],
			Name: entity.InstagramBusinessAccountName[i],
		})
	}
	return &Customer{
		ID:            entity.ID,
		Name:          entity.Name,
		WordpressUrl:  entity.WordpressUrl,
		AccessToken:   entity.FacebookToken,
		InstagramList: instagramList,
	}
}
