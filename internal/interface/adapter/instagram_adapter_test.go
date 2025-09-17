package adapter

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/zuxt268/homing/internal/infrastructure/driver"
)

func TestInstagramAdapter_GetAccount(t *testing.T) {
	domain := os.Getenv("TEST_DOMAIN")
	accessToken := os.Getenv("TEST_ACCESS_TOKEN")
	fmt.Println("TEST_ACCESS_TOKEN:", accessToken)
	fmt.Println("TEST_DOMAIN:", domain)

	httpClient := &http.Client{}
	client := driver.NewClient(httpClient)

	adapter := NewInstagramAdapter(client)

	account, err := adapter.GetAccount(context.Background(), accessToken)
	if err != nil {
		t.Fatalf("GetAccount failed: %v", err)
	}
	fmt.Println("account:", account)

	for _, a := range account {
		fmt.Println("account:", a)
		posts, err := adapter.GetPosts(context.Background(), accessToken, a.InstagramAccountID)
		if err != nil {
			t.Fatalf("GetPosts failed: %v", err)
		}
		fmt.Println(posts)
	}
}
