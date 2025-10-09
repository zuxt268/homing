# ビルドステージ
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# バイナリのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/homing ./cmd/homing

# 実行ステージ
FROM alpine:latest

WORKDIR /root/

# 証明書のインストール（外部APIへのHTTPS通信用）とヘルスチェック用ツール
RUN apk --no-cache add ca-certificates wget

# ビルドステージからバイナリをコピー
COPY --from=builder /app/homing .

# ポートの公開
EXPOSE 8090

# アプリケーションの実行
CMD ["./homing"]