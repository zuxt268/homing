package repository

import (
	"context"
	"fmt"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"gorm.io/gorm"
)

type GoogleBusinessRepository interface {
	Get(ctx context.Context, f GoogleBusinessFilter) (*domain.GoogleBusinesses, error)
	FindAll(ctx context.Context, f GoogleBusinessFilter) ([]*domain.GoogleBusinesses, error)
	Count(ctx context.Context, f GoogleBusinessFilter) (int64, error)
	Exists(ctx context.Context, f GoogleBusinessFilter) (bool, error)
	Update(ctx context.Context, item *domain.GoogleBusinesses, f GoogleBusinessFilter) error
	Create(ctx context.Context, googleBusiness *domain.GoogleBusinesses) error
	Delete(ctx context.Context, f GoogleBusinessFilter) error
}

type googleBusinessRepository struct {
	db *gorm.DB
}

func NewGoogleBusinessRepository(db *gorm.DB) GoogleBusinessRepository {
	return &googleBusinessRepository{
		db: db,
	}
}

func (r *googleBusinessRepository) Get(ctx context.Context, f GoogleBusinessFilter) (*domain.GoogleBusinesses, error) {
	var gb model.GoogleBusiness
	err := f.Mod(r.getDB(ctx)).Find(&gb).Error
	if err != nil {
		return nil, err
	}
	return &domain.GoogleBusinesses{
		ID:        gb.ID,
		Name:      gb.Name,
		Title:     gb.Title,
		CreatedAt: gb.CreatedAt,
	}, nil
}

func (r *googleBusinessRepository) FindAll(ctx context.Context, f GoogleBusinessFilter) ([]*domain.GoogleBusinesses, error) {
	var gbList []*model.GoogleBusiness
	err := f.Mod(r.getDB(ctx)).Find(&gbList).Error
	if err != nil {
		return nil, err
	}
	googleBusinessList := make([]*domain.GoogleBusinesses, 0, len(gbList))
	for _, gb := range gbList {
		googleBusinessList = append(googleBusinessList, &domain.GoogleBusinesses{
			ID:        gb.ID,
			Name:      gb.Name,
			Title:     gb.Title,
			CreatedAt: gb.CreatedAt,
		})
	}
	return googleBusinessList, nil
}

func (r *googleBusinessRepository) Count(ctx context.Context, f GoogleBusinessFilter) (int64, error) {
	var total int64
	f.Offset = nil
	f.Limit = nil
	err := f.Mod(r.getDB(ctx)).Model(model.GoogleBusiness{}).Count(&total).Error
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return total, nil
}

func (r *googleBusinessRepository) Exists(ctx context.Context, f GoogleBusinessFilter) (bool, error) {
	var gbList []*model.GoogleBusiness
	err := f.Mod(r.getDB(ctx)).Find(&gbList).Error
	if err != nil {
		return false, err
	}
	return len(gbList) > 0, nil
}

func (r *googleBusinessRepository) Update(ctx context.Context, googleBusiness *domain.GoogleBusinesses, f GoogleBusinessFilter) error {
	m := &model.GoogleBusiness{
		ID:    googleBusiness.ID,
		Name:  googleBusiness.Name,
		Title: googleBusiness.Title,
	}
	return r.getDB(ctx).Omit("created_at").Save(m).Error
}

func (r *googleBusinessRepository) Create(ctx context.Context, googleBusiness *domain.GoogleBusinesses) error {
	m := model.GoogleBusiness{
		Name:  googleBusiness.Name,
		Title: googleBusiness.Title,
	}
	if err := r.getDB(ctx).Create(&m).Error; err != nil {
		return err
	}
	googleBusiness.ID = m.ID
	return nil
}

func (r *googleBusinessRepository) Delete(ctx context.Context, f GoogleBusinessFilter) error {
	return f.Mod(r.getDB(ctx)).Delete(model.GoogleBusiness{}).Error
}

func (r *googleBusinessRepository) getDB(ctx context.Context) *gorm.DB {
	if v, ok := ctx.Value(TxKey{}).(*gorm.DB); ok {
		return v.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

type GoogleBusinessFilter struct {
	ID            *int
	Name          *string
	Title         *string
	Limit         *int
	Offset        *int
	All           *bool
	PartialName   *string
	PartialTitle  *string
	OrderByIDDesc *bool
}

func (p *GoogleBusinessFilter) Mod(db *gorm.DB) *gorm.DB {
	if p.All != nil && *p.All {
		return db.Where("1 = 1")
	}
	if p.ID != nil {
		db = db.Where("id = ?", *p.ID)
	}
	if p.Name != nil {
		db = db.Where("name = ?", *p.Name)
	}
	if p.Title != nil {
		db = db.Where("title = ?", *p.Title)
	}

	if p.PartialName != nil || p.PartialTitle != nil {
		var orConditions []string
		var orValues []interface{}
		if p.PartialName != nil {
			orConditions = append(orConditions, "name like ?")
			orValues = append(orValues, "%"+*p.PartialName+"%")
		}
		if p.PartialTitle != nil {
			orConditions = append(orConditions, "title like ?")
			orValues = append(orValues, "%"+*p.PartialTitle+"%")
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
