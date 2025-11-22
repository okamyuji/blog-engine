package entity

import "github.com/uptrace/bun"

// PostTag 記事とタグの多対多関係を表す中間テーブル
type PostTag struct {
	bun.BaseModel `bun:"table:post_tags,alias:pt"`

	PostID int64 `bun:"post_id,pk,notnull"`
	TagID  int64 `bun:"tag_id,pk,notnull"`

	// Relations
	Post *Post `bun:"rel:belongs-to,join:post_id=id"`
	Tag  *Tag  `bun:"rel:belongs-to,join:tag_id=id"`
}
