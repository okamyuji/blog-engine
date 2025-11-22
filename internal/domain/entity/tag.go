package entity

import (
	"time"

	"github.com/uptrace/bun"
)

// Tag タグエンティティ
type Tag struct {
	bun.BaseModel `bun:"table:tags,alias:t"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name,notnull"`
	Slug      string    `bun:"slug,unique,notnull"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`

	// Relations
	Posts []*Post `bun:"m2m:post_tags,join:Tag=Post"`
}
