package adapter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zuxt268/homing/internal/config"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/dto/external"
)

type InstagramAdapter interface {
	GetPosts25(ctx context.Context, token, instagramID string) ([]domain.InstagramPost, error)
	GetPostsAll(ctx context.Context, token, instagramID string) ([]domain.InstagramPost, error)
	GetAccount(ctx context.Context, token, instagramID string) (*domain.InstagramAccount, error)
	DebugToken(ctx context.Context, userToken string) (*external.DebugTokenResponse, error)
}

func NewInstagramAdapter(httpDriver driver.HttpDriver) InstagramAdapter {
	return &instagramAdapter{
		httpDriver:   httpDriver,
		clientID:     config.Env.ClientID,
		clientSecret: config.Env.ClientSecret,
	}
}

const (
	baseURL = "https://graph.facebook.com/v23.0"
)

type instagramAdapter struct {
	httpDriver   driver.HttpDriver
	clientID     string
	clientSecret string
}

func (a *instagramAdapter) GetPosts25(ctx context.Context, token string, instagramID string) ([]domain.InstagramPost, error) {
	req := &external.InstagramRequest{
		AccessToken: token,
		Fields:      "media{id,permalink,caption,timestamp,media_type,media_url,children{media_type,media_url}}",
	}
	endpoint := baseURL + "/" + instagramID
	resp, err := a.httpDriver.Get(ctx, endpoint, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	var postsDto external.InstagramGetPostsResponse
	if err := json.Unmarshal(resp, &postsDto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal instagram posts response: %w", err)
	}
	return external.ToInstagramPostsEntity(&postsDto), nil
}

func (a *instagramAdapter) GetPostsAll(ctx context.Context, token, instagramID string) ([]domain.InstagramPost, error) {
	result := make([]domain.InstagramPost, 0)

	req := &external.InstagramRequest{
		AccessToken: token,
		Fields:      "media{id,permalink,caption,timestamp,media_type,media_url,children{media_type,media_url}}",
	}
	endpoint := baseURL + "/" + instagramID
	resp, err := a.httpDriver.Get(ctx, endpoint, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	var postsDto external.InstagramGetPostsResponse
	if err := json.Unmarshal(resp, &postsDto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal posts response: %w", err)
	}

	entityList := external.ToInstagramPostsEntity(&postsDto)
	result = append(result, entityList...)
	nextURL := postsDto.Media.Paging.Next

	for {
		if nextURL == "" {
			break
		}
		respBody, err := a.httpDriver.Get(ctx, nextURL, nil, nil)
		if err != nil {
			return nil, err
		}
		var postsDto external.InstagramGetPostsNextResponse
		if err := json.Unmarshal(respBody, &postsDto); err != nil {
			return nil, fmt.Errorf("failed to unmarshal posts response: %w", err)
		}
		posts := external.NextResponseToInstagramPostsEntity(&postsDto)
		result = append(result, posts...)
		nextURL = postsDto.Paging.Next
	}

	return result, nil
}

func (a *instagramAdapter) DebugToken(ctx context.Context, token string) (*external.DebugTokenResponse, error) {
	appToken := fmt.Sprintf("%s|%s", a.clientID, a.clientSecret)
	endpoint := "https://graph.facebook.com/debug_token"
	req := external.DebugTokenRequest{
		AccessToken: appToken,
		InputToken:  token,
	}

	respBody, err := a.httpDriver.Get(ctx, endpoint, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to debug token: %w", err)
	}

	var dto external.DebugTokenResponse
	if err := json.Unmarshal(respBody, &dto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal debug token response: %w", err)
	}

	return &dto, nil
}

func (a *instagramAdapter) GetAccount(ctx context.Context, token, instagramID string) (*domain.InstagramAccount, error) {
	req := external.InstagramRequest{
		AccessToken: token,
		Fields:      "name,username",
		Limit:       100,
	}
	endpoint := baseURL + "/" + instagramID
	respBody, err := a.httpDriver.Get(ctx, endpoint, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get instagram account: %w", err)
	}
	var accountDto external.InstagramAccount
	if err := json.Unmarshal(respBody, &accountDto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal instagram account response: %w", err)
	}
	return &domain.InstagramAccount{
		InstagramAccountUserName: accountDto.Username,
		InstagramAccountName:     accountDto.Name,
		InstagramAccountID:       accountDto.ID,
	}, nil
}
