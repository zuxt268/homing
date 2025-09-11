package external

import "github.com/zuxt268/homing/internal/domain/entity"

type WordpressPostPayload struct {
	Email         string `json:"email"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	FeaturedMedia int    `json:"featured_media"`
}

type WordpressHeader struct {
}

type WordpressPostResponse struct {
	PostId  int    `json:"post_id"`
	PostUrl string `json:"post_url"`
	Message string `json:"message"`
}

type WordpressPostInput struct {
	Domain  string
	MediaID int
	Post    entity.InstagramPost
}
