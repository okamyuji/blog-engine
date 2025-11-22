package entity

import (
	"time"

	"github.com/uptrace/bun"
)

// Category カテゴリエンティティ
type Category struct {
	bun.BaseModel `bun:"table:categories,alias:c"`

	ID          int64     `bun:"id,pk,autoincrement"`
	Name        string    `bun:"name,notnull"`
	Slug        string    `bun:"slug,unique,notnull"`
	Description string    `bun:"description,type:text"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`

	// Relations
	Posts []*Post `bun:"rel:has-many,join:id=category_id"`
}
