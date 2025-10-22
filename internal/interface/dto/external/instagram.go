package external

import (
	"github.com/zuxt268/homing/internal/domain"
)

type InstagramRequest struct {
	AccessToken string `param:"access_token"`
	Fields      string `param:"fields"`
	Limit       int    `param:"limit"`
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
			Next     string `json:"next"`
			Previous string `json:"previous"`
		} `json:"paging"`
	} `json:"media"`
	Id string `json:"id"`
}

type InstagramGetPostsNextResponse struct {
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
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"paging"`
}

func ToInstagramPostsEntity(dto *InstagramGetPostsResponse) []domain.InstagramPost {
	var posts []domain.InstagramPost
	for _, post := range dto.Media.Data {
		children := make([]domain.InstagramPostChildren, 0, len(post.Children.Data))
		for _, child := range post.Children.Data {
			children = append(children, domain.InstagramPostChildren{
				MediaType: child.MediaType,
				MediaURL:  child.MediaUrl,
				ID:        child.Id,
			})
		}
		posts = append(posts, domain.InstagramPost{
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

func NextResponseToInstagramPostsEntity(dto *InstagramGetPostsNextResponse) []domain.InstagramPost {
	var posts []domain.InstagramPost
	for _, post := range dto.Data {
		children := make([]domain.InstagramPostChildren, 0, len(post.Children.Data))
		for _, child := range post.Children.Data {
			children = append(children, domain.InstagramPostChildren{
				MediaType: child.MediaType,
				MediaURL:  child.MediaUrl,
				ID:        child.Id,
			})
		}
		posts = append(posts, domain.InstagramPost{
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
