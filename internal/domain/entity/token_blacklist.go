package entity

import (
	"time"

	"github.com/uptrace/bun"
)

// TokenBlacklist JWTブラックリストエンティティ
type TokenBlacklist struct {
	bun.BaseModel `bun:"table:token_blacklist,alias:tb"`

	ID        int64     `bun:"id,pk,autoincrement"`
	TokenJTI  string    `bun:"token_jti,unique,notnull"`
	ExpiresAt time.Time `bun:"expires_at,notnull"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

// IsExpired トークンが期限切れかどうかを判定
func (t *TokenBlacklist) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
