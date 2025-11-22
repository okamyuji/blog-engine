package repository

import (
	"context"
	"my-blog-engine/internal/domain/entity"
	"time"
)

// TokenRepository トークンブラックリストリポジトリのインターフェース
type TokenRepository interface {
	// Add トークンをブラックリストに追加
	Add(ctx context.Context, jti string, expiresAt time.Time) error

	// Exists トークンがブラックリストに存在するかチェック
	Exists(ctx context.Context, jti string) (bool, error)

	// CleanupExpired 期限切れトークンを削除
	CleanupExpired(ctx context.Context) error

	// FindByJTI JTIでトークンを検索
	FindByJTI(ctx context.Context, jti string) (*entity.TokenBlacklist, error)
}
