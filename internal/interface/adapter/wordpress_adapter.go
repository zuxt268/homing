package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/zuxt268/homing/internal/domain/entity"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/dto/external"
)

type WordpressAdapter interface {
	Post(ctx context.Context, in external.WordpressPostInput) (*entity.Post, error)
	//FileUpload(ctx context.Context, file multipart.File, header *multipart.FileHeader) error
}

func NewWordpressAdapter(httpDriver driver.HttpDriver, adminEmail string) WordpressAdapter {
	return &wordpressAdapter{
		httpDriver: httpDriver,
		adminEmail: adminEmail,
	}
}

type wordpressAdapter struct {
	httpDriver driver.HttpDriver
	adminEmail string
}

func (a *wordpressAdapter) Post(ctx context.Context, in external.WordpressPostInput) (*entity.Post, error) {
	reqBody := external.WordpressPostPayload{
		Email:         a.adminEmail,
		Title:         in.Post.GetTitle(),
		Content:       in.Post.GetContent(),
		FeaturedMedia: in.MediaID,
	}
	header := external.WordpressHeader{} // TODO
	u, err := url.Parse(in.Domain)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("rest_route", "/rodut/v1/create_post")
	u.RawQuery = q.Encode()

	resp, err := a.httpDriver.Post(ctx, u.String(), &reqBody, &header)
	if err != nil {
		return nil, fmt.Errorf("記事の投稿に失敗: %w", err)
	}
	var postDto external.WordpressPostResponse
	if err := json.Unmarshal(resp, &postDto); err != nil {
		return nil, fmt.Errorf("JSONの変換に失敗: %w", err)
	}

	return &entity.Post{
		ID: postDto.PostId,
	}, nil
}
