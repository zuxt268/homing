package repository

import (
	"context"
	"fmt"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"gorm.io/gorm"
)

type BusinessInstagramRepository interface {
	Get(ctx context.Context, f BusinessInstagramFilter) (*domain.BusinessInstagram, error)
	FindAll(ctx context.Context, f BusinessInstagramFilter) ([]*domain.BusinessInstagram, error)
	Count(ctx context.Context, f BusinessInstagramFilter) (int64, error)
	Exists(ctx context.Context, f BusinessInstagramFilter) (bool, error)
	Update(ctx context.Context, item *domain.BusinessInstagram, f BusinessInstagramFilter) error
	Create(ctx context.Context, businessInstagram *domain.BusinessInstagram) error
	Delete(ctx context.Context, f BusinessInstagramFilter) error
}

type businessInstagramRepository struct {
	db *gorm.DB
}

func NewBusinessInstagramRepository(db *gorm.DB) BusinessInstagramRepository {
	return &businessInstagramRepository{
		db: db,
	}
}

func (r *businessInstagramRepository) Get(ctx context.Context, f BusinessInstagramFilter) (*domain.BusinessInstagram, error) {
	var bi model.BusinessInstagram
	err := f.Mod(r.getDB(ctx)).Find(&bi).Error
	if err != nil {
		return nil, err
	}
	return &domain.BusinessInstagram{
		ID:            bi.ID,
		Name:          bi.Name,
		Memo:          bi.Memo,
		InstagramID:   bi.InstagramID,
		InstagramName: bi.InstagramName,
		BusinessName:  bi.BusinessName,
		BusinessTitle: bi.BusinessTitle,
		StartDate:     bi.StartDate,
		Status:        domain.Status(bi.Status),
		UpdatedAt:     bi.UpdatedAt,
		CreatedAt:     bi.CreatedAt,
	}, nil
}

func (r *businessInstagramRepository) FindAll(ctx context.Context, f BusinessInstagramFilter) ([]*domain.BusinessInstagram, error) {
	var biList []*model.BusinessInstagram
	err := f.Mod(r.getDB(ctx)).Find(&biList).Error
	if err != nil {
		return nil, err
	}
	businessInstagramList := make([]*domain.BusinessInstagram, 0, len(biList))
	for _, bi := range biList {
		businessInstagramList = append(businessInstagramList, &domain.BusinessInstagram{
			ID:            bi.ID,
			Name:          bi.Name,
			Memo:          bi.Memo,
			InstagramID:   bi.InstagramID,
			InstagramName: bi.InstagramName,
			BusinessName:  bi.BusinessName,
			BusinessTitle: bi.BusinessTitle,
			StartDate:     bi.StartDate,
			Status:        domain.Status(bi.Status),
			UpdatedAt:     bi.UpdatedAt,
			CreatedAt:     bi.CreatedAt,
		})
	}
	return businessInstagramList, nil
}

func (r *businessInstagramRepository) Count(ctx context.Context, f BusinessInstagramFilter) (int64, error) {
	var total int64
	f.Offset = nil
	f.Limit = nil
	err := f.Mod(r.getDB(ctx)).Model(model.BusinessInstagram{}).Count(&total).Error
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return total, nil
}

func (r *businessInstagramRepository) Exists(ctx context.Context, f BusinessInstagramFilter) (bool, error) {
	var biList []*model.BusinessInstagram
	err := f.Mod(r.getDB(ctx)).Find(&biList).Error
	if err != nil {
		return false, err
	}
	return len(biList) > 0, nil
}

func (r *businessInstagramRepository) Update(ctx context.Context, businessInstagram *domain.BusinessInstagram, f BusinessInstagramFilter) error {
	m := &model.BusinessInstagram{
		ID:            businessInstagram.ID,
		Name:          businessInstagram.Name,
		Memo:          businessInstagram.Memo,
		InstagramID:   businessInstagram.InstagramID,
		InstagramName: businessInstagram.InstagramName,
		BusinessName:  businessInstagram.BusinessName,
		BusinessTitle: businessInstagram.BusinessTitle,
		StartDate:     businessInstagram.StartDate,
		Status:        int(businessInstagram.Status),
	}
	return r.getDB(ctx).Omit("created_at").Save(m).Error
}

func (r *businessInstagramRepository) Create(ctx context.Context, businessInstagram *domain.BusinessInstagram) error {
	m := model.BusinessInstagram{
		Name:          businessInstagram.Name,
		Memo:          businessInstagram.Memo,
		InstagramID:   businessInstagram.InstagramID,
		InstagramName: businessInstagram.InstagramName,
		BusinessName:  businessInstagram.BusinessName,
		BusinessTitle: businessInstagram.BusinessTitle,
		StartDate:     businessInstagram.StartDate,
		Status:        int(businessInstagram.Status),
	}
	if err := r.getDB(ctx).Create(&m).Error; err != nil {
		return err
	}
	businessInstagram.ID = m.ID
	businessInstagram.CreatedAt = m.CreatedAt
	businessInstagram.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *businessInstagramRepository) Delete(ctx context.Context, f BusinessInstagramFilter) error {
	return f.Mod(r.getDB(ctx)).Delete(model.BusinessInstagram{}).Error
}

func (r *businessInstagramRepository) getDB(ctx context.Context) *gorm.DB {
	if v, ok := ctx.Value(TxKey{}).(*gorm.DB); ok {
		return v.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

type BusinessInstagramFilter struct {
	ID            *int
	Name          *string
	InstagramID   *string
	InstagramName *string
	BusinessName  *string
	Status        *int
	Limit         *int
	Offset        *int
	All           *bool

	PartialName          *string
	PartialInstagramName *string
	PartialBusinessName  *string
	OrderByIDDesc        *bool
}

func (p *BusinessInstagramFilter) Mod(db *gorm.DB) *gorm.DB {
	if p.All != nil && *p.All {
		return db.Where("1")
	}
	if p.ID != nil {
		db = db.Where("id = ?", *p.ID)
	}
	if p.Name != nil {
		db = db.Where("name = ?", *p.Name)
	}
	if p.InstagramID != nil {
		db = db.Where("instagram_id = ?", *p.InstagramID)
	}
	if p.InstagramName != nil {
		db = db.Where("instagram_name = ?", *p.InstagramName)
	}
	if p.BusinessName != nil {
		db = db.Where("business_name = ?", *p.BusinessName)
	}
	if p.Status != nil {
		db = db.Where("status = ?", *p.Status)
	}

	if p.PartialName != nil || p.PartialInstagramName != nil || p.PartialBusinessName != nil {
		var orConditions []string
		var orValues []interface{}
		if p.PartialName != nil {
			orConditions = append(orConditions, "name like ?")
			orValues = append(orValues, "%"+*p.PartialName+"%")
		}
		if p.PartialInstagramName != nil {
			orConditions = append(orConditions, "instagram_name like ?")
			orValues = append(orValues, "%"+*p.PartialInstagramName+"%")
		}
		if p.PartialBusinessName != nil {
			orConditions = append(orConditions, "business_name like ?")
			orValues = append(orValues, "%"+*p.PartialBusinessName+"%")
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
