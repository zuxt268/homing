package repository

import (
	"context"
	"fmt"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"gorm.io/gorm"
)

type GooglePostRepository interface {
	Get(ctx context.Context, f GooglePostFilter) (*domain.GooglePost, error)
	FindAll(ctx context.Context, f GooglePostFilter) ([]*domain.GooglePost, error)
	Count(ctx context.Context, f GooglePostFilter) (int64, error)
	Exists(ctx context.Context, f GooglePostFilter) (bool, error)
	Update(ctx context.Context, item *domain.GooglePost, f GooglePostFilter) error
	Create(ctx context.Context, googlePost *domain.GooglePost) error
	Delete(ctx context.Context, f GooglePostFilter) error
}

type googlePostRepository struct {
	db *gorm.DB
}

func NewGooglePostRepository(db *gorm.DB) GooglePostRepository {
	return &googlePostRepository{
		db: db,
	}
}

func (r *googlePostRepository) Get(ctx context.Context, f GooglePostFilter) (*domain.GooglePost, error) {
	var gp model.GooglePost
	err := f.Mod(r.getDB(ctx)).Find(&gp).Error
	if err != nil {
		return nil, err
	}
	return &domain.GooglePost{
		ID:                gp.ID,
		GoogleBusinessURL: gp.GoogleBusinessURL,
		InstagramURL:      gp.InstagramURL,
		MediaID:           gp.MediaID,
		CustomerID:        gp.CustomerID,
		Name:              gp.Name,
		MediaFormat:       gp.MediaFormat,
		GoogleURL:         gp.GoogleURL,
		CreateTime:        gp.CreateTime,
		CreatedAt:         gp.CreatedAt,
	}, nil
}

func (r *googlePostRepository) FindAll(ctx context.Context, f GooglePostFilter) ([]*domain.GooglePost, error) {
	var gpList []*model.GooglePost
	err := f.Mod(r.getDB(ctx)).Find(&gpList).Error
	if err != nil {
		return nil, err
	}
	googlePostList := make([]*domain.GooglePost, 0, len(gpList))
	for _, gp := range gpList {
		googlePostList = append(googlePostList, &domain.GooglePost{
			ID:                gp.ID,
			GoogleBusinessURL: gp.GoogleBusinessURL,
			InstagramURL:      gp.InstagramURL,
			MediaID:           gp.MediaID,
			CustomerID:        gp.CustomerID,
			Name:              gp.Name,
			MediaFormat:       gp.MediaFormat,
			GoogleURL:         gp.GoogleURL,
			CreateTime:        gp.CreateTime,
			CreatedAt:         gp.CreatedAt,
		})
	}
	return googlePostList, nil
}

func (r *googlePostRepository) Count(ctx context.Context, f GooglePostFilter) (int64, error) {
	var total int64
	f.Offset = nil
	f.Limit = nil
	err := f.Mod(r.getDB(ctx)).Model(model.GooglePost{}).Count(&total).Error
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return total, nil
}

func (r *googlePostRepository) Exists(ctx context.Context, f GooglePostFilter) (bool, error) {
	var gpList []*model.GooglePost
	err := f.Mod(r.getDB(ctx)).Find(&gpList).Error
	if err != nil {
		return false, err
	}
	return len(gpList) > 0, nil
}

func (r *googlePostRepository) Update(ctx context.Context, googlePost *domain.GooglePost, f GooglePostFilter) error {
	m := &model.GooglePost{
		ID:                googlePost.ID,
		GoogleBusinessURL: googlePost.GoogleBusinessURL,
		InstagramURL:      googlePost.InstagramURL,
		MediaID:           googlePost.MediaID,
		CustomerID:        googlePost.CustomerID,
		Name:              googlePost.Name,
		MediaFormat:       googlePost.MediaFormat,
		GoogleURL:         googlePost.GoogleURL,
		CreateTime:        googlePost.CreateTime,
	}
	return r.getDB(ctx).Omit("created_at").Save(m).Error
}

func (r *googlePostRepository) Create(ctx context.Context, googlePost *domain.GooglePost) error {
	m := model.GooglePost{
		GoogleBusinessURL: googlePost.GoogleBusinessURL,
		InstagramURL:      googlePost.InstagramURL,
		MediaID:           googlePost.MediaID,
		CustomerID:        googlePost.CustomerID,
		Name:              googlePost.Name,
		MediaFormat:       googlePost.MediaFormat,
		GoogleURL:         googlePost.GoogleURL,
		CreateTime:        googlePost.CreateTime,
	}
	if err := r.getDB(ctx).Create(&m).Error; err != nil {
		return err
	}
	googlePost.ID = m.ID
	return nil
}

func (r *googlePostRepository) Delete(ctx context.Context, f GooglePostFilter) error {
	return f.Mod(r.getDB(ctx)).Delete(model.GooglePost{}).Error
}

func (r *googlePostRepository) getDB(ctx context.Context) *gorm.DB {
	if v, ok := ctx.Value(TxKey{}).(*gorm.DB); ok {
		return v.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

type GooglePostFilter struct {
	ID                *int
	GoogleBusinessURL *string
	InstagramURL      *string
	MediaID           *string
	CustomerID        *int
	Name              *string
	MediaFormat       *string
	GoogleURL         *string
	CreateTime        *string
	Limit             *int
	Offset            *int
	All               *bool

	PartialGoogleBusinessURL *string
	PartialInstagramURL      *string
	OrderByIDDesc            *bool
}

func (p *GooglePostFilter) Mod(db *gorm.DB) *gorm.DB {
	if p.All != nil && *p.All {
		return db.Where("1")
	}
	if p.ID != nil {
		db = db.Where("id = ?", *p.ID)
	}
	if p.GoogleBusinessURL != nil {
		db = db.Where("google_business_url = ?", *p.GoogleBusinessURL)
	}
	if p.InstagramURL != nil {
		db = db.Where("instagram_url = ?", *p.InstagramURL)
	}
	if p.MediaID != nil {
		db = db.Where("media_id = ?", *p.MediaID)
	}
	if p.CustomerID != nil {
		db = db.Where("customer_id = ?", *p.CustomerID)
	}
	if p.Name != nil {
		db = db.Where("name = ?", *p.Name)
	}
	if p.MediaFormat != nil {
		db = db.Where("media_format = ?", *p.MediaFormat)
	}
	if p.GoogleURL != nil {
		db = db.Where("google_url = ?", *p.GoogleURL)
	}
	if p.CreateTime != nil {
		db = db.Where("create_time = ?", *p.CreateTime)
	}

	if p.PartialGoogleBusinessURL != nil || p.PartialInstagramURL != nil {
		var orConditions []string
		var orValues []interface{}
		if p.PartialGoogleBusinessURL != nil {
			orConditions = append(orConditions, "google_business_url like ?")
			orValues = append(orValues, "%"+*p.PartialGoogleBusinessURL+"%")
		}
		if p.PartialInstagramURL != nil {
			orConditions = append(orConditions, "instagram_url like ?")
			orValues = append(orValues, "%"+*p.PartialInstagramURL+"%")
		}
		if len(orConditions) > 0 {
			query := orConditions[0]
			for i := 1; i < len(orConditions); i++ {
				query += " OR " + orConditions[i]
			}
			db = db.Where(query, orValues...)
		}
	}
	if p.OrderByIDDesc != nil {
		db = db.Order("id desc")
	}
	if p.Limit != nil {
		db = db.Limit(*p.Limit)
		if p.Offset != nil {
			db = db.Offset(*p.Offset)
		}
	}
	return db
}