package repository

import (
	"context"

	"github.com/zuxt268/homing/internal/interface/dto/model"
	"gorm.io/gorm"
)

type TokenRepository interface {
	First(ctx context.Context) (string, error)
	DeleteInsert(ctx context.Context, token string) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (r *tokenRepository) First(ctx context.Context) (string, error) {
	token := model.Token{}
	err := r.db.WithContext(ctx).First(&token).Error
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (r *tokenRepository) DeleteInsert(ctx context.Context, token string) error {
	err := r.db.WithContext(ctx).Delete(&model.Token{}, "1").Error
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(&model.Token{
		Token: token,
	}).Error
}
