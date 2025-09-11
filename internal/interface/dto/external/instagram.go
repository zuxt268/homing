package external

import "github.com/zuxt268/homing/internal/domain/entity"

type InstagramRequest struct {
	AccessToken string `param:"access_token"`
	Fields      string `param:"fields"`
}

type InstagramGetAccountResponse struct {
	Accounts struct {
		Data []struct {
			Name                     string `json:"name"`
			InstagramBusinessAccount struct {
				Name     string `json:"name"`
				Username string `json:"username"`
				Id       string `json:"id"`
			} `json:"instagram_business_account"`
			Id string `json:"id"`
		} `json:"data"`
	} `json:"accounts"`
	Id string `json:"id"`
}

type InstagramGetPostsResponse struct {
	Media struct {
		Data []struct {
			Id        string `json:"id"`
			Permalink string `json:"permalink"`
			Timestamp string `json:"timestamp"`
			MediaType string `json:"media_type"`
			MediaUrl  string `json:"media_url"`
			Children  struct {
				Data []struct {
					MediaType string `json:"media_type"`
					MediaUrl  string `json:"media_url"`
					Id        string `json:"id"`
				} `json:"data"`
			} `json:"children,omitempty"`
			Caption string `json:"caption,omitempty"`
		} `json:"data"`
		Paging struct {
			Cursors struct {
				Before string `json:"before"`
				After  string `json:"after"`
			} `json:"cursors"`
		} `json:"paging"`
	} `json:"media"`
	Id string `json:"id"`
}

func ToInstagramAccountEntity(dto *InstagramGetAccountResponse) *entity.InstagramAccount {
	return &entity.InstagramAccount{}
}

func ToInstagramPostsEntity(dto *InstagramGetPostsResponse) []entity.InstagramPost {
	var posts []entity.InstagramPost

	return posts
}
