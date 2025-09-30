package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zuxt268/homing/internal/interface/util"
)

func TestCustomerRepository_GetCustomer(t *testing.T) {
	repo := NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("存在するカスタマーを取得", func(t *testing.T) {
		customer, err := repo.GetCustomer(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, 1, customer.ID)
		assert.Equal(t, "田中太郎", customer.Name)
		assert.Equal(t, "https://tanaka-blog.com", customer.WordpressUrl)
		assert.Equal(t, "facebook_token_123", customer.FacebookToken)
	})

	t.Run("存在しないカスタマーの場合はnilを返す", func(t *testing.T) {
		customer, err := repo.GetCustomer(ctx, 9999)
		assert.NoError(t, err)
		assert.Nil(t, customer)
	})

	t.Run("無効なIDの場合", func(t *testing.T) {
		customer, err := repo.GetCustomer(ctx, -1)
		assert.NoError(t, err)
		assert.Nil(t, customer)
	})
}

func TestCustomerRepository_FindAllCustomers(t *testing.T) {
	repo := NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("フィルターなしで全カスタマーを取得", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{})
		assert.NoError(t, err)
		assert.Len(t, customers, 4)

		// IDの昇順で並んでいることを確認
		expectedNames := []string{"田中太郎", "山田花子", "佐藤次郎", "鈴木美香"}
		for i, customer := range customers {
			assert.Equal(t, i+1, customer.ID)
			assert.Equal(t, expectedNames[i], customer.Name)
		}
	})

	t.Run("名前でフィルター", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			Name: util.Pointer("田中太郎"),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, "田中太郎", customers[0].Name)
		assert.Equal(t, 1, customers[0].ID)
	})

	t.Run("IDでフィルター", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			ID: util.Pointer(2),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, "山田花子", customers[0].Name)
		assert.Equal(t, 2, customers[0].ID)
	})

	t.Run("メールアドレスでフィルター", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			Email: util.Pointer("sato@example.com"),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, "佐藤次郎", customers[0].Name)
		assert.Equal(t, 3, customers[0].ID)
	})

	t.Run("WordpressURLでフィルター", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			WordpressURL: util.Pointer("https://suzuki-beauty.com"),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, "鈴木美香", customers[0].Name)
		assert.Equal(t, 4, customers[0].ID)
	})

	t.Run("FacebookTokenでフィルター", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			FacebookToken: util.Pointer("facebook_token_456"),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, "山田花子", customers[0].Name)
		assert.Equal(t, 2, customers[0].ID)
	})

	t.Run("InstagramBusinessAccountIDでフィルター", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			InstagramBusinessAccountID: util.Pointer("ig_business_789"),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, "佐藤次郎", customers[0].Name)
		assert.Equal(t, 3, customers[0].ID)
	})

	t.Run("InstagramBusinessAccountNameでフィルター", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			InstagramBusinessAccountName: util.Pointer("yamada_fashion"),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, "山田花子", customers[0].Name)
		assert.Equal(t, 2, customers[0].ID)
	})

	t.Run("InstagramTokenStatusでフィルター", func(t *testing.T) {
		// トークンステータスが1のカスタマーを検索
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			InstagramTokenStatus: util.Pointer(1),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 3) // 田中、山田、鈴木

		expectedNames := []string{"田中太郎", "山田花子", "鈴木美香"}
		for i, customer := range customers {
			assert.Equal(t, expectedNames[i], customer.Name)
		}
	})

	t.Run("DeleteHashでフィルター", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			DeleteHash: util.Pointer(false),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 4) // 全員がfalse
	})

	t.Run("PaymentTypeでフィルター", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			PaymentType: util.Pointer("monthly"),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 2) // 田中、鈴木

		expectedNames := []string{"田中太郎", "鈴木美香"}
		for i, customer := range customers {
			assert.Equal(t, expectedNames[i], customer.Name)
		}
	})

	t.Run("複数フィルターの組み合わせ", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			PaymentType:          util.Pointer("monthly"),
			InstagramTokenStatus: util.Pointer(1),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 2) // 田中、鈴木
	})

	t.Run("存在しない条件でフィルター", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{
			Name: util.Pointer("存在しない名前"),
		})
		assert.NoError(t, err)
		assert.Len(t, customers, 0)
	})
}

func TestCustomerRepository_Integration(t *testing.T) {
	repo := NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("全ての基本機能の統合テスト", func(t *testing.T) {
		// 1. 全カスタマーを取得
		allCustomers, err := repo.FindAllCustomers(ctx, CustomerFilter{})
		require.NoError(t, err)
		require.Len(t, allCustomers, 4)

		// 2. 各カスタマーを個別に取得して比較
		for _, customer := range allCustomers {
			individualCustomer, err := repo.GetCustomer(ctx, customer.ID)
			require.NoError(t, err)
			require.NotNil(t, individualCustomer)

			assert.Equal(t, customer.ID, individualCustomer.ID)
			assert.Equal(t, customer.Name, individualCustomer.Name)
			assert.Equal(t, customer.WordpressUrl, individualCustomer.WordpressUrl)
			assert.Equal(t, customer.FacebookToken, individualCustomer.FacebookToken)
		}
	})

	t.Run("データの整合性チェック", func(t *testing.T) {
		customers, err := repo.FindAllCustomers(ctx, CustomerFilter{})
		require.NoError(t, err)

		// サンプルデータの内容を確認
		expectedData := map[int]struct {
			name               string
			wordpressURL       string
			facebookToken      string
			instagramAccountID string
		}{
			1: {"田中太郎", "https://tanaka-blog.com", "facebook_token_123", "ig_business_123"},
			2: {"山田花子", "https://yamada-store.com", "facebook_token_456", "ig_business_456"},
			3: {"佐藤次郎", "https://sato-cafe.com", "", "ig_business_789"}, // FacebookTokenはNULL
			4: {"鈴木美香", "https://suzuki-beauty.com", "facebook_token_789", "ig_business_012"},
		}

		for _, customer := range customers {
			expected := expectedData[customer.ID]
			assert.Equal(t, expected.name, customer.Name)
			assert.Equal(t, expected.wordpressURL, customer.WordpressUrl)
			assert.Equal(t, expected.facebookToken, customer.FacebookToken)
		}
	})
}
