package adapter

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zuxt268/homing/internal/domain/entity"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/dto/external"
)

type WordpressAdapter interface {
	Post(ctx context.Context, in external.WordpressPostInput) (*entity.Post, error)
	FileUpload(ctx context.Context, in external.WordpressFileUploadInput) (*external.WordpressFileUploadResponse, error)
}

func NewWordpressAdapter(
	httpDriver driver.HttpDriver,
	adminEmail string,
	secretPhrase string,
) WordpressAdapter {
	return &wordpressAdapter{
		httpDriver:   httpDriver,
		adminEmail:   adminEmail,
		secretPhrase: secretPhrase,
	}
}

type wordpressAdapter struct {
	httpDriver   driver.HttpDriver
	adminEmail   string
	secretPhrase string
}

func (a *wordpressAdapter) Post(ctx context.Context, in external.WordpressPostInput) (*entity.Post, error) {
	reqBody := external.WordpressPostPayload{
		Email:         a.adminEmail,
		Title:         in.Post.GetTitle(),
		Content:       in.Post.GetContent(),
		FeaturedMedia: in.FeaturedMediaID,
	}
	apiKey := in.Customer.GenerateAPIKey(a.secretPhrase)

	header, err := external.GetWordpressHeader(reqBody, apiKey)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(in.Customer.WordpressUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("rest_route", "/rodut/v1/create-post")
	u.RawQuery = q.Encode()

	endpoint := "https://" + u.String()

	resp, err := a.httpDriver.Post(ctx, endpoint, &reqBody, header)
	if err != nil {
		return nil, fmt.Errorf("記事の投稿に失敗: %w", err)
	}
	var postDto external.WordpressPostResponse
	if err := json.Unmarshal(resp, &postDto); err != nil {
		return nil, fmt.Errorf("JSONの変換に失敗: %w", err)
	}

	return &entity.Post{
		ID:           postDto.PostId,
		WordpressURL: postDto.PostUrl,
	}, nil
}

func (a *wordpressAdapter) FileUpload(ctx context.Context, in external.WordpressFileUploadInput) (*external.WordpressFileUploadResponse, error) {
	file, err := os.Open(in.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	fileName := filepath.Base(in.Path)
	apiKey := in.Customer.GenerateAPIKey(a.secretPhrase)

	// MIMEタイプを取得（Pythonの実装に合わせる）
	mimeType := mime.TypeByExtension(filepath.Ext(fileName))
	if mimeType == "" || !strings.HasPrefix(mimeType, "video/") {
		mimeType = "video/mp4" // Instagram動画はたいていmp4に寄せる
	}

	// HMAC署名を作成
	headers := signUploadHeaders(a.adminEmail, fileName, apiKey)

	// multipart/form-dataを作成
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// ファイルフィールドをMIMEタイプ付きで追加
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName))
	h.Set("Content-Type", mimeType)
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// emailフィールドを追加
	err = writer.WriteField("email", a.adminEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to write email field: %w", err)
	}
	defer func() {
		_ = writer.Close()
	}()

	// WordPressのアップロードURLを構築
	u, err := url.Parse(in.Customer.WordpressUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WordPress URL: %w", err)
	}

	q := u.Query()
	q.Set("rest_route", "/rodut/v1/upload-media")
	u.RawQuery = q.Encode()

	endpoint := "https://" + u.String()

	// HTTPリクエストを作成
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Timestamp", headers["X-Timestamp"])
	req.Header.Set("X-Signature", headers["X-Signature"])

	// リクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("upload failed with status code: %d", resp.StatusCode)
	}

	var uploadResponse external.WordpressFileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResponse); err != nil {
		return nil, fmt.Errorf("failed to decode upload response: %w", err)
	}

	return &uploadResponse, nil
}

// signUploadHeaders creates HMAC signature headers for file upload
// multipart/form-data 用: 署名対象は 'timestamp.email.filename'
func signUploadHeaders(email, filename, apiKeyHex string) map[string]string {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	message := fmt.Sprintf("%s.%s.%s", ts, email, filename)

	// Python implementation: hmac.new(api_key_hex.encode("utf-8"), message, hashlib.sha256)
	mac := hmac.New(sha256.New, []byte(apiKeyHex))
	mac.Write([]byte(message))
	signature := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"X-Timestamp": ts,
		"X-Signature": signature,
	}
}
