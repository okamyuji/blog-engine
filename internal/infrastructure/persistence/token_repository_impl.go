package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/domain/repository"

	"github.com/uptrace/bun"
)

// tokenRepositoryImpl TokenRepositoryの実装
type tokenRepositoryImpl struct {
	db *bun.DB
}

// NewTokenRepository 新しいTokenRepositoryを作成
func NewTokenRepository(db *bun.DB) repository.TokenRepository {
	return &tokenRepositoryImpl{db: db}
}

// Add トークンをブラックリストに追加
func (r *tokenRepositoryImpl) Add(ctx context.Context, jti string, expiresAt time.Time) error {
	token := &entity.TokenBlacklist{
		TokenJTI:  jti,
		ExpiresAt: expiresAt,
	}

	_, err := r.db.NewInsert().
		Model(token).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}

	return nil
}

// Exists トークンがブラックリストに存在するかチェック
func (r *tokenRepositoryImpl) Exists(ctx context.Context, jti string) (bool, error) {
	count, err := r.db.NewSelect().
		Model((*entity.TokenBlacklist)(nil)).
		Where("token_jti = ?", jti).
		Where("expires_at > ?", time.Now()).
		Count(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check token existence: %w", err)
	}

	return count > 0, nil
}

// CleanupExpired 期限切れトークンを削除
func (r *tokenRepositoryImpl) CleanupExpired(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*entity.TokenBlacklist)(nil)).
		Where("expires_at <= ?", time.Now()).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return nil
}

// FindByJTI JTIでトークンを検索
func (r *tokenRepositoryImpl) FindByJTI(ctx context.Context, jti string) (*entity.TokenBlacklist, error) {
	token := new(entity.TokenBlacklist)
	err := r.db.NewSelect().
		Model(token).
		Where("token_jti = ?", jti).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("token not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find token: %w", err)
	}

	return token, nil
}
