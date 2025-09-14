package adapter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zuxt268/homing/internal/domain/entity"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/dto/external"
)

type InstagramAdapter interface {
	GetAccount(ctx context.Context, accessToken string) ([]entity.InstagramAccount, error)
	GetPosts(ctx context.Context, accessToken string, accountID string) ([]entity.InstagramPost, error)
}

func NewInstagramAdapter(httpDriver driver.HttpDriver) InstagramAdapter {
	return &instagramAdapter{
		httpDriver: httpDriver,
	}
}

const (
	baseURL = "https://graph.facebook.com/v23.0"
)

type instagramAdapter struct {
	httpDriver driver.HttpDriver
}

func (a *instagramAdapter) GetAccount(ctx context.Context, accessToken string) ([]entity.InstagramAccount, error) {
	req := external.InstagramRequest{
		AccessToken: accessToken,
		Fields:      "accounts{name,instagram_business_account{name,username}}",
	}
	endpoint := baseURL + "/me"
	respBody, err := a.httpDriver.Get(ctx, endpoint, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get instagram account: %w", err)
	}
	var accountDto external.InstagramGetAccountResponse
	if err := json.Unmarshal(respBody, &accountDto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal instagram account response: %w", err)
	}
	return external.ToInstagramAccountEntity(&accountDto), nil
}

func (a *instagramAdapter) GetPosts(ctx context.Context, accessToken string, accountID string) ([]entity.InstagramPost, error) {
	req := &external.InstagramRequest{
		AccessToken: accessToken,
		Fields:      "media{id,permalink,caption,timestamp,media_type,media_url,children{media_type,media_url}}",
	}
	endpoint := baseURL + "/" + accountID
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
