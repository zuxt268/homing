package repository

import (
	"context"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"gorm.io/gorm"
)

type PostRepository interface {
	ExistPost(ctx context.Context, filter PostFilter) (bool, error)
	CreatePost(ctx context.Context, post *model.Post) error
	GetPosts(ctx context.Context, filter PostFilter) ([]domain.Post, error)
	CountPosts(ctx context.Context, filter PostFilter) (int64, error)
}

type postRepository struct {
	db *gorm.DB
}

type PostFilter struct {
	ID                   *int
	MediaID              *string
	CustomerID           *int
	Timestamp            *string
	MediaURL             *string
	Permalink            *string
	WordpressLink        *string
	OrderByCreatedAtDesc *bool
	Limit                *int
	Offset               *int
}

func (p *PostFilter) Mod(db *gorm.DB) *gorm.DB {
	if p.ID != nil {
		db = db.Where("id = ?", *p.ID)
	}
	if p.MediaID != nil {
		db = db.Where("media_id = ?", *p.MediaID)
	}
	if p.CustomerID != nil {
		db = db.Where("customer_id = ?", *p.CustomerID)
	}
	if p.Timestamp != nil {
		db = db.Where("timestamp = ?", *p.Timestamp)
	}
	if p.MediaURL != nil {
		db = db.Where("media_url = ?", *p.MediaURL)
	}
	if p.Permalink != nil {
		db = db.Where("permalink = ?", *p.Permalink)
	}
	if p.WordpressLink != nil {
		db = db.Where("wordpress_link = ?", *p.WordpressLink)
	}
	if p.OrderByCreatedAtDesc != nil {
		db = db.Order("created_at desc")
	}
	if p.Limit != nil {
		db = db.Limit(*p.Limit)
		if p.Offset != nil {
			db = db.Offset(*p.Offset)
		}
	}
	return db
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) ExistPost(ctx context.Context, filter PostFilter) (bool, error) {
	var posts []*model.Post

	err := filter.Mod(r.db).WithContext(ctx).Find(&posts).Error
	if err != nil {
		return false, err
	}
	return len(posts) > 0, nil
}

func (r *postRepository) CreatePost(ctx context.Context, post *model.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepository) GetPosts(ctx context.Context, filter PostFilter) ([]domain.Post, error) {
	var posts []*model.Post
	err := filter.Mod(r.db).WithContext(ctx).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	result := make([]domain.Post, len(posts))
	for i, post := range posts {
		result[i] = domain.Post{
			ID:           post.ID,
			WordpressURL: post.WordpressLink,
			InstagramURL: post.Permalink,
			CreatedAt:    post.CreatedAt,
		}
	}
	return result, nil
}

func (r *postRepository) CountPosts(ctx context.Context, filter PostFilter) (int64, error) {
	var total int64
	filter.Offset = nil
	err := filter.Mod(r.db).WithContext(ctx).Model(model.Post{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
