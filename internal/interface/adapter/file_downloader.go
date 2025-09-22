package adapter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type FileDownloader interface {
	Download(ctx context.Context, url string) (string, error)
	MakeTempDirectory() error
	DeleteTempDirectory() error
}

type fileDownloader struct {
	httpClient *http.Client
	tempDir    string
}

func NewFileDownloader(client *http.Client) FileDownloader {
	return &fileDownloader{
		httpClient: client,
	}
}

func (f *fileDownloader) Download(ctx context.Context, urlStr string) (string, error) {
	if f.tempDir == "" {
		if err := f.MakeTempDirectory(); err != nil {
			return "", err
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	fileName := path.Base(u.Path)
	if fileName == "." || fileName == "/" {
		fileName = "downloaded_file"
	}

	// セパレータを削除（安全なファイル名に変換）
	safeFileName := filepath.Base(fileName)
	filePath := filepath.Join(f.tempDir, safeFileName)

	// os.CreateTemp は `*` を末尾に含めないといけない
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return file.Name(), nil
}

func (f *fileDownloader) MakeTempDirectory() error {
	tempDir, err := os.MkdirTemp("", "homing_download_*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	f.tempDir = tempDir
	return nil
}

func (f *fileDownloader) DeleteTempDirectory() error {
	if f.tempDir == "" {
		return nil
	}

	err := os.RemoveAll(f.tempDir)
	if err != nil {
		return fmt.Errorf("failed to delete temp directory: %w", err)
	}

	f.tempDir = ""
	return nil
}
