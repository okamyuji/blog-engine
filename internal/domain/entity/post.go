package entity

import (
	"time"

	"github.com/uptrace/bun"
)

// PostStatus 投稿のステータスを表す型
type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPublished PostStatus = "published"
)

// Post ブログ記事エンティティ
type Post struct {
	bun.BaseModel `bun:"table:posts,alias:p"`

	ID           int64      `bun:"id,pk,autoincrement"`
	Title        string     `bun:"title,notnull"`
	Slug         string     `bun:"slug,unique,notnull"`
	Content      string     `bun:"content,notnull,type:text"`
	RenderedHTML string     `bun:"rendered_html,type:text"`
	Status       PostStatus `bun:"status,notnull,default:'draft'"`
	AuthorID     int64      `bun:"author_id,notnull"`
	CategoryID   *int64     `bun:"category_id"`
	CreatedAt    time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	PublishedAt  *time.Time `bun:"published_at"`

	// Relations
	Author   *User     `bun:"rel:belongs-to,join:author_id=id"`
	Category *Category `bun:"rel:belongs-to,join:category_id=id"`
	Tags     []*Tag    `bun:"m2m:post_tags,join:Post=Tag"`
}

// IsPublished 公開済みかどうかを判定
func (p *Post) IsPublished() bool {
	return p.Status == StatusPublished
}

// Publish 記事を公開する
func (p *Post) Publish() {
	now := time.Now()
	p.Status = StatusPublished
	if p.PublishedAt == nil {
		p.PublishedAt = &now
	}
}

// Unpublish 記事を下書きに戻す
func (p *Post) Unpublish() {
	p.Status = StatusDraft
}
