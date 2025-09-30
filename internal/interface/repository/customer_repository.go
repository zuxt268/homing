package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	SaveCustomer(ctx context.Context, customer *domain.Customer) error
	GetCustomer(ctx context.Context, id int) (*domain.Customer, error)
	FindAllCustomers(ctx context.Context, filter CustomerFilter) ([]*domain.Customer, error)
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (a *customerRepository) SaveCustomer(ctx context.Context, customer *domain.Customer) error {
	return a.db.WithContext(ctx).Save(&model.Customer{
		ID:                           customer.ID,
		Name:                         customer.Name,
		Email:                        customer.Email,
		Password:                     customer.Password,
		WordpressURL:                 customer.WordpressUrl,
		FacebookToken:                customer.FacebookToken,
		StartDate:                    customer.StartDate,
		InstagramBusinessAccountID:   strings.Join(customer.InstagramBusinessAccountID, ","),
		InstagramBusinessAccountName: strings.Join(customer.InstagramBusinessAccountName, ","),
		InstagramTokenStatus:         customer.InstagramTokenStatus,
		DeleteHash:                   customer.DeleteHash,
		PaymentType:                  customer.PaymentType,
		Type:                         customer.Type,
	}).Error
}

func (a *customerRepository) GetCustomer(ctx context.Context, id int) (*domain.Customer, error) {
	var customer model.Customer
	err := a.db.WithContext(ctx).
		Where("id = ?", id).
		First(&customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 見つからなかった場合は nil を返す
		}
		return nil, err
	}
	return &domain.Customer{
		ID:                           customer.ID,
		Name:                         customer.Name,
		WordpressUrl:                 customer.WordpressURL,
		FacebookToken:                customer.FacebookToken,
		Email:                        customer.Email,
		Password:                     customer.Password,
		StartDate:                    customer.StartDate,
		InstagramBusinessAccountID:   strings.Split(customer.InstagramBusinessAccountID, ","),
		InstagramBusinessAccountName: strings.Split(customer.InstagramBusinessAccountName, ","),
		InstagramTokenStatus:         customer.InstagramTokenStatus,
		DeleteHash:                   customer.DeleteHash,
		PaymentType:                  customer.PaymentType,
		Type:                         customer.Type,
	}, nil
}

type CustomerFilter struct {
	ID                           *int
	Name                         *string
	Email                        *string
	Password                     *string
	WordpressURL                 *string
	FacebookToken                *string
	StartDate                    *time.Time
	InstagramBusinessAccountID   *string
	InstagramBusinessAccountName *string
	InstagramTokenStatus         *int
	DeleteHash                   *bool
	PaymentType                  *string
	Type                         *int
}

func (f *CustomerFilter) Mod(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.Name != nil {
		db = db.Where("name = ?", *f.Name)
	}
	if f.Email != nil {
		db = db.Where("email = ?", *f.Email)
	}
	if f.Password != nil {
		db = db.Where("password = ?", *f.Password)
	}
	if f.WordpressURL != nil {
		db = db.Where("wordpress_url = ?", *f.WordpressURL)
	}
	if f.FacebookToken != nil {
		db = db.Where("facebook_token = ?", *f.FacebookToken)
	}
	if f.InstagramBusinessAccountID != nil {
		db = db.Where("instagram_business_account_id = ?", *f.InstagramBusinessAccountID)
	}
	if f.InstagramBusinessAccountName != nil {
		db = db.Where("instagram_business_account_name = ?", *f.InstagramBusinessAccountName)
	}
	if f.InstagramTokenStatus != nil {
		db = db.Where("instagram_token_status = ?", *f.InstagramTokenStatus)
	}
	if f.DeleteHash != nil {
		db = db.Where("delete_hash = ?", *f.DeleteHash)
	}
	if f.PaymentType != nil {
		db = db.Where("payment_type = ?", *f.PaymentType)
	}
	return db
}

func (a *customerRepository) FindAllCustomers(ctx context.Context, filter CustomerFilter) ([]*domain.Customer, error) {
	var customers []*model.Customer
	db := filter.Mod(a.db).WithContext(ctx).Find(&customers)
	if db.Error != nil {
		return nil, db.Error
	}
	result := make([]*domain.Customer, 0, len(customers))
	for _, customer := range customers {
		result = append(result, &domain.Customer{
			ID:                         customer.ID,
			Name:                       customer.Name,
			WordpressUrl:               customer.WordpressURL,
			FacebookToken:              customer.FacebookToken,
			InstagramBusinessAccountID: strings.Split(customer.InstagramBusinessAccountID, ","),
		})
	}
	return result, nil
}
