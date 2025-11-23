package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/infrastructure/auth"
	"my-blog-engine/internal/infrastructure/database"
	"my-blog-engine/internal/infrastructure/persistence"
	"my-blog-engine/internal/infrastructure/renderer"
	"my-blog-engine/internal/interface/handler"
	"my-blog-engine/internal/interface/middleware"
	"my-blog-engine/internal/usecase"
)

func main() {
	// ロガー設定
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting blog engine server...")

	// 設定読み込み
	cfg := loadConfig()

	// データベース接続
	db, err := database.NewConnection(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			slog.Error("Failed to close database", "error", err)
		}
	}()

	slog.Info("Database connected successfully")

	// Repository初期化
	userRepo := persistence.NewUserRepository(db)
	postRepo := persistence.NewPostRepository(db)
	categoryRepo := persistence.NewCategoryRepository(db)
	tagRepo := persistence.NewTagRepository(db)
	tokenRepo := persistence.NewTokenRepository(db)

	// Infrastructure初期化
	passwordHasher := auth.NewPasswordHasher()
	jwtManager, err := auth.NewJWTManager(auth.JWTConfig{
		SecretKey:     cfg.JWTSecret,
		AccessExpiry:  cfg.JWTAccessExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
	})
	if err != nil {
		log.Fatal("Failed to create JWT manager:", err)
	}

	mermaidRenderer := renderer.NewMermaidRenderer()
	mdRenderer := renderer.NewMarkdownRenderer(mermaidRenderer)

	// UseCase初期化
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtManager, passwordHasher, cfg.JWTAccessExpiry)
	postUseCase := usecase.NewPostUseCase(postRepo, categoryRepo, tagRepo, mdRenderer)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)
	tagUseCase := usecase.NewTagUseCase(tagRepo)

	// Handler初期化
	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(authUseCase)
	postHandler := handler.NewPostHandler(postUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	tagHandler := handler.NewTagHandler(tagUseCase)
	publicHandler := handler.NewPublicHandler(postUseCase, categoryUseCase)

	// Middleware初期化
	authMiddleware := middleware.NewAuthMiddleware(authUseCase)
	rateLimiter := middleware.NewRateLimiter(100, 200)

	// ルーター設定
	mux := http.NewServeMux()

	// 公開HTMLページ
	mux.HandleFunc("/", publicHandler.Home)

	// 公開エンドポイント
	mux.HandleFunc("/health", healthHandler.Check)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/refresh", authHandler.RefreshToken)

	// 認証が必要なエンドポイント
	mux.Handle("/api/auth/logout", authMiddleware.Authenticate(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("/api/auth/me", authMiddleware.Authenticate(http.HandlerFunc(authHandler.Me)))

	// 記事管理エンドポイント(編集権限が必要)
	mux.Handle("/api/admin/posts",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole(entity.RoleAdmin, entity.RoleEditor)(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						postHandler.List(w, r)
					case http.MethodPost:
						postHandler.Create(w, r)
					case http.MethodPut:
						postHandler.Update(w, r)
					case http.MethodDelete:
						postHandler.Delete(w, r)
					default:
						http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)

	mux.Handle("/api/admin/posts/publish",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole(entity.RoleAdmin, entity.RoleEditor)(
				http.HandlerFunc(postHandler.Publish),
			),
		),
	)

	mux.Handle("/api/admin/posts/unpublish",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole(entity.RoleAdmin, entity.RoleEditor)(
				http.HandlerFunc(postHandler.Unpublish),
			),
		),
	)

	// 公開記事エンドポイント
	mux.HandleFunc("/api/posts", postHandler.ListPublished)
	mux.HandleFunc("/api/posts/id", postHandler.GetByID)
	mux.HandleFunc("/api/posts/slug", postHandler.GetBySlug)

	// カテゴリエンドポイント
	mux.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			categoryHandler.List(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/categories/id", categoryHandler.GetByID)
	mux.HandleFunc("/api/categories/slug", categoryHandler.GetBySlug)

	mux.Handle("/api/admin/categories",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole(entity.RoleAdmin, entity.RoleEditor)(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						categoryHandler.List(w, r)
					case http.MethodPost:
						categoryHandler.Create(w, r)
					case http.MethodPut:
						categoryHandler.Update(w, r)
					case http.MethodDelete:
						categoryHandler.Delete(w, r)
					default:
						http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)

	// タグエンドポイント
	mux.HandleFunc("/api/tags", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tagHandler.List(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/tags/id", tagHandler.GetByID)
	mux.HandleFunc("/api/tags/slug", tagHandler.GetBySlug)

	mux.Handle("/api/admin/tags",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole(entity.RoleAdmin, entity.RoleEditor)(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						tagHandler.List(w, r)
					case http.MethodPost:
						tagHandler.Create(w, r)
					case http.MethodPut:
						tagHandler.Update(w, r)
					case http.MethodDelete:
						tagHandler.Delete(w, r)
					default:
						http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)

	// ミドルウェアチェーン
	handler := middleware.Recovery(
		middleware.Logging(
			middleware.SecurityHeaders(
				rateLimiter.Limit(mux),
			),
		),
	)

	// サーバー設定
	server := &http.Server{
		Addr:         cfg.ServerHost + ":" + cfg.ServerPort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// graceful shutdown設定
	go func() {
		slog.Info("Server starting", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	slog.Info("Server exited")
}

// Config アプリケーション設定
type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	ServerHost       string
	ServerPort       string
}

// loadConfig 環境変数から設定を読み込む
func loadConfig() Config {
	return Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "3306"),
		DBUser:           getEnv("DB_USER", "bloguser"),
		DBPassword:       getEnv("DB_PASSWORD", "blogpass"),
		DBName:           getEnv("DB_NAME", "blogdb"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key-min-32-chars-long-change-this-in-production"),
		JWTAccessExpiry:  parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"), 15*time.Minute),
		JWTRefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"), 168*time.Hour),
		ServerHost:       getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
	}
}

// getEnv 環境変数を取得(デフォルト値付き)
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration 文字列をtime.Durationにパース
func parseDuration(s string, defaultValue time.Duration) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultValue
	}
	return d
}
