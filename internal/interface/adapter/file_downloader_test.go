package adapter

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileDownloader_Download(t *testing.T) {
	u := "https://picsum.photos/200/300"
	downloader := NewFileDownloader()
	path, err := downloader.Download(context.Background(), u)
	assert.NoError(t, err)
	fmt.Println(path)

	//err = downloader.DeleteTempDirectory()
	assert.NoError(t, err)
}
