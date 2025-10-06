package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/external"
	"github.com/zuxt268/homing/internal/interface/dto/req"
)

type MockSlack struct {
	mock.Mock
}

func (m *MockSlack) SendTokenExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSlack) Alert(ctx context.Context, message string, wi domain.WordpressInstagram) error {
	args := m.Called(ctx, message, wi)
	return args.Error(0)
}

func (m *MockSlack) Success(ctx context.Context, wi *domain.WordpressInstagram, wordpressURL, instagramURL string) error {
	args := m.Called(ctx, wi, wordpressURL, instagramURL)
	return args.Error(0)
}

func (m *MockSlack) SendMessage(ctx context.Context, payload external.SlackRequest) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func TestTokenUsecase_GetToken(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(*MockTokenRepository, *MockInstagramAdapter, *MockSlack)
		wantErr bool
		check   func(t *testing.T, slack *MockSlack)
	}{
		{
			name: "トークンが有効期限まで10日以上ある場合",
			setup: func(tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, slack *MockSlack) {
				tokenRepo.On("First", ctx).Return("valid_token", nil)
				resp := &external.DebugTokenResponse{}
				resp.Data.ExpiresAt = time.Now().AddDate(0, 0, 20).Unix() // 20日後
				igAdapter.On("DebugToken", ctx, "valid_token").Return(resp, nil)
			},
			wantErr: false,
			check: func(t *testing.T, slack *MockSlack) {
				// Slackには通知されない
				slack.AssertNotCalled(t, "SendTokenExpired", mock.Anything)
			},
		},
		{
			name: "トークンが有効期限まで10日未満の場合",
			setup: func(tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, slack *MockSlack) {
				tokenRepo.On("First", ctx).Return("expiring_token", nil)
				resp := &external.DebugTokenResponse{}
				resp.Data.ExpiresAt = time.Now().AddDate(0, 0, 5).Unix() // 5日後
				igAdapter.On("DebugToken", ctx, "expiring_token").Return(resp, nil)
				slack.On("SendTokenExpired", ctx).Return(nil)
			},
			wantErr: false,
			check: func(t *testing.T, slack *MockSlack) {
				// Slackに通知される
				slack.AssertCalled(t, "SendTokenExpired", ctx)
			},
		},
		{
			name: "トークンが既に期限切れの場合",
			setup: func(tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, slack *MockSlack) {
				tokenRepo.On("First", ctx).Return("expired_token", nil)
				resp := &external.DebugTokenResponse{}
				resp.Data.ExpiresAt = time.Now().AddDate(0, 0, -1).Unix() // 1日前
				igAdapter.On("DebugToken", ctx, "expired_token").Return(resp, nil)
				slack.On("SendTokenExpired", ctx).Return(nil)
			},
			wantErr: false,
			check: func(t *testing.T, slack *MockSlack) {
				// Slackに通知される
				slack.AssertCalled(t, "SendTokenExpired", ctx)
			},
		},
		{
			name: "トークン取得失敗",
			setup: func(tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, slack *MockSlack) {
				tokenRepo.On("First", ctx).Return("", errors.New("token not found"))
			},
			wantErr: true,
			check:   nil,
		},
		{
			name: "DebugToken API呼び出し失敗",
			setup: func(tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, slack *MockSlack) {
				tokenRepo.On("First", ctx).Return("invalid_token", nil)
				igAdapter.On("DebugToken", ctx, "invalid_token").Return(nil, errors.New("API error"))
			},
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenRepo := new(MockTokenRepository)
			igAdapter := new(MockInstagramAdapter)
			slack := new(MockSlack)

			tt.setup(tokenRepo, igAdapter, slack)

			usecase := NewTokenUsecase(igAdapter, slack, tokenRepo)
			result, err := usecase.GetToken(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.Token)
			}

			if tt.check != nil {
				tt.check(t, slack)
			}

			tokenRepo.AssertExpectations(t)
			igAdapter.AssertExpectations(t)
		})
	}
}

func TestTokenUsecase_UpdateToken(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		request req.UpdateToken
		setup   func(*MockTokenRepository)
		wantErr bool
	}{
		{
			name: "トークンを正常に更新",
			request: req.UpdateToken{
				Token: "new_token_value",
			},
			setup: func(tokenRepo *MockTokenRepository) {
				tokenRepo.On("DeleteInsert", ctx, "new_token_value").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "トークン更新失敗",
			request: req.UpdateToken{
				Token: "invalid_token",
			},
			setup: func(tokenRepo *MockTokenRepository) {
				tokenRepo.On("DeleteInsert", ctx, "invalid_token").Return(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "空のトークンで更新",
			request: req.UpdateToken{
				Token: "",
			},
			setup: func(tokenRepo *MockTokenRepository) {
				tokenRepo.On("DeleteInsert", ctx, "").Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenRepo := new(MockTokenRepository)
			igAdapter := new(MockInstagramAdapter)
			slack := new(MockSlack)

			tt.setup(tokenRepo)

			usecase := NewTokenUsecase(igAdapter, slack, tokenRepo)
			err := usecase.UpdateToken(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tokenRepo.AssertExpectations(t)
		})
	}
}