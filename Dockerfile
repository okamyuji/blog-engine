# ビルドステージ
FROM golang:1.25-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git make

WORKDIR /app

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /blog-engine ./cmd/blog

# 実行ステージ
FROM alpine:latest

# Mermaidレンダリング用にNode.js、npm、Chromiumをインストール
RUN apk add --no-cache \
    ca-certificates \
    nodejs \
    npm \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont

# Mermaid CLIをインストール
RUN npm install -g @mermaid-js/mermaid-cli

# Puppeteer用の環境変数設定
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

WORKDIR /app

# ビルド成果物をコピー
COPY --from=builder /blog-engine /app/blog-engine

# テンプレートと静的ファイルをコピー
COPY templates ./templates
COPY static ./static

# 非rootユーザーで実行
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["/app/blog-engine"]

