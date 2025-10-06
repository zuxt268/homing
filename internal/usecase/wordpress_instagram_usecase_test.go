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
	"github.com/zuxt268/homing/internal/interface/repository"
)

// Mock repositories and adapters
type MockWordpressInstagramRepository struct {
	mock.Mock
}

func (m *MockWordpressInstagramRepository) Get(ctx context.Context, f repository.WordpressInstagramFilter) (*domain.WordpressInstagram, error) {
	args := m.Called(ctx, f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WordpressInstagram), args.Error(1)
}

func (m *MockWordpressInstagramRepository) FindAll(ctx context.Context, f repository.WordpressInstagramFilter) ([]*domain.WordpressInstagram, error) {
	args := m.Called(ctx, f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.WordpressInstagram), args.Error(1)
}

func (m *MockWordpressInstagramRepository) Exists(ctx context.Context, f repository.WordpressInstagramFilter) (bool, error) {
	args := m.Called(ctx, f)
	return args.Bool(0), args.Error(1)
}

func (m *MockWordpressInstagramRepository) Update(ctx context.Context, item *domain.WordpressInstagram, f repository.WordpressInstagramFilter) error {
	args := m.Called(ctx, item, f)
	return args.Error(0)
}

func (m *MockWordpressInstagramRepository) Create(ctx context.Context, wordpressInstagram *domain.WordpressInstagram) error {
	args := m.Called(ctx, wordpressInstagram)
	return args.Error(0)
}

func (m *MockWordpressInstagramRepository) Delete(ctx context.Context, f repository.WordpressInstagramFilter) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) First(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockTokenRepository) DeleteInsert(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

type MockInstagramAdapter struct {
	mock.Mock
}

func (m *MockInstagramAdapter) GetPosts(ctx context.Context, token, instagramID string) ([]domain.InstagramPost, error) {
	args := m.Called(ctx, token, instagramID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InstagramPost), args.Error(1)
}

func (m *MockInstagramAdapter) GetAccount(ctx context.Context, token, instagramID string) (*domain.InstagramAccount, error) {
	args := m.Called(ctx, token, instagramID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InstagramAccount), args.Error(1)
}

func (m *MockInstagramAdapter) DebugToken(ctx context.Context, userToken string) (*external.DebugTokenResponse, error) {
	args := m.Called(ctx, userToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.DebugTokenResponse), args.Error(1)
}

type MockWordpressAdapter struct {
	mock.Mock
}

func (m *MockWordpressAdapter) GetTitle(ctx context.Context, domain string) (string, error) {
	args := m.Called(ctx, domain)
	return args.String(0), args.Error(1)
}

func (m *MockWordpressAdapter) Post(ctx context.Context, in external.WordpressPostInput) (*domain.Post, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockWordpressAdapter) FileUpload(ctx context.Context, in external.WordpressFileUploadInput) (*external.WordpressFileUploadResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.WordpressFileUploadResponse), args.Error(1)
}

