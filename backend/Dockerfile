# マルチステージビルド
FROM golang:1.24-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /app

# Go modulesファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# マイグレーションツールをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate/main.go

# シードツールをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed ./cmd/seed/main.go

# 本番用の軽量イメージ
FROM alpine:latest AS production

# 必要なパッケージをインストール
RUN apk --no-cache add ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/seed .

# ポートを公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./main"]

# 開発用イメージ
FROM golang:1.24-alpine AS development

# 必要なパッケージをインストール
RUN apk add --no-cache git ca-certificates tzdata make curl

# Air（ホットリロード）をインストール - Go 1.24互換バージョン
RUN go install github.com/air-verse/air@v1.52.3

# golangci-lintをインストール
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.0

# 作業ディレクトリを設定
WORKDIR /app

# Go modulesファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー（開発時はボリュームマウントでオーバーライド）
COPY . .

# ポートを公開
EXPOSE 8080

# ヘルスチェック用のスクリプト
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# デフォルトコマンド（docker-compose.ymlでオーバーライド）
CMD ["go", "run", "main.go"]