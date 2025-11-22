package usecase_test

import (
	"context"
	"testing"
	"time"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/infrastructure/auth"
	"my-blog-engine/internal/infrastructure/persistence"
	"my-blog-engine/internal/usecase"
	"my-blog-engine/tests/integration/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthUseCase(t *testing.T) (usecase.AuthUseCase, func()) {
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

	return authUseCase, cleanup
}

func TestAuthUseCase_Login(t *testing.T) {
	ctx := context.Background()

	// テストユーザー作成
	db, dbCleanup := testhelper.SetupTestDB(t)
	defer dbCleanup()

	userRepo := persistence.NewUserRepository(db)
	passwordHasher := auth.NewPasswordHasher()

	hash, err := passwordHasher.Hash("password123")
	require.NoError(t, err)

	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hash,
		Role:         entity.RoleAdmin,
		Status:       entity.StatusActive,
	}
	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	// 再度UseCaseを初期化（同じDBを使用）
	tokenRepo := persistence.NewTokenRepository(db)
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

	// ログインテスト
	response, err := authUseCase.Login(ctx, "testuser", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Equal(t, user.Username, response.User.Username)
	assert.Empty(t, response.User.PasswordHash) // パスワードハッシュ除外確認
}

func TestAuthUseCase_Login_InvalidCredentials(t *testing.T) {
	authUseCase, cleanup := setupAuthUseCase(t)
	defer cleanup()

	ctx := context.Background()

	// 存在しないユーザーでログイン
	_, err := authUseCase.Login(ctx, "nonexistent", "password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthUseCase_Logout(t *testing.T) {
	ctx := context.Background()

	// テストユーザー作成とログイン
	db, dbCleanup := testhelper.SetupTestDB(t)
	defer dbCleanup()

	userRepo := persistence.NewUserRepository(db)
	tokenRepo := persistence.NewTokenRepository(db)
	passwordHasher := auth.NewPasswordHasher()

	hash, err := passwordHasher.Hash("password123")
	require.NoError(t, err)

	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hash,
		Role:         entity.RoleAdmin,
		Status:       entity.StatusActive,
	}
	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

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

	response, err := authUseCase.Login(ctx, "testuser", "password123")
	require.NoError(t, err)

	// ログアウト
	err = authUseCase.Logout(ctx, response.AccessToken)
	assert.NoError(t, err)

	// トークン検証（ブラックリスト登録済みのため失敗するはず）
	_, err = authUseCase.ValidateToken(ctx, response.AccessToken)
	assert.Error(t, err)
}

func TestAuthUseCase_ValidateToken(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	userRepo := persistence.NewUserRepository(db)
	tokenRepo := persistence.NewTokenRepository(db)
	passwordHasher := auth.NewPasswordHasher()

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

	response, err := authUseCase.Login(ctx, "testuser", "password123")
	require.NoError(t, err)

	// トークン検証
	validatedUser, err := authUseCase.ValidateToken(ctx, response.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, user.Username, validatedUser.Username)
	assert.Equal(t, user.Role, validatedUser.Role)
	assert.Empty(t, validatedUser.PasswordHash)
}
