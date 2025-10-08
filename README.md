# Homing

**Instagram投稿を自動的にWordPressへ連携するバックエンドサービス**

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://golang.org/)
[![Echo](https://img.shields.io/badge/Echo-v4-00ADD8)](https://echo.labstack.com/)
[![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white)](https://www.mysql.com/)

## 🎯 プロジェクトの背景と目的

SNSマーケティングを活用する企業にとって、InstagramからWebサイトへのコンテンツ連携は手作業では時間がかかり、ヒューマンエラーも発生しやすい課題でした。本プロジェクトは、この課題を解決するために開発した**マルチテナント対応のSaaS型自動連携システム**です。

### ビジネス価値

- **作業時間の削減**: 手作業での投稿作業を完全自動化（1投稿あたり5分→0分）
- **リアルタイム性**: 定期実行により、新規投稿を即座にWebサイトへ反映
- **スケーラビリティ**: 複数顧客・複数アカウントを一元管理
- **運用の安定性**: エラー通知とリトライ機能により、連携漏れを防止

## 📋 概要

Homingは、複数の顧客のInstagramビジネスアカウントから投稿を取得し、自動的にWordPressサイトへ連携するGoアプリケーションです。**クリーンアーキテクチャ**を採用し、保守性と拡張性を重視した設計になっています。

### 主な機能

- ✅ Instagram Graph APIからの投稿自動取得
- ✅ WordPressへの画像/動画アップロードと記事投稿
- ✅ 複数顧客・複数アカウント対応（マルチテナント）
- ✅ 連携開始日以降の投稿のみ同期（柔軟なフィルタリング）
- ✅ 20件並列処理による同期（セマフォパターン）
- ✅ Slack通知によるエラーアラート・成功通知
- ✅ 重複投稿の防止（冪等性の保証）
- ✅ Graceful Shutdown対応
- ✅ Docker/Docker Compose対応

## 💡 技術的なアピールポイント

### 1. クリーンアーキテクチャの実践

**Domain Driven Design (DDD)** に基づいた4層アーキテクチャを採用し、ビジネスロジックと技術的関心事を分離。

- **テスタビリティ**: 依存性注入により、モックを使った単体テストが容易
- **保守性**: 各層の責務が明確で、変更の影響範囲が限定的
- **拡張性**: 新しい外部サービス連携も、Adapterパターンで容易に追加可能

### 2. 複数の外部API統合

3つの異なる外部APIを統合し、エンドツーエンドの自動化を実現。

- **Instagram Graph API**: OAuth認証、ページング処理
- **WordPress REST API**: マルチパートファイルアップロード、Basic認証
- **Slack Incoming Webhook**: リアルタイムアラート

### 3. データ整合性の保証

- **重複防止**: メディアIDベースの冪等性チェック
- **トランザクション管理**: GORMを使用したDBトランザクション
- **エラーハンドリング**: 段階的なリトライとSlack通知

### 4. プロダクションレディな実装

- **環境変数管理**: direnv / envconfig
- **マイグレーション**: sql-migrateによるスキーマバージョン管理
- **API仕様書**: Swagger/OpenAPIによる自動生成
- **ホットリロード**: Airによる開発効率向上
- **統合テスト**: Testcontainersで実際のMySQLを使用したテスト
- **Docker対応**: Docker ComposeによるコンテナベースのデプロイとMakefileによる自動化

### 5. コード品質

- **総コード行数**: 4,500行以上（Go）
- **テストカバレッジ**: リポジトリ層を中心にテスト実装
- **Go標準スタイル**: gofmtに準拠

## アーキテクチャ

クリーンアーキテクチャ(DDD)を採用しています。

```
internal/
├── domain/          # ドメインモデル（Customer, Post, WordpressInstagram, Instagram）
├── usecase/         # ビジネスロジック
│   ├── customer_usecase.go                # 顧客同期ロジック
│   ├── token_usecase.go                   # トークン管理ロジック
│   └── wordpress_instagram_usecase.go     # WordPress-Instagram連携管理
├── interface/       # 外部とのインターフェース
│   ├── adapter/     # 外部API連携
│   │   ├── instagram_adapter.go           # Instagram Graph API
│   │   ├── wordpress_adapter.go           # WordPress REST API
│   │   ├── slack.go                       # Slack通知
│   │   └── file_downloader.go             # ファイルダウンロード
│   ├── handler/     # HTTPハンドラー（APIエンドポイント）
│   ├── repository/  # データベースアクセス
│   │   ├── customer_repository.go
│   │   ├── post_repository.go
│   │   ├── token_repository.go
│   │   └── wordpress_instagram_repository.go
│   ├── dto/         # データ転送オブジェクト
│   │   ├── req/     # リクエストDTO
│   │   ├── res/     # レスポンスDTO
│   │   ├── model/   # データベースモデル
│   │   └── external/ # 外部API用DTO
│   └── util/        # ユーティリティ関数
├── infrastructure/  # インフラ層
│   ├── database/    # DB接続・マイグレーション
│   ├── driver/      # HTTPクライアント
│   └── server/      # Echoサーバー設定
├── config/          # 環境変数管理
└── di/              # 依存性注入
```

## 技術スタック

- **Go**: 1.25.0
- **Webフレームワーク**: Echo v4
- **ORM**: GORM
- **データベース**: MySQL 8.0
- **マイグレーション**: sql-migrate
- **ホットリロード**: Air
- **API仕様**: Swagger/OpenAPI
- **テスト**: Testify, Testcontainers

## セットアップ

### 前提条件

- Go 1.25.0以上
- MySQL 8.0以上
- direnv（推奨）

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd homing
```

### 2. 環境変数の設定

`.envrc`ファイルを作成し、以下の環境変数を設定してください。

```bash
export ADDRESS=:8090
export SECRET_PHRASE=your_secret_phrase
export ADMIN_EMAIL=admin@example.com
export SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
export DB_USER=root
export DB_PASSWORD=your_password
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=homing_db
```

direnvを使用している場合:
```bash
direnv allow
```

### 3. データベースの作成

```bash
mysql -u root -p
CREATE DATABASE homing_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. マイグレーションの実行

```bash
# sql-migrateのインストール
go install github.com/rubenv/sql-migrate/...@latest

# マイグレーション実行
sql-migrate up
```

### 5. 依存関係のインストール

```bash
go mod download
```

### 6. アプリケーションの起動

#### 開発環境（ホットリロード）

```bash
# Airのインストール
go install github.com/air-verse/air@latest

# 起動
make dev
```

#### 本番環境

```bash
go build -o homing ./cmd/homing/main.go
./homing
```

## API仕様

サーバー起動後、Swagger UIでAPI仕様を確認できます。

```
http://localhost:8090/swagger/index.html
```

### 主要エンドポイント

#### 同期
| メソッド | パス | 説明 |
|---------|------|------|
| POST | `/api/sync` | 全顧客の投稿を同期（20件並列処理） |
| POST | `/api/sync/{id}` | 特定顧客の投稿を同期 |

#### トークン管理
| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/token` | Instagram APIトークンを取得 |
| POST | `/api/token` | Instagram APIトークンを保存 |
| POST | `/api/token/check` | トークンの認証情報を確認 |

#### WordPress-Instagram連携管理
| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/wordpress-instagram` | 連携情報一覧取得（ページング、フィルタリング対応） |
| GET | `/api/wordpress-instagram/{id}` | 連携情報詳細取得 |
| POST | `/api/wordpress-instagram` | 連携情報作成 |
| PUT | `/api/wordpress-instagram/{id}` | 連携情報更新 |
| DELETE | `/api/wordpress-instagram/{id}` | 連携情報削除 |

## データベーススキーマ

### wordpress_instagrams テーブル

WordPress-Instagram連携情報を管理します。

| カラム名 | 型 | 説明 |
|---------|---|------|
| id | INT | 主キー |
| name | VARCHAR(255) | アカウント名 |
| wordpress | VARCHAR(255) | WordPress URL |
| instagram_id | VARCHAR(255) | Instagram ビジネスアカウントID |
| memo | TEXT | メモ |
| start_date | DATETIME | 連携開始日 |
| status | INT | ステータス（0=無効, 1=有効） |
| delete_hash | TINYINT | 削除フラグ |
| customer_type | INT | 顧客種別 |
| update_at | DATETIME | 更新日時 |
| create_at | DATETIME | 作成日時 |

### token テーブル

Instagram APIトークンを管理します。

| カラム名 | 型 | 説明 |
|---------|---|------|
| id | INT | 主キー |
| token | VARCHAR(500) | Instagram Graph APIトークン |
| update_at | DATETIME | 更新日時 |
| create_at | DATETIME | 作成日時 |

### posts テーブル

連携済み投稿を管理します。

| カラム名 | 型 | 説明 |
|---------|---|------|
| id | INT | 主キー |
| media_id | VARCHAR(45) | Instagram メディアID |
| customer_id | INT | 顧客ID（wordpress_instagrams.id + 100000） |
| timestamp | VARCHAR(45) | 投稿日時 |
| media_url | MEDIUMTEXT | メディアURL |
| permalink | VARCHAR(255) | Instagram パーマリンク |
| wordpress_link | VARCHAR(255) | WordPress 投稿URL |
| created_at | DATETIME | レコード作成日時 |

## 開発

### テストの実行

```bash
make test
```

### Swagger仕様の再生成

コードのコメントを修正した後、Swagger仕様を再生成します。

```bash
make swag
```

### コードスタイル

- gofmtでフォーマット
- golangci-lintでリント（推奨）

## 運用

### ログ

- 標準出力にJSON形式でログを出力
- Echoのミドルウェアでリクエスト/レスポンスをログ記録

### エラー通知

同期処理でエラーが発生した場合、Slackに通知されます。

通知内容:
- 顧客ID
- 顧客名
- エラーメッセージ

### 同期処理の仕様

1. **重複チェック**: 既に連携済みの投稿はスキップ
2. **日付フィルタ**: `start_date`以前の投稿はスキップ
3. **メディア処理**:
   - 画像/動画を一時ディレクトリにダウンロード
   - WordPressへアップロード
   - アップロード後、一時ディレクトリを削除
4. **トランザクション**: 投稿記録をDBに保存

## トラブルシューティング

### マイグレーションエラー

```bash
# 現在の状態確認
sql-migrate status

# ロールバック
sql-migrate down

# 再実行
sql-migrate up
```

### データベース接続エラー

- 環境変数が正しく設定されているか確認
- MySQLが起動しているか確認
- データベースが作成されているか確認

```bash
mysql -u $DB_USER -p$DB_PASSWORD -h $DB_HOST -P $DB_PORT -e "SHOW DATABASES;"
```

### Instagram API エラー

- Facebook トークンの有効期限を確認
- Instagram Graph APIの権限を確認
- レート制限に達していないか確認

## 🎓 学んだこと・工夫した点

### 設計面

1. **依存性の逆転原則 (DIP)**
   - Usecaseは具象ではなくインターフェースに依存
   - テスト時のモック化が容易で、外部APIに依存しないテスト実装を実現

2. **リポジトリパターン**
   - データアクセスロジックをカプセル化
   - 柔軟なクエリ条件をFilterパターンで実現

3. **Adapterパターン**
   - 各外部APIの差異を吸収し、統一されたインターフェースを提供
   - 将来的な外部サービス変更時の影響範囲を最小化

### 実装面

1. **エラーハンドリング**
   - カスタムエラー型を定義し、エラーの種類を判別可能に
   - Slackへの通知で、運用チームへリアルタイム通知

2. **データ整合性**
   - メディアIDによる冪等性チェックで重複投稿を防止
   - 連携開始日フィルターで、過去データの不要な連携を回避

3. **パフォーマンス**
   - **セマフォパターンによる20件並列処理**で大量アカウントの高速同期を実現
   - HTTPクライアントのタイムアウト設定
   - 一時ファイルの適切なクリーンアップ

### チーム開発を想定した工夫

1. **セルフドキュメンティング**
   - Swagger/OpenAPIによる自動生成API仕様
   - 日本語コメントによる可読性向上

2. **開発体験**
   - Airによるホットリロードで開発効率向上
   - Makefileでよく使うコマンドを簡略化

3. **テスト戦略**
   - Testcontainersで本番環境に近い統合テスト
   - リポジトリ層の網羅的なテストケース

## 🚀 今後の改善案

- [ ] テストカバレッジの向上
- [ ] ログ構造化とログレベル管理
- [ ] メトリクス監視（Prometheus対応）
- [ ] レート制限対策（指数バックオフ）
- [ ] Webhook対応（Instagram投稿時の即座連携）

## 📊 プロジェクト統計

- **開発期間**: [期間を記載]
- **総コード行数**: 4,500行以上（Go）
- **ファイル数**: 47ファイル
- **テストファイル数**: 7ファイル
- **外部API連携数**: 3サービス（Instagram Graph API、WordPress REST API、Slack Webhook）
- **エンドポイント数**: 10以上（同期、トークン管理、WordPress-Instagram連携CRUD）
- **マイグレーションファイル数**: 5ファイル
