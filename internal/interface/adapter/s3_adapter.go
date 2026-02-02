package adapter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Adapter interface {
	UploadFromURL(ctx context.Context, sourceURL string) (publicURL string, err error)
}

type s3Adapter struct {
	client *s3.Client
	bucket string
	region string
	prefix string
}

func NewS3Adapter(bucket, region, prefix string) (S3Adapter, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("AWS設定読み込みエラー: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &s3Adapter{
		client: client,
		bucket: bucket,
		region: region,
		prefix: prefix,
	}, nil
}

func (a *s3Adapter) UploadFromURL(ctx context.Context, sourceURL string) (string, error) {
	// URLからメディアをダウンロード
	req, err := http.NewRequestWithContext(ctx, "GET", sourceURL, nil)
	if err != nil {
		return "", fmt.Errorf("リクエスト作成エラー: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ダウンロードエラー: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ダウンロード失敗 (ステータス: %d)", resp.StatusCode)
	}

	// URLから拡張子を取得
	ext := extFromURL(sourceURL)

	// Content-Type をレスポンスヘッダーから取得（動画・画像両対応）
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = guessContentType(ext)
	}

	// 拡張子が取れなかった場合は Content-Type から推定
	if ext == "" {
		ext = extFromContentType(contentType)
	}

	key := fmt.Sprintf("%s%s%s", a.prefix, uuid.New().String(), ext)

	// レスポンスボディをメモリに読み込み（Content-Length が必要なため）
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("レスポンス読み込みエラー: %v", err)
	}

	_, err = a.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(a.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(data))),
	})
	if err != nil {
		return "", fmt.Errorf("S3アップロードエラー: %v", err)
	}

	publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", a.bucket, a.region, key)

	return publicURL, nil
}

func extFromURL(sourceURL string) string {
	u, err := url.Parse(sourceURL)
	if err != nil {
		return ""
	}
	ext := path.Ext(u.Path)
	// クエリパラメータ付きの拡張子をクリーンアップ
	if idx := strings.Index(ext, "?"); idx != -1 {
		ext = ext[:idx]
	}
	return ext
}

func guessContentType(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	default:
		return "application/octet-stream"
	}
}

func extFromContentType(contentType string) string {
	ct := strings.ToLower(strings.Split(contentType, ";")[0])
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "video/quicktime":
		return ".mov"
	case "video/x-msvideo":
		return ".avi"
	default:
		return ".bin"
	}
}