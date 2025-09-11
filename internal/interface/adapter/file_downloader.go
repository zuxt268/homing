package adapter

import "net/http"

type FileDownloader interface {
	Download(ctx context.Context, url string, localPath string) error
}

type localFileDownloader struct {
	httpClient *http.Client
}

func NewFileDownloader(client *http.Client) FileDownloader {
	return &localFileDownloader{httpClient: client}
}
