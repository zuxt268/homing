package domain

import (
	"fmt"
	"testing"
)

func TestInstagramPost(t *testing.T) {
	post := InstagramPost{
		Timestamp: "2026-02-06T03:01:50+0000",
	}
	fmt.Println(post.GetPostDate())
}
