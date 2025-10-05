package external

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zuxt268/homing/internal/domain"
)

type WordpressPostPayload struct {
	Email         string `json:"email"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	FeaturedMedia int    `json:"featured_media"`
}

func GetWordpressHeader(payload any, apiKeyHex string) (map[string]string, error) {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	ts := fmt.Sprintf("%d", time.Now().Unix())

	var buf bytes.Buffer
	buf.WriteString(ts)
	buf.WriteByte('.')
	buf.Write(bodyBytes)
	message := buf.Bytes()

	// HMAC-SHA256 (鍵は hex文字列をそのまま使う)
	mac := hmac.New(sha256.New, []byte(apiKeyHex))
	mac.Write(message)
	signature := hex.EncodeToString(mac.Sum(nil))

	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"X-Timestamp":  ts,
		"X-Signature":  signature,
	}

	return headers, nil
}

type WordpressPostResponse struct {
	PostId  int    `json:"post_id"`
	PostUrl string `json:"post_url"`
	Message string `json:"message"`
}

type WordpressPostInput struct {
	WordpressInstagram domain.WordpressInstagram
	FeaturedMediaID    int
	Post               domain.InstagramPost
}

type WordpressFileUploadInput struct {
	Path               string
	WordpressInstagram domain.WordpressInstagram
}

type WordpressFileUploadPayload struct {
	Email string `json:"email"`
}

type WordpressFileUploadResponse struct {
	Id        int    `json:"id"`
	SourceUrl string `json:"source_url"`
	MimeType  string `json:"mime_type"`
}
