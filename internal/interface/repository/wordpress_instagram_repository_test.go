package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/model"
)

func TestWordpressInstagramRepository_Create(t *testing.T) {
	repo := NewWordpressInstagramRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		wi      *domain.WordpressInstagram
		wantErr bool
	}{
		{
			name: "正常にWordpressInstagramを作成",
			wi: &domain.WordpressInstagram{
				Name:               "Test Site",
				WordpressDomain:    "https://test.example.com",
				WordpressSiteTitle: "Test Site Title",
				InstagramID:        "123456789",
				InstagramName:      "testuser",
				Memo:               "Test memo",
				StartDate:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:             1,
				DeleteHash:         false,
				CustomerType:       1,
			},
			wantErr: false,
		},
		{
			name: "空のMemoでWordpressInstagramを作成",
			wi: &domain.WordpressInstagram{
				Name:               "Test Site 2",
				WordpressDomain:    "https://test2.example.com",
				WordpressSiteTitle: "Test Site 2 Title",
				InstagramID:        "987654321",
				InstagramName:      "testuser2",
				Memo:               "",
				StartDate:          time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
				Status:             0,
				DeleteHash:         true,
				CustomerType:       2,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.wi)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotZero(t, tt.wi.ID)

			// 作成されたデータを確認
			got, err := repo.Get(ctx, WordpressInstagramFilter{ID: &tt.wi.ID})
			assert.NoError(t, err)
			assert.Equal(t, tt.wi.Name, got.Name)
			assert.Equal(t, tt.wi.WordpressDomain, got.WordpressDomain)
			assert.Equal(t, tt.wi.WordpressSiteTitle, got.WordpressSiteTitle)
			assert.Equal(t, tt.wi.InstagramID, got.InstagramID)
			assert.Equal(t, tt.wi.InstagramName, got.InstagramName)
			assert.Equal(t, tt.wi.Memo, got.Memo)
			assert.Equal(t, tt.wi.Status, got.Status)
			assert.Equal(t, tt.wi.DeleteHash, got.DeleteHash)
			assert.Equal(t, tt.wi.CustomerType, got.CustomerType)
		})
	}
}

func TestWordpressInstagramRepository_Get(t *testing.T) {
	repo := NewWordpressInstagramRepository(db)
	ctx := context.Background()

	// テストデータを作成
	testWI := &domain.WordpressInstagram{
		Name:               "Get Test Site",
		WordpressDomain:    "https://gettest.example.com",
		WordpressSiteTitle: "Get Test Title",
		InstagramID:        "111222333",
		InstagramName:      "gettestuser",
		Memo:               "Get test memo",
		StartDate:          time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		Status:             1,
		DeleteHash:         false,
		CustomerType:       1,
	}
	err := repo.Create(ctx, testWI)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		filter  WordpressInstagramFilter
		wantErr bool
		check   func(t *testing.T, got *domain.WordpressInstagram)
	}{
		{
			name:    "IDで取得",
			filter:  WordpressInstagramFilter{ID: &testWI.ID},
			wantErr: false,
			check: func(t *testing.T, got *domain.WordpressInstagram) {
				assert.Equal(t, testWI.ID, got.ID)
				assert.Equal(t, testWI.Name, got.Name)
				assert.Equal(t, testWI.InstagramName, got.InstagramName)
			},
		},
		{
			name:    "InstagramIDで取得",
			filter:  WordpressInstagramFilter{InstagramID: &testWI.InstagramID},
			wantErr: false,
			check: func(t *testing.T, got *domain.WordpressInstagram) {
				assert.Equal(t, testWI.InstagramID, got.InstagramID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(ctx, tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestWordpressInstagramRepository_Update(t *testing.T) {
	repo := NewWordpressInstagramRepository(db)
	ctx := context.Background()

	// テストデータを作成
	testWI := &domain.WordpressInstagram{
		Name:               "Update Test Site",
		WordpressDomain:    "https://updatetest.example.com",
		WordpressSiteTitle: "Update Test Title",
		InstagramID:        "444555666",
		InstagramName:      "updatetestuser",
		Memo:               "Update test memo",
		StartDate:          time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		Status:             1,
		DeleteHash:         false,
		CustomerType:       1,
	}
	err := repo.Create(ctx, testWI)
	assert.NoError(t, err)

	t.Run("全フィールドを更新", func(t *testing.T) {
		testWI.Name = "Updated Name"
		testWI.WordpressDomain = "https://updated.example.com"
		testWI.WordpressSiteTitle = "Updated Title"
		testWI.InstagramID = "999888777"
		testWI.InstagramName = "updateduser"
		testWI.Memo = "Updated memo"
		testWI.Status = 2
		testWI.DeleteHash = true
		testWI.CustomerType = 2

		err := repo.Update(ctx, testWI, WordpressInstagramFilter{ID: &testWI.ID})
		assert.NoError(t, err)

		// 更新されたデータを確認
		got, err := repo.Get(ctx, WordpressInstagramFilter{ID: &testWI.ID})
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", got.Name)
		assert.Equal(t, "https://updated.example.com", got.WordpressDomain)
		assert.Equal(t, "Updated Title", got.WordpressSiteTitle)
		assert.Equal(t, "999888777", got.InstagramID)
		assert.Equal(t, "updateduser", got.InstagramName)
		assert.Equal(t, "Updated memo", got.Memo)
		assert.Equal(t, domain.Status(2), got.Status)
		assert.Equal(t, true, got.DeleteHash)
		assert.Equal(t, domain.CustomerType(2), got.CustomerType)
	})
}

func TestWordpressInstagramRepository_Delete(t *testing.T) {
	repo := NewWordpressInstagramRepository(db)
	ctx := context.Background()

	// テストデータを作成
	testWI := &domain.WordpressInstagram{
		Name:               "Delete Test Site",
		WordpressDomain:    "https://deletetest.example.com",
		WordpressSiteTitle: "Delete Test Title",
		InstagramID:        "777888999",
		InstagramName:      "deletetestuser",
		Memo:               "Delete test memo",
		StartDate:          time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
		Status:             1,
		DeleteHash:         false,
		CustomerType:       1,
	}
	err := repo.Create(ctx, testWI)
	assert.NoError(t, err)

	t.Run("削除成功", func(t *testing.T) {
		err := repo.Delete(ctx, WordpressInstagramFilter{ID: &testWI.ID})
		assert.NoError(t, err)

		// 削除されたことを確認
		var count int64
		db.Model(&model.WordpressInstagram{}).Where("id = ?", testWI.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestWordpressInstagramRepository_Exists(t *testing.T) {
	repo := NewWordpressInstagramRepository(db)
	ctx := context.Background()

	// テストデータを作成
	testWI := &domain.WordpressInstagram{
		Name:               "Exists Test Site",
		WordpressDomain:    "https://existstest.example.com",
		WordpressSiteTitle: "Exists Test Title",
		InstagramID:        "555666777",
		InstagramName:      "existstestuser",
		Memo:               "Exists test memo",
		StartDate:          time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		Status:             1,
		DeleteHash:         false,
		CustomerType:       1,
	}
	err := repo.Create(ctx, testWI)
	fmt.Println(testWI.ID)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		filter  WordpressInstagramFilter
		want    bool
		wantErr bool
	}{
		{
			name:    "存在するIDで検索",
			filter:  WordpressInstagramFilter{ID: &testWI.ID},
			want:    true,
			wantErr: false,
		},
		{
			name:    "存在しないIDで検索",
			filter:  WordpressInstagramFilter{ID: intPtr(99999)},
			want:    false,
			wantErr: false,
		},
		{
			name:    "存在するInstagramIDで検索",
			filter:  WordpressInstagramFilter{InstagramID: &testWI.InstagramID},
			want:    true,
			wantErr: false,
		},
		{
			name:    "存在しないInstagramIDで検索",
			filter:  WordpressInstagramFilter{InstagramID: strPtr("nonexistent")},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Exists(ctx, tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ヘルパー関数
func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
