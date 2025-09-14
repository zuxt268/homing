package external

import (
	"github.com/zuxt268/homing/internal/domain/entity"
)

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

func ToInstagramAccountEntity(dto *InstagramGetAccountResponse) []entity.InstagramAccount {
	accounts := make([]entity.InstagramAccount, 0, len(dto.Accounts.Data))
	for _, account := range dto.Accounts.Data {
		if account.InstagramBusinessAccount.Name != "" {
			accounts = append(accounts, entity.InstagramAccount{
				InstagramAccountName:     account.InstagramBusinessAccount.Name,
				InstagramAccountID:       account.InstagramBusinessAccount.Id,
				InstagramAccountUsername: account.InstagramBusinessAccount.Username,
			})
		}
	}
	return accounts
}

func ToInstagramPostsEntity(dto *InstagramGetPostsResponse) []entity.InstagramPost {
	var posts []entity.InstagramPost
	for _, post := range dto.Media.Data {
		children := make([]entity.InstagramPostChildren, 0, len(post.Children.Data))
		for _, child := range post.Children.Data {
			children = append(children, entity.InstagramPostChildren{
				MediaType: child.MediaType,
				MediaURL:  child.MediaUrl,
				ID:        child.Id,
			})
		}
		posts = append(posts, entity.InstagramPost{
			ID:        post.Id,
			Permalink: post.Permalink,
			Caption:   post.Caption,
			Timestamp: post.Timestamp,
			MediaType: post.MediaType,
			MediaURL:  post.MediaUrl,
			Children:  children,
		})
	}
	return posts
}
