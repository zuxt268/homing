package usecase

import (
	"context"
	"time"

	"github.com/zuxt268/homing/internal/interface/adapter"
	"github.com/zuxt268/homing/internal/interface/dto/req"
	"github.com/zuxt268/homing/internal/interface/dto/res"
	"github.com/zuxt268/homing/internal/interface/repository"
)

type TokenUsecase interface {
	GetToken(ctx context.Context) (*res.Token, error)
	UpdateToken(ctx context.Context, body req.UpdateToken) error
}

type tokenUsecase struct {
	instagramAdapter adapter.InstagramAdapter
	slack            adapter.Slack
	tokenRepo        repository.TokenRepository
}

func NewTokenUsecase(
	instagramAdapter adapter.InstagramAdapter,
	slack adapter.Slack,
	tokenRepo repository.TokenRepository,
) TokenUsecase {
	return &tokenUsecase{
		instagramAdapter: instagramAdapter,
		slack:            slack,
		tokenRepo:        tokenRepo,
	}
}

func (u *tokenUsecase) GetToken(ctx context.Context) (*res.Token, error) {
	token, err := u.tokenRepo.First(ctx)
	if err != nil {
		return nil, err
	}
	debug, err := u.instagramAdapter.DebugToken(ctx, token)
	if err != nil {
		return nil, err
	}
	expiredAt := time.Unix(debug.Data.ExpiresAt, 0)
	tenDaysLater := time.Now().AddDate(0, 0, 10)

	if tenDaysLater.After(expiredAt) {
		_ = u.slack.SendTokenExpired(ctx)
	}

	return &res.Token{
		Token:    token,
		ExpireAt: expiredAt,
	}, nil
}

func (u *tokenUsecase) UpdateToken(ctx context.Context, req req.UpdateToken) error {
	return u.tokenRepo.DeleteInsert(ctx, req.Token)
}