package persistence_test

import (
	"context"
	"testing"
	"time"

	"my-blog-engine/internal/infrastructure/persistence"
	"my-blog-engine/tests/integration/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenRepository_Add(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTokenRepository(db)
	ctx := context.Background()

	jti := "test-jti-123"
	expiresAt := time.Now().Add(1 * time.Hour)

	err := repo.Add(ctx, jti, expiresAt)
	require.NoError(t, err)
}

func TestTokenRepository_Exists(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTokenRepository(db)
	ctx := context.Background()

	jti := "test-jti-123"
	expiresAt := time.Now().Add(1 * time.Hour)

	// 追加
	err := repo.Add(ctx, jti, expiresAt)
	require.NoError(t, err)

	// 存在確認
	exists, err := repo.Exists(ctx, jti)
	require.NoError(t, err)
	assert.True(t, exists)

	// 存在しないJTI
	exists, err = repo.Exists(ctx, "nonexistent")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestTokenRepository_CleanupExpired(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTokenRepository(db)
	ctx := context.Background()

	// 期限切れトークン追加
	expiredJTI := "expired-jti"
	err := repo.Add(ctx, expiredJTI, time.Now().Add(-1*time.Hour))
	require.NoError(t, err)

	// 有効なトークン追加
	validJTI := "valid-jti"
	err = repo.Add(ctx, validJTI, time.Now().Add(1*time.Hour))
	require.NoError(t, err)

	// クリーンアップ実行
	err = repo.CleanupExpired(ctx)
	require.NoError(t, err)

	// 期限切れトークンが削除されたか確認
	exists, err := repo.Exists(ctx, expiredJTI)
	require.NoError(t, err)
	assert.False(t, exists)

	// 有効なトークンは残っているか確認
	exists, err = repo.Exists(ctx, validJTI)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestTokenRepository_FindByJTI(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTokenRepository(db)
	ctx := context.Background()

	jti := "test-jti-123"
	expiresAt := time.Now().UTC().Add(1 * time.Hour).Truncate(time.Microsecond)

	err := repo.Add(ctx, jti, expiresAt)
	require.NoError(t, err)

	token, err := repo.FindByJTI(ctx, jti)
	require.NoError(t, err)
	assert.Equal(t, jti, token.TokenJTI)
	// UTC時刻で比較（マイクロ秒精度）
	assert.True(t, expiresAt.Equal(token.ExpiresAt.UTC().Truncate(time.Microsecond)),
		"expected %v, got %v", expiresAt, token.ExpiresAt)
}

func TestTokenRepository_FindByJTI_NotFound(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTokenRepository(db)
	ctx := context.Background()

	_, err := repo.FindByJTI(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
