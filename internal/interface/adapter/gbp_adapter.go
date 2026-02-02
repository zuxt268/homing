package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zuxt268/homing/internal/interface/dto/external"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/mybusinessbusinessinformation/v1"
	"google.golang.org/api/option"
)

type Business struct {
	Name               string
	Title              string
	PostalCode         string
	AdministrativeArea string
	Locality           string
	AddressLines       []string
	Description        string
}

type GbpAdapter interface {
	GetAllBusinesses(ctx context.Context, accountName string) ([]Business, error)
	UploadMedia(ctx context.Context, accountName, businessName string, sourceURL string) (*external.GoogleBusinessMediaUploadResponse, error)
	GetBusiness(ctx context.Context, businessName string) (Business, error)
}

type gbpAdapter struct {
	client *http.Client
}

func NewGbpAdapter(credentialsData []byte) (GbpAdapter, error) {
	ctx := context.Background()

	// GBP 用スコープ
	config, err := google.ConfigFromJSON(credentialsData,
		"https://www.googleapis.com/auth/business.manage",
	)
	if err != nil {
		return nil, fmt.Errorf("OAuth設定エラー: %v", err)
	}

	client := getOAuthClient(ctx, config)

	return &gbpAdapter{
		client: client,
	}, nil
}

func (a *gbpAdapter) GetAllBusinesses(ctx context.Context, accountName string) ([]Business, error) {
	businessSvc, err := mybusinessbusinessinformation.NewService(ctx, option.WithHTTPClient(a.client))
	if err != nil {
		return nil, fmt.Errorf("ビジネス情報API 初期化エラー: %v", err)
	}

	pageToken := ""
	var businesses []Business

	for {
		call := businessSvc.Accounts.Locations.List(accountName).
			ReadMask("name,title,storefrontAddress,profile").
			PageSize(100)

		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		locResp, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("ビジネス取得エラー: %v", err)
		}

		for _, loc := range locResp.Locations {
			business := Business{
				Name:  loc.Name,
				Title: loc.Title,
			}

			if loc.StorefrontAddress != nil {
				business.PostalCode = loc.StorefrontAddress.PostalCode
				business.AdministrativeArea = loc.StorefrontAddress.AdministrativeArea
				business.Locality = loc.StorefrontAddress.Locality
				business.AddressLines = loc.StorefrontAddress.AddressLines
			}

			if loc.Profile != nil {
				business.Description = loc.Profile.Description
			}

			businesses = append(businesses, business)
		}

		if locResp.NextPageToken == "" {
			break
		}
		pageToken = locResp.NextPageToken
	}

	return businesses, nil
}

func (a *gbpAdapter) UploadMedia(ctx context.Context, accountName, businessName string, sourceURL string) (*external.GoogleBusinessMediaUploadResponse, error) {
	locationID := extractLocationID(businessName)
	parent := fmt.Sprintf("%s/locations/%s", accountName, locationID)

	// sourceUrl 方式で media.create を呼ぶ
	createURL := fmt.Sprintf("https://mybusiness.googleapis.com/v4/%s/media", parent)
	reqBody := map[string]any{
		"mediaFormat": "PHOTO",
		"locationAssociation": map[string]string{
			"category": "ADDITIONAL",
		},
		"sourceUrl": sourceURL,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("JSON作成エラー: %v", err)
	}

	req, err := http.NewRequest("POST", createURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("media.createエラー: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み込みエラー: %v", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("アップロード失敗 (ステータス: %d): %s", resp.StatusCode, string(body))
	}

	var uploadResponse external.GoogleBusinessMediaUploadResponse
	if err := json.Unmarshal(body, &uploadResponse); err != nil {
		return nil, fmt.Errorf("レスポンスパースエラー: %v", err)
	}

	return &uploadResponse, nil
}

func (a *gbpAdapter) GetBusiness(ctx context.Context, businessName string) (Business, error) {
	businessSvc, err := mybusinessbusinessinformation.NewService(ctx, option.WithHTTPClient(a.client))
	if err != nil {
		return Business{}, err
	}
	fmt.Println(businessName)
	locResp, err := businessSvc.Locations.
		Get(businessName).
		ReadMask("name,title,storefrontAddress,profile").
		Do()
	if err != nil {
		return Business{}, err
	}
	business := Business{
		Name:  locResp.Name,
		Title: locResp.Title,
	}
	if locResp.StorefrontAddress != nil {
		business.PostalCode = locResp.StorefrontAddress.PostalCode
		business.AdministrativeArea = locResp.StorefrontAddress.AdministrativeArea
		business.Locality = locResp.StorefrontAddress.Locality
		business.AddressLines = locResp.StorefrontAddress.AddressLines
	}

	if locResp.Profile != nil {
		business.Description = locResp.Profile.Description
	}
	return business, nil
}

// ヘルパー関数: locationNameからlocation IDを抽出
func extractLocationID(locationName string) string {
	// "accounts/xxx/locations/12345" -> "12345"
	parts := filepath.Base(locationName)
	return parts
}

// ヘルパー関数: OAuth2クライアント取得
func getOAuthClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tokenFile := "./credentials/token.json"

	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		// 初回認証
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		log.Printf("以下のURLにアクセスしてGoogleログインしてください:\n%s\n", authURL)

		var code string
		fmt.Print("認証コードを貼り付けてください: ")
		fmt.Scan(&code)

		tok, err = config.Exchange(ctx, code)
		if err != nil {
			log.Fatalf("認証コード交換エラー: %v", err)
		}
		saveToken(tokenFile, tok)
		log.Println("token.json 保存完了")
	}

	// トークンソースを作成（自動更新機能付き）
	tokenSource := config.TokenSource(ctx, tok)

	// トークンを取得（必要に応じて自動更新される）
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("トークン取得エラー: %v", err)
	}

	// トークンが更新された場合は保存
	if newToken.AccessToken != tok.AccessToken {
		saveToken(tokenFile, newToken)
		log.Println("トークン自動更新・保存完了")
	}

	return oauth2.NewClient(ctx, tokenSource)
}

// token.json を読む
func tokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// トークンを保存
func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("token 保存エラー: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
