package repository

import (
	"context"
	"time"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"gorm.io/gorm"
)

type WordpressInstagramRepository interface {
	Get(ctx context.Context, f WordpressInstagramFilter) (*domain.WordpressInstagram, error)
	FindAll(ctx context.Context, f WordpressInstagramFilter) ([]*domain.WordpressInstagram, error)
	Exists(ctx context.Context, f WordpressInstagramFilter) (bool, error)
	Update(ctx context.Context, item *domain.WordpressInstagram, f WordpressInstagramFilter) error
	Create(ctx context.Context, wordpressInstagram *domain.WordpressInstagram) error
	Delete(ctx context.Context, f WordpressInstagramFilter) error
}

type wordpressInstagramRepository struct {
	db *gorm.DB
}

func NewWordpressInstagramRepository(db *gorm.DB) WordpressInstagramRepository {
	return &wordpressInstagramRepository{
		db: db,
	}
}

func (r *wordpressInstagramRepository) Get(ctx context.Context, f WordpressInstagramFilter) (*domain.WordpressInstagram, error) {
	var wi model.WordpressInstagram
	err := f.Mod(r.getDB(ctx)).Find(&wi).Error
	if err != nil {
		return nil, err
	}
	return &domain.WordpressInstagram{
		ID:                 wi.ID,
		Name:               wi.Name,
		WordpressDomain:    wi.WordpressDomain,
		WordpressSiteTitle: wi.WordpressSiteTitle,
		InstagramID:        wi.InstagramID,
		InstagramName:      wi.InstagramName,
		Memo:               wi.Memo,
		StartDate:          wi.StartDate,
		Status:             domain.Status(wi.Status),
		DeleteHash:         wi.DeleteHash,
		CustomerType:       domain.CustomerType(wi.CustomerType),
	}, nil
}

func (r *wordpressInstagramRepository) FindAll(ctx context.Context, f WordpressInstagramFilter) ([]*domain.WordpressInstagram, error) {
	var wiList []*model.WordpressInstagram
	err := f.Mod(r.getDB(ctx)).Find(&wiList).Error
	if err != nil {
		return nil, err
	}
	wordpressInstagramList := make([]*domain.WordpressInstagram, 0, len(wiList))
	for _, wi := range wiList {
		wordpressInstagramList = append(wordpressInstagramList, &domain.WordpressInstagram{
			ID:                 wi.ID,
			Name:               wi.Name,
			WordpressDomain:    wi.WordpressDomain,
			WordpressSiteTitle: wi.WordpressSiteTitle,
			InstagramID:        wi.InstagramID,
			InstagramName:      wi.InstagramName,
			Memo:               wi.Memo,
			StartDate:          wi.StartDate,
			Status:             domain.Status(wi.Status),
			DeleteHash:         wi.DeleteHash,
			CustomerType:       domain.CustomerType(wi.CustomerType),
		})
	}
	return wordpressInstagramList, nil
}

func (r *wordpressInstagramRepository) Exists(ctx context.Context, f WordpressInstagramFilter) (bool, error) {
	var wiList []*model.WordpressInstagram
	err := f.Mod(r.getDB(ctx)).Find(&wiList).Error
	if err != nil {
		return false, err
	}
	return len(wiList) > 0, nil
}

func (r *wordpressInstagramRepository) Update(ctx context.Context, wordpressInstagram *domain.WordpressInstagram, f WordpressInstagramFilter) error {
	return f.Mod(r.getDB(ctx)).Updates(&model.WordpressInstagram{
		ID:                 wordpressInstagram.ID,
		Name:               wordpressInstagram.Name,
		WordpressDomain:    wordpressInstagram.WordpressDomain,
		WordpressSiteTitle: wordpressInstagram.WordpressSiteTitle,
		InstagramID:        wordpressInstagram.InstagramID,
		InstagramName:      wordpressInstagram.InstagramName,
		Memo:               wordpressInstagram.Memo,
		StartDate:          wordpressInstagram.StartDate,
		Status:             int(wordpressInstagram.Status),
		DeleteHash:         wordpressInstagram.DeleteHash,
		CustomerType:       int(wordpressInstagram.CustomerType),
	}).Error
}

func (r *wordpressInstagramRepository) Create(ctx context.Context, wordpressInstagram *domain.WordpressInstagram) error {
	return r.getDB(ctx).Create(&model.WordpressInstagram{
		Name:               wordpressInstagram.Name,
		WordpressDomain:    wordpressInstagram.WordpressDomain,
		WordpressSiteTitle: wordpressInstagram.WordpressSiteTitle,
		InstagramID:        wordpressInstagram.InstagramID,
		InstagramName:      wordpressInstagram.InstagramName,
		Memo:               wordpressInstagram.Memo,
		StartDate:          wordpressInstagram.StartDate,
		Status:             int(wordpressInstagram.Status),
		DeleteHash:         wordpressInstagram.DeleteHash,
		CustomerType:       int(wordpressInstagram.CustomerType),
	}).Error
}

func (r *wordpressInstagramRepository) Delete(ctx context.Context, f WordpressInstagramFilter) error {
	return f.Mod(r.getDB(ctx)).Delete(model.WordpressInstagram{}).Error
}

func (r *wordpressInstagramRepository) getDB(ctx context.Context) *gorm.DB {
	if v, ok := ctx.Value(TxKey{}).(*gorm.DB); ok {
		return v.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

type WordpressInstagramFilter struct {
	ID                 *int
	Name               *string
	WordpressDomain    *string
	WordpressSiteTitle *string
	InstagramID        *string
	InstagramName      *string
	Memo               *string
	StartDate          *time.Time
	Status             *int
	DeleteHash         *bool
	CustomerType       *int
	Limit              *int
	Offset             *int
}

func (p *WordpressInstagramFilter) Mod(db *gorm.DB) *gorm.DB {
	if p.ID != nil {
		db = db.Where("id = ?", *p.ID)
	}
	if p.Name != nil {
		db = db.Where("name = ?", *p.Name)
	}
	if p.WordpressDomain != nil {
		db = db.Where("wordpress_domain = ?", *p.WordpressDomain)
	}
	if p.WordpressSiteTitle != nil {
		db = db.Where("wordpress_site_title = ?", *p.WordpressSiteTitle)
	}
	if p.InstagramID != nil {
		db = db.Where("instagram_id = ?", *p.InstagramID)
	}
	if p.InstagramName != nil {
		db = db.Where("instagram_name = ?", *p.InstagramName)
	}
	if p.Memo != nil {
		db = db.Where("memo = ?", *p.Memo)
	}
	if p.StartDate != nil {
		db = db.Where("start_date >= ?", *p.StartDate)
	}
	if p.Status != nil {
		db = db.Where("status = ?", *p.Status)
	}
	if p.DeleteHash != nil {
		db = db.Where("delete_hash = ?", *p.DeleteHash)
	}
	if p.CustomerType != nil {
		db = db.Where("customer_type = ?", *p.CustomerType)
	}
	if p.Limit != nil {
		db = db.Limit(*p.Limit)
		if p.Offset != nil {
			db = db.Offset(*p.Offset)
		}
	}
	return db
}
