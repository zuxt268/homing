package adapter

import (
	"context"
	"net/http"
)

type FileDownloader interface {
	Download(ctx context.Context, url string, localPath string) error
}

type fileDownloader struct {
	httpClient *http.Client
}

func NewFileDownloader(client *http.Client) FileDownloader {
	return &fileDownloader{httpClient: client}
}

func (f *fileDownloader) Download(ctx context.Context, url string, localPath string) error {
	return nil
}
