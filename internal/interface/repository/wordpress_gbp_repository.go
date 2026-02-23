package repository

import (
	"context"
	"fmt"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"gorm.io/gorm"
)

type WordpressGbpRepository interface {
	Get(ctx context.Context, f WordpressGbpFilter) (*domain.WordpressGbp, error)
	FindAll(ctx context.Context, f WordpressGbpFilter) ([]*domain.WordpressGbp, error)
	Count(ctx context.Context, f WordpressGbpFilter) (int64, error)
	Exists(ctx context.Context, f WordpressGbpFilter) (bool, error)
	Update(ctx context.Context, item *domain.WordpressGbp, f WordpressGbpFilter) error
	Create(ctx context.Context, wordpressGbp *domain.WordpressGbp) error
	Delete(ctx context.Context, f WordpressGbpFilter) error
}

type wordpressGbpRepository struct {
	db *gorm.DB
}

func NewWordpressGbpRepository(db *gorm.DB) WordpressGbpRepository {
	return &wordpressGbpRepository{
		db: db,
	}
}

func (r *wordpressGbpRepository) Get(ctx context.Context, f WordpressGbpFilter) (*domain.WordpressGbp, error) {
	var wg model.WordpressGbp
	err := f.Mod(r.getDB(ctx)).Find(&wg).Error
	if err != nil {
		return nil, err
	}
	return &domain.WordpressGbp{
		ID:              wg.ID,
		Name:            wg.Name,
		Memo:            wg.Memo,
		WordpressDomain: wg.WordpressDomain,
		BusinessName:    wg.BusinessName,
		BusinessTitle:   wg.BusinessTitle,
		MapsURL:         wg.MapsURL,
		StartDate:       wg.StartDate,
		Status:          domain.Status(wg.Status),
		UpdatedAt:       wg.UpdatedAt,
		CreatedAt:       wg.CreatedAt,
	}, nil
}

func (r *wordpressGbpRepository) FindAll(ctx context.Context, f WordpressGbpFilter) ([]*domain.WordpressGbp, error) {
	var wgList []*model.WordpressGbp
	err := f.Mod(r.getDB(ctx)).Find(&wgList).Error
	if err != nil {
		return nil, err
	}
	wordpressGbpList := make([]*domain.WordpressGbp, 0, len(wgList))
	for _, wg := range wgList {
		wordpressGbpList = append(wordpressGbpList, &domain.WordpressGbp{
			ID:              wg.ID,
			Name:            wg.Name,
			Memo:            wg.Memo,
			WordpressDomain: wg.WordpressDomain,
			BusinessName:    wg.BusinessName,
			BusinessTitle:   wg.BusinessTitle,
			MapsURL:         wg.MapsURL,
			Status:          domain.Status(wg.Status),
			UpdatedAt:       wg.UpdatedAt,
			CreatedAt:       wg.CreatedAt,
		})
	}
	return wordpressGbpList, nil
}

func (r *wordpressGbpRepository) Count(ctx context.Context, f WordpressGbpFilter) (int64, error) {
	var total int64
	f.Offset = nil
	f.Limit = nil
	err := f.Mod(r.getDB(ctx)).Model(model.WordpressGbp{}).Count(&total).Error
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return total, nil
}

func (r *wordpressGbpRepository) Exists(ctx context.Context, f WordpressGbpFilter) (bool, error) {
	var wgList []*model.WordpressGbp
	err := f.Mod(r.getDB(ctx)).Find(&wgList).Error
	if err != nil {
		return false, err
	}
	return len(wgList) > 0, nil
}

func (r *wordpressGbpRepository) Update(ctx context.Context, wordpressGbp *domain.WordpressGbp, f WordpressGbpFilter) error {
	m := &model.WordpressGbp{
		ID:              wordpressGbp.ID,
		Name:            wordpressGbp.Name,
		Memo:            wordpressGbp.Memo,
		WordpressDomain: wordpressGbp.WordpressDomain,
		BusinessName:    wordpressGbp.BusinessName,
		BusinessTitle:   wordpressGbp.BusinessTitle,
		MapsURL:         wordpressGbp.MapsURL,
		StartDate:       wordpressGbp.StartDate,
		Status:          int(wordpressGbp.Status),
	}
	return r.getDB(ctx).Omit("created_at").Save(m).Error
}

func (r *wordpressGbpRepository) Create(ctx context.Context, wordpressGbp *domain.WordpressGbp) error {
	m := model.WordpressGbp{
		Name:            wordpressGbp.Name,
		Memo:            wordpressGbp.Memo,
		WordpressDomain: wordpressGbp.WordpressDomain,
		BusinessName:    wordpressGbp.BusinessName,
		BusinessTitle:   wordpressGbp.BusinessTitle,
		MapsURL:         wordpressGbp.MapsURL,
		StartDate:       wordpressGbp.StartDate,
		Status:          int(wordpressGbp.Status),
	}
	if err := r.getDB(ctx).Create(&m).Error; err != nil {
		return err
	}
	wordpressGbp.ID = m.ID
	wordpressGbp.CreatedAt = m.CreatedAt
	wordpressGbp.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *wordpressGbpRepository) Delete(ctx context.Context, f WordpressGbpFilter) error {
	return f.Mod(r.getDB(ctx)).Delete(model.WordpressGbp{}).Error
}

func (r *wordpressGbpRepository) getDB(ctx context.Context) *gorm.DB {
	if v, ok := ctx.Value(TxKey{}).(*gorm.DB); ok {
		return v.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

type WordpressGbpFilter struct {
	ID              *int
	Name            *string
	WordpressDomain *string
	BusinessName    *string
	Status          *int
	Limit           *int
	Offset          *int
	All             *bool

	PartialName   *string
	OrderByIDDesc *bool
}

func (p *WordpressGbpFilter) Mod(db *gorm.DB) *gorm.DB {
	if p.All != nil && *p.All {
		return db.Where("1")
	}
	if p.ID != nil {
		db = db.Where("id = ?", *p.ID)
	}
	if p.Name != nil {
		db = db.Where("name = ?", *p.Name)
	}
	if p.WordpressDomain != nil {
		db = db.Where("wordpress_domain = ?", *p.WordpressDomain)
	}
	if p.BusinessName != nil {
		db = db.Where("business_name = ?", *p.BusinessName)
	}
	if p.Status != nil {
		db = db.Where("status = ?", *p.Status)
	}
	if p.PartialName != nil {
		db = db.Where("name like ?", "%"+*p.PartialName+"%")
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
