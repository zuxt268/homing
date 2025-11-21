package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/mybusinessaccountmanagement/v1"
	"google.golang.org/api/mybusinessbusinessinformation/v1"
	"google.golang.org/api/option"
)

// 保存用トークンファイル名
const tokenFile = "./credentials/token.json"

// トークンを保存
func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("token 保存エラー: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
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

// token.json がなければブラウザからOAuth認証
func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		// 初回認証
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		fmt.Println("以下のURLにアクセスしてGoogleログインしてください:\n", authURL)

		var code string
		fmt.Print("認証コードを貼り付けてください: ")
		fmt.Scan(&code)

		tok, err = config.Exchange(context.Background(), code)
		if err != nil {
			log.Fatalf("認証コード交換エラー: %v", err)
		}
		saveToken(tokenFile, tok)
		fmt.Println("token.json 保存完了")
	}

	// トークンソースを作成（自動更新機能付き）
	tokenSource := config.TokenSource(context.Background(), tok)

	// トークンを取得（必要に応じて自動更新される）
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("トークン取得エラー: %v", err)
	}

	// トークンが更新された場合は保存
	if newToken.AccessToken != tok.AccessToken {
		saveToken(tokenFile, newToken)
		fmt.Println("トークン自動更新・保存完了")
	}

	return oauth2.NewClient(context.Background(), tokenSource)
}

func main() {
	ctx := context.Background()

	// 👇 credentials.json を読み込む
	b, err := os.ReadFile("credentials/client_secret.json")
	if err != nil {
		log.Fatalf("credentials.json 読み込みエラー: %v", err)
	}

	// GBP 用スコープ
	config, err := google.ConfigFromJSON(b,
		"https://www.googleapis.com/auth/business.manage",
	)
	if err != nil {
		log.Fatalf("OAuth設定エラー: %v", err)
	}

	client := getClient(config)

	// GBP API クライアント（アカウント管理用）
	accountSvc, err := mybusinessaccountmanagement.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("アカウント管理API 初期化エラー: %v", err)
	}

	// GBP API クライアント（ビジネス情報用）
	businessSvc, err := mybusinessbusinessinformation.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("ビジネス情報API 初期化エラー: %v", err)
	}

	// アカウント一覧を取得
	resp, err := accountSvc.Accounts.List().Do()
	if err != nil {
		log.Fatalf("アカウント取得エラー: %v", err)
	}

	fmt.Println("==== GBP アカウント一覧 ====")
	for _, acct := range resp.Accounts {
		fmt.Println("アカウント名:", acct.Name)
		fmt.Println("表示名:", acct.AccountName)
		fmt.Println()

		// 各アカウント配下のロケーション（ビジネス）一覧を取得
		fmt.Println("  --- ビジネス一覧 ---")
		locResp, err := businessSvc.Accounts.Locations.List(acct.Name).
			ReadMask("name,title,storefrontAddress,profile").
			Do()
		if err != nil {
			log.Printf("  ビジネス取得エラー: %v\n", err)
			continue
		}

		if locResp == nil || len(locResp.Locations) == 0 {
			fmt.Println("  （ビジネスなし）")
		} else {
			for _, loc := range locResp.Locations {
				fmt.Println("  ビジネス名:", loc.Name)
				if loc.Title != "" {
					fmt.Println("  店舗名:", loc.Title)
				}
				if loc.StorefrontAddress != nil {
					fmt.Printf("  住所: %s\n", formatAddress(loc.StorefrontAddress))
				}

				// プロフィール情報（説明文）
				if loc.Profile != nil && loc.Profile.Description != "" {
					fmt.Printf("  説明: %s\n", loc.Profile.Description)
				}

				// メディア（写真）と投稿を取得
				fetchMediaAndLocalPosts(client, acct.Name, loc.Name)

				fmt.Println("  ---")
			}
		}
		fmt.Println("========================")
	}
}

// 住所を整形
func formatAddress(addr *mybusinessbusinessinformation.PostalAddress) string {
	return fmt.Sprintf("%s %s %s %s",
		addr.PostalCode,
		addr.AdministrativeArea,
		addr.Locality,
		addr.AddressLines,
	)
}

// locationNameからlocation IDを抽出
func extractLocationID(locationName string) string {
	// "locations/12345" -> "12345"
	parts := strings.Split(locationName, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return locationName
}

// メディア（画像）を取得（My Business API v4経由）
func fetchMediaAndLocalPosts(client *http.Client, accountName, locationName string) {
	fmt.Println("\n  === メディア（写真）===")

	// My Business API v4のエンドポイントを使用
	locationID := extractLocationID(locationName)
	mediaURL := fmt.Sprintf("https://mybusiness.googleapis.com/v4/%s/locations/%s/media", accountName, locationID)

	resp, err := client.Get(mediaURL)
	if err != nil {
		fmt.Printf("  メディア取得エラー: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("  ステータスコード: %d\n", resp.StatusCode)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err == nil {
			// mediaItems を探す
			if mediaItems, ok := result["mediaItems"].([]interface{}); ok && len(mediaItems) > 0 {
				fmt.Printf("  写真数: %d件\n", len(mediaItems))
				for i, item := range mediaItems {
					mediaItem := item.(map[string]interface{})
					if name, ok := mediaItem["name"].(string); ok {
						fmt.Printf("  [%d] メディア名: %s\n", i+1, name)
					}
					if url, ok := mediaItem["googleUrl"].(string); ok {
						fmt.Printf("      URL: %s\n", url)
					}
					if sourceUrl, ok := mediaItem["sourceUrl"].(string); ok {
						fmt.Printf("      ソースURL: %s\n", sourceUrl)
					}
					if mediaFormat, ok := mediaItem["mediaFormat"].(string); ok {
						fmt.Printf("      形式: %s\n", mediaFormat)
					}
					if createTime, ok := mediaItem["createTime"].(string); ok {
						fmt.Printf("      作成日時: %s\n", createTime)
					}
				}
			} else {
				fmt.Println("  （写真なし）")
				if len(body) < 500 {
					fmt.Printf("  レスポンス: %s\n", string(body))
				}
			}
		} else {
			fmt.Printf("  JSONパースエラー: %v\n", err)
		}
	} else {
		fmt.Printf("  エラー: %s\n", string(body)[:min(500, len(body))])
	}
}
