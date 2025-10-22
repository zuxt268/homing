package adapter

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/zuxt268/homing/internal/infrastructure/driver"
)

func TestInstagramAdapter_GetPostsAll(t *testing.T) {
	client := &http.Client{}
	httpDriver := driver.NewClient(client)
	adapter := &instagramAdapter{
		httpDriver:   httpDriver,
		clientID:     os.Getenv("CLIENT_ID"),
		clientSecret: os.Getenv("CLIENT_SECRET"),
	}
	data, err := adapter.GetPostsAll(context.Background(), os.Getenv("TOKEN"), os.Getenv("INSTAGRAM_ID"))
	if err != nil {
		t.Errorf("GetPostsAll returned an error %v", err)
	}
	fmt.Println(len(data))
}
