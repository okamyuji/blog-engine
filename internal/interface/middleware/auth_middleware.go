package middleware

import (
	"context"
	"net/http"
	"strings"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/usecase"
)

// contextKey コンテキストキーの型
type contextKey string

const (
	// UserContextKeyユーザーコンテキストキー
	UserContextKey contextKey = "user"
)

// AuthMiddleware 認証ミドルウェア
type AuthMiddleware struct {
	authUseCase usecase.AuthUseCase
}

// NewAuthMiddleware 新しいAuthMiddlewareを作成
func NewAuthMiddleware(authUseCase usecase.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
	}
}

// Authenticate 認証を行うミドルウェア
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorizationヘッダーを取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Bearer トークンを抽出
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// トークン検証
		user, err := m.authUseCase.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// ユーザーをコンテキストに追加
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole 指定されたロールを必要とするミドルウェア
func (m *AuthMiddleware) RequireRole(roles ...entity.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserContextKey).(*entity.User)
			if !ok {
				http.Error(w, "User not found in context", http.StatusUnauthorized)
				return
			}

			// ロールチェック
			hasRole := false
			for _, role := range roles {
				if user.HasRole(role) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext コンテキストからユーザーを取得
func GetUserFromContext(ctx context.Context) (*entity.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*entity.User)
	return user, ok
}