func TestWordpressInstagramUsecase_CreateWordpressInstagram(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		request req.CreateWordpressInstagram
		setup   func(*MockWordpressInstagramRepository, *MockTokenRepository, *MockInstagramAdapter, *MockWordpressAdapter)
		wantErr bool
	}{
		{
			name: "正常にWordpressInstagramを作成",
			request: req.CreateWordpressInstagram{
				Name:            "Test Site",
				WordpressDomain: "https://test.example.com",
				InstagramID:     "123456789",
				Memo:            "Test memo",
				StartDate:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:          1,
				DeleteHash:      false,
				CustomerType:    1,
			},
			setup: func(wiRepo *MockWordpressInstagramRepository, tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, wpAdapter *MockWordpressAdapter) {
				tokenRepo.On("First", ctx).Return("test_token", nil)
				igAdapter.On("GetAccount", ctx, "test_token", "123456789").Return(&domain.InstagramAccount{
					InstagramAccountUserName: "testuser",
					InstagramAccountName:     "Test User",
					InstagramAccountID:       "123456789",
				}, nil)
				wpAdapter.On("GetTitle", ctx, "https://test.example.com").Return("Test Site Title", nil)
				wiRepo.On("Create", ctx, mock.AnythingOfType("*domain.WordpressInstagram")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "トークン取得失敗",
			request: req.CreateWordpressInstagram{
				Name:            "Test Site",
				WordpressDomain: "https://test.example.com",
				InstagramID:     "123456789",
				Memo:            "Test memo",
				StartDate:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:          1,
				DeleteHash:      false,
				CustomerType:    1,
			},
			setup: func(wiRepo *MockWordpressInstagramRepository, tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, wpAdapter *MockWordpressAdapter) {
				tokenRepo.On("First", ctx).Return("", errors.New("token not found"))
			},
			wantErr: true,
		},
		{
			name: "Instagramアカウント取得失敗",
			request: req.CreateWordpressInstagram{
				Name:            "Test Site",
				WordpressDomain: "https://test.example.com",
				InstagramID:     "invalid_id",
				Memo:            "Test memo",
				StartDate:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:          1,
				DeleteHash:      false,
				CustomerType:    1,
			},
			setup: func(wiRepo *MockWordpressInstagramRepository, tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, wpAdapter *MockWordpressAdapter) {
				tokenRepo.On("First", ctx).Return("test_token", nil)
				igAdapter.On("GetAccount", ctx, "test_token", "invalid_id").Return(nil, errors.New("account not found"))
			},
			wantErr: true,
		},
		{
			name: "WordPressタイトル取得失敗",
			request: req.CreateWordpressInstagram{
				Name:            "Test Site",
				WordpressDomain: "https://invalid.example.com",
				InstagramID:     "123456789",
				Memo:            "Test memo",
				StartDate:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:          1,
				DeleteHash:      false,
				CustomerType:    1,
			},
			setup: func(wiRepo *MockWordpressInstagramRepository, tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, wpAdapter *MockWordpressAdapter) {
				tokenRepo.On("First", ctx).Return("test_token", nil)
				igAdapter.On("GetAccount", ctx, "test_token", "123456789").Return(&domain.InstagramAccount{
					InstagramAccountUserName: "testuser",
					InstagramAccountName:     "Test User",
					InstagramAccountID:       "123456789",
				}, nil)
				wpAdapter.On("GetTitle", ctx, "https://invalid.example.com").Return("", errors.New("wordpress site not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wiRepo := new(MockWordpressInstagramRepository)
			tokenRepo := new(MockTokenRepository)
			igAdapter := new(MockInstagramAdapter)
			wpAdapter := new(MockWordpressAdapter)

			tt.setup(wiRepo, tokenRepo, igAdapter, wpAdapter)

			usecase := NewWordpressInstagramUsecase(wiRepo, tokenRepo, igAdapter, wpAdapter)
			result, err := usecase.CreateWordpressInstagram(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
				assert.Equal(t, tt.request.WordpressDomain, result.WordpressDomain)
				assert.NotEmpty(t, result.WordpressSiteTitle)
				assert.NotEmpty(t, result.InstagramName)
			}

			wiRepo.AssertExpectations(t)
			tokenRepo.AssertExpectations(t)
			igAdapter.AssertExpectations(t)
			wpAdapter.AssertExpectations(t)
		})
	}
}

func TestWordpressInstagramUsecase_UpdateWordpressInstagram(t *testing.T) {
	ctx := context.Background()
	testID := 1

	tests := []struct {
		name    string
		request req.UpdateWordpressInstagram
		setup   func(*MockWordpressInstagramRepository, *MockTokenRepository, *MockInstagramAdapter, *MockWordpressAdapter)
		wantErr bool
	}{
		{
			name: "名前のみ更新",
			request: req.UpdateWordpressInstagram{
				ID:   &testID,
				Name: stringPtr("Updated Name"),
			},
			setup: func(wiRepo *MockWordpressInstagramRepository, tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, wpAdapter *MockWordpressAdapter) {
				wiRepo.On("Get", ctx, mock.AnythingOfType("repository.WordpressInstagramFilter")).Return(&domain.WordpressInstagram{
					ID:                 testID,
					Name:               "Original Name",
					WordpressDomain:    "https://test.example.com",
					WordpressSiteTitle: "Test Title",
					InstagramID:        "123456789",
					InstagramName:      "testuser",
					Status:             1,
					CustomerType:       1,
				}, nil)
				wiRepo.On("Update", ctx, mock.AnythingOfType("*domain.WordpressInstagram"), mock.AnythingOfType("repository.WordpressInstagramFilter")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "WordpressDomainを更新（WordpressSiteTitleも自動更新）",
			request: req.UpdateWordpressInstagram{
				ID:        &testID,
				Wordpress: stringPtr("https://updated.example.com"),
			},
			setup: func(wiRepo *MockWordpressInstagramRepository, tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, wpAdapter *MockWordpressAdapter) {
				wiRepo.On("Get", ctx, mock.AnythingOfType("repository.WordpressInstagramFilter")).Return(&domain.WordpressInstagram{
					ID:                 testID,
					Name:               "Test Name",
					WordpressDomain:    "https://test.example.com",
					WordpressSiteTitle: "Test Title",
					InstagramID:        "123456789",
					InstagramName:      "testuser",
					Status:             1,
					CustomerType:       1,
				}, nil)
				wpAdapter.On("GetTitle", ctx, "https://updated.example.com").Return("Updated Title", nil)
				wiRepo.On("Update", ctx, mock.AnythingOfType("*domain.WordpressInstagram"), mock.AnythingOfType("repository.WordpressInstagramFilter")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "InstagramIDを更新（InstagramNameも自動更新）",
			request: req.UpdateWordpressInstagram{
				ID:          &testID,
				InstagramID: stringPtr("987654321"),
			},
			setup: func(wiRepo *MockWordpressInstagramRepository, tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, wpAdapter *MockWordpressAdapter) {
				wiRepo.On("Get", ctx, mock.AnythingOfType("repository.WordpressInstagramFilter")).Return(&domain.WordpressInstagram{
					ID:                 testID,
					Name:               "Test Name",
					WordpressDomain:    "https://test.example.com",
					WordpressSiteTitle: "Test Title",
					InstagramID:        "123456789",
					InstagramName:      "testuser",
					Status:             1,
					CustomerType:       1,
				}, nil)
				tokenRepo.On("First", ctx).Return("test_token", nil)
				igAdapter.On("GetAccount", ctx, "test_token", "987654321").Return(&domain.InstagramAccount{
					InstagramAccountUserName: "updateduser",
					InstagramAccountName:     "Updated User",
					InstagramAccountID:       "987654321",
				}, nil)
				wiRepo.On("Update", ctx, mock.AnythingOfType("*domain.WordpressInstagram"), mock.AnythingOfType("repository.WordpressInstagramFilter")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "データ取得失敗",
			request: req.UpdateWordpressInstagram{
				ID:   &testID,
				Name: stringPtr("Updated Name"),
			},
			setup: func(wiRepo *MockWordpressInstagramRepository, tokenRepo *MockTokenRepository, igAdapter *MockInstagramAdapter, wpAdapter *MockWordpressAdapter) {
				wiRepo.On("Get", ctx, mock.AnythingOfType("repository.WordpressInstagramFilter")).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wiRepo := new(MockWordpressInstagramRepository)
			tokenRepo := new(MockTokenRepository)
			igAdapter := new(MockInstagramAdapter)
			wpAdapter := new(MockWordpressAdapter)

			tt.setup(wiRepo, tokenRepo, igAdapter, wpAdapter)

			usecase := NewWordpressInstagramUsecase(wiRepo, tokenRepo, igAdapter, wpAdapter)
			result, err := usecase.UpdateWordpressInstagram(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			wiRepo.AssertExpectations(t)
			tokenRepo.AssertExpectations(t)
			igAdapter.AssertExpectations(t)
			wpAdapter.AssertExpectations(t)
		})
	}
}

func TestWordpressInstagramUsecase_GetWordpressInstagram(t *testing.T) {
	ctx := context.Background()
	testID := 1

	wiRepo := new(MockWordpressInstagramRepository)
	tokenRepo := new(MockTokenRepository)
	igAdapter := new(MockInstagramAdapter)
	wpAdapter := new(MockWordpressAdapter)

	wiRepo.On("Get", ctx, mock.AnythingOfType("repository.WordpressInstagramFilter")).Return(&domain.WordpressInstagram{
		ID:                 testID,
		Name:               "Test Name",
		WordpressDomain:    "https://test.example.com",
		WordpressSiteTitle: "Test Title",
		InstagramID:        "123456789",
		InstagramName:      "testuser",
		Memo:               "Test memo",
		StartDate:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Status:             1,
		DeleteHash:         false,
		CustomerType:       1,
	}, nil)

	usecase := NewWordpressInstagramUsecase(wiRepo, tokenRepo, igAdapter, wpAdapter)
	result, err := usecase.GetWordpressInstagram(ctx, testID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testID, result.ID)
	assert.Equal(t, "Test Name", result.Name)
	assert.Equal(t, "https://test.example.com", result.WordpressDomain)
	assert.Equal(t, "Test Title", result.WordpressSiteTitle)
	assert.Equal(t, "123456789", result.InstagramID)
	assert.Equal(t, "testuser", result.InstagramName)

	wiRepo.AssertExpectations(t)
}

func TestWordpressInstagramUsecase_DeleteWordpressInstagram(t *testing.T) {
	ctx := context.Background()
	testID := 1

	wiRepo := new(MockWordpressInstagramRepository)
	tokenRepo := new(MockTokenRepository)
	igAdapter := new(MockInstagramAdapter)
	wpAdapter := new(MockWordpressAdapter)

	wiRepo.On("Delete", ctx, mock.AnythingOfType("repository.WordpressInstagramFilter")).Return(nil)

	usecase := NewWordpressInstagramUsecase(wiRepo, tokenRepo, igAdapter, wpAdapter)
	err := usecase.DeleteWordpressInstagram(ctx, testID)

	assert.NoError(t, err)
	wiRepo.AssertExpectations(t)
}

// ヘルパー関数
func stringPtr(s string) *string {
	return &s
}