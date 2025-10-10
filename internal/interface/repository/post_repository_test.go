package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"github.com/zuxt268/homing/internal/interface/util"
)

func TestPostRepository_ExistPost(t *testing.T) {
	repo := NewPostRepository(db)
	ctx := context.Background()

	t.Run("存在する投稿をMediaIDで検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			MediaID: util.Pointer("media_001"),
		})
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("存在する投稿をIDで検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			ID: util.Pointer(1),
		})
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("存在する投稿をCustomerIDで検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			CustomerID: util.Pointer(1),
		})
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("存在する投稿をTimestampで検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			Timestamp: util.Pointer("1640995200"),
		})
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("存在する投稿をMediaURLで検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			MediaURL: util.Pointer("https://instagram.com/p/sample1.jpg"),
		})
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("存在する投稿をPermalinkで検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			Permalink: util.Pointer("https://instagram.com/p/sample1/"),
		})
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("存在する投稿をWordpressLinkで検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			WordpressLink: util.Pointer("https://tanaka-blog.com/post/instagram-1/"),
		})
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("複数の条件で検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			CustomerID: util.Pointer(1),
			MediaID:    util.Pointer("media_001"),
		})
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("存在しないMediaIDで検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			MediaID: util.Pointer("non_existent_media"),
		})
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("存在しないIDで検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			ID: util.Pointer(9999),
		})
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("存在しないCustomerIDで検索", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			CustomerID: util.Pointer(9999),
		})
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("複数の条件で検索（一致しない）", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{
			CustomerID: util.Pointer(1),
			MediaID:    util.Pointer("media_003"), // Customer 1には属さないMediaID
		})
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("空のフィルターで検索（すべての投稿が対象）", func(t *testing.T) {
		exists, err := repo.ExistPost(ctx, PostFilter{})
		assert.NoError(t, err)
		assert.True(t, exists) // サンプルデータが存在するため
	})
}

func TestPostRepository_SavePost(t *testing.T) {
	repo := NewPostRepository(db)
	ctx := context.Background()

	t.Run("新しい投稿を保存", func(t *testing.T) {
		newPost := &model.Post{
			MediaID:       "test_media_001",
			CustomerID:    1,
			Timestamp:     "2025-09-29T10:19:24+0000",
			MediaURL:      "https://instagram.com/p/test001.jpg",
			Permalink:     "https://instagram.com/p/test001/",
			WordpressLink: "https://tanaka-blog.com/post/test-1/",
			CreatedAt:     time.Now(),
		}

		err := repo.CreatePost(ctx, newPost)
		assert.NoError(t, err)
		assert.NotZero(t, newPost.ID) // IDが自動で設定されることを確認

		// 保存された投稿が実際に存在することを確認
		exists, err := repo.ExistPost(ctx, PostFilter{
			MediaID: util.Pointer("test_media_001"),
		})
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("必要な情報のみで投稿を保存", func(t *testing.T) {

		minimumPost := &model.Post{
			MediaID:    "test_media_002",
			CustomerID: 2,
			Timestamp:  "2025-09-29T10:19:24+0000",
			CreatedAt:  time.Now(),
		}

		err := repo.CreatePost(ctx, minimumPost)
		assert.NoError(t, err)
		assert.NotZero(t, minimumPost.ID)

		// 保存された投稿が実際に存在することを確認
		exists, err := repo.ExistPost(ctx, PostFilter{
			MediaID: util.Pointer("test_media_002"),
		})
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("異なるCustomerIDで投稿を保存", func(t *testing.T) {

		posts := []*model.Post{
			{
				MediaID:    "test_media_003",
				CustomerID: 3,
				Timestamp:  "2025-09-29T10:19:24+0000",
				CreatedAt:  time.Now(),
			},
			{
				MediaID:    "test_media_004",
				CustomerID: 4,
				Timestamp:  "2025-09-29T10:19:24+0000",
				CreatedAt:  time.Now(),
			},
		}

		for _, post := range posts {
			err := repo.CreatePost(ctx, post)
			assert.NoError(t, err)
			assert.NotZero(t, post.ID)
		}

		// 保存された投稿が実際に存在することを確認
		for _, post := range posts {
			exists, err := repo.ExistPost(ctx, PostFilter{
				MediaID: util.Pointer(post.MediaID),
			})
			assert.NoError(t, err)
			assert.True(t, exists)
		}
	})
}

func TestPostRepository_Integration(t *testing.T) {
	repo := NewPostRepository(db)
	ctx := context.Background()

	t.Run("保存と存在確認の統合テスト", func(t *testing.T) {
		testPost := &model.Post{
			MediaID:       "integration_test_001",
			CustomerID:    1,
			Timestamp:     "2025-09-29T10:19:24+0000",
			MediaURL:      "https://instagram.com/p/integration001.jpg",
			Permalink:     "https://instagram.com/p/integration001/",
			WordpressLink: "https://tanaka-blog.com/post/integration-1/",
			CreatedAt:     time.Now(),
		}

		// 1. 投稿が存在しないことを確認
		exists, err := repo.ExistPost(ctx, PostFilter{
			MediaID: util.Pointer(testPost.MediaID),
		})
		require.NoError(t, err)
		assert.False(t, exists)

		// 2. 投稿を保存
		err = repo.CreatePost(ctx, testPost)
		require.NoError(t, err)
		require.NotZero(t, testPost.ID)

		// 3. 投稿が存在することを確認
		exists, err = repo.ExistPost(ctx, PostFilter{
			MediaID: util.Pointer(testPost.MediaID),
		})
		require.NoError(t, err)
		assert.True(t, exists)

		// 4. 様々な条件で存在確認
		exists, err = repo.ExistPost(ctx, PostFilter{
			ID: util.Pointer(testPost.ID),
		})
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = repo.ExistPost(ctx, PostFilter{
			CustomerID: util.Pointer(testPost.CustomerID),
			MediaID:    util.Pointer(testPost.MediaID),
		})
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("サンプルデータの整合性チェック", func(t *testing.T) {
		// サンプルデータが正しく存在することを確認
		expectedSampleData := []struct {
			mediaID    string
			customerID int
			timestamp  string
		}{
			{"media_001", 1, "1640995200"},
			{"media_002", 1, "1641081600"},
			{"media_003", 2, "1644753600"},
			{"media_004", 2, "1644840000"},
			{"media_005", 3, "1647334800"},
			{"media_006", 4, "1649154900"},
		}

		for _, expected := range expectedSampleData {
			exists, err := repo.ExistPost(ctx, PostFilter{
				MediaID:    util.Pointer(expected.mediaID),
				CustomerID: util.Pointer(expected.customerID),
				Timestamp:  util.Pointer(expected.timestamp),
			})
			require.NoError(t, err)
			assert.True(t, exists, "Sample data not found: MediaID=%s, CustomerID=%d", expected.mediaID, expected.customerID)
		}
	})

	t.Run("Customer別投稿数の確認", func(t *testing.T) {
		// 各Customerの投稿数を確認
		expectedCounts := map[int]int{
			1: 2, // 田中太郎
			2: 2, // 山田花子
			3: 1, // 佐藤次郎
			4: 1, // 鈴木美香
		}

		for customerID := range expectedCounts {
			exists, err := repo.ExistPost(ctx, PostFilter{
				CustomerID: util.Pointer(customerID),
			})
			require.NoError(t, err)
			assert.True(t, exists, "No posts found for customer %d", customerID)

			// 実際の件数チェックは簡易的に存在確認のみ（件数取得メソッドがないため）
		}
	})
}
