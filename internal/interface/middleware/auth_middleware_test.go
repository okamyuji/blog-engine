package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/infrastructure/auth"
	"my-blog-engine/internal/infrastructure/persistence"
	"my-blog-engine/internal/interface/middleware"
	"my-blog-engine/internal/usecase"
	"my-blog-engine/tests/integration/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthMiddleware(t *testing.T) (*middleware.AuthMiddleware, *entity.User, string, func()) {
	t.Helper()

	db, cleanup := testhelper.SetupTestDB(t)

	userRepo := persistence.NewUserRepository(db)
	tokenRepo := persistence.NewTokenRepository(db)

	passwordHasher := auth.NewPasswordHasher()
	jwtManager, err := auth.NewJWTManager(auth.JWTConfig{
		SecretKey:     "test-secret-key-min-32-chars-long",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 168 * time.Hour,
	})
	require.NoError(t, err)

	authUseCase := usecase.NewAuthUseCase(
		userRepo,
		tokenRepo,
		jwtManager,
		passwordHasher,
		15*time.Minute,
	)

	authMiddleware := middleware.NewAuthMiddleware(authUseCase)

	// テストユーザー作成
	ctx := context.Background()
	hash, err := passwordHasher.Hash("password123")
	require.NoError(t, err)

	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hash,
		Role:         entity.RoleEditor,
		Status:       entity.StatusActive,
	}
	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	// ログインしてトークン取得
	response, err := authUseCase.Login(ctx, "testuser", "password123")
	require.NoError(t, err)

	return authMiddleware, user, response.AccessToken, cleanup
}

func TestAuthMiddleware_Authenticate_Success(t *testing.T) {
	authMiddleware, user, token, cleanup := setupAuthMiddleware(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	handler := authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextUser, ok := middleware.GetUserFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, user.Username, contextUser.Username)
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_Authenticate_MissingHeader(t *testing.T) {
	authMiddleware, _, _, cleanup := setupAuthMiddleware(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler := authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not reach here")
	}))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_Authenticate_InvalidFormat(t *testing.T) {
	authMiddleware, _, _, cleanup := setupAuthMiddleware(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")

	rec := httptest.NewRecorder()

	handler := authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not reach here")
	}))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_RequireRole_Success(t *testing.T) {
	authMiddleware, user, token, cleanup := setupAuthMiddleware(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	handler := authMiddleware.Authenticate(
		authMiddleware.RequireRole(entity.RoleEditor)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				contextUser, ok := middleware.GetUserFromContext(r.Context())
				assert.True(t, ok)
				assert.Equal(t, user.Username, contextUser.Username)
				w.WriteHeader(http.StatusOK)
			}),
		),
	)

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_RequireRole_Forbidden(t *testing.T) {
	authMiddleware, _, token, cleanup := setupAuthMiddleware(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	// Editorユーザーでadminロールが必要なハンドラーにアクセス
	handler := authMiddleware.Authenticate(
		authMiddleware.RequireRole(entity.RoleAdmin)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("Should not reach here")
			}),
		),
	)

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}
