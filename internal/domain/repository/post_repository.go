package repository

import (
	"context"
	"my-blog-engine/internal/domain/entity"
)

// PostRepository 記事リポジトリのインターフェース
type PostRepository interface {
	// Create 新しい記事を作成
	Create(ctx context.Context, post *entity.Post) error

	// FindByID IDで記事を検索
	FindByID(ctx context.Context, id int64) (*entity.Post, error)

	// FindBySlug スラッグで記事を検索
	FindBySlug(ctx context.Context, slug string) (*entity.Post, error)

	// Update 記事を更新
	Update(ctx context.Context, post *entity.Post) error

	// Delete 記事を削除
	Delete(ctx context.Context, id int64) error

	// List 記事一覧を取得
	List(ctx context.Context, limit, offset int) ([]*entity.Post, error)

	// ListPublished 公開済み記事一覧を取得
	ListPublished(ctx context.Context, limit, offset int) ([]*entity.Post, error)

	// ListByCategory カテゴリ別記事一覧を取得
	ListByCategory(ctx context.Context, categoryID int64, limit, offset int) ([]*entity.Post, error)

	// ListByTag タグ別記事一覧を取得
	ListByTag(ctx context.Context, tagID int64, limit, offset int) ([]*entity.Post, error)

	// ListByAuthor 著者別記事一覧を取得
	ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*entity.Post, error)

	// Count 記事数を取得
	Count(ctx context.Context) (int, error)

	// CountPublished 公開済み記事数を取得
	CountPublished(ctx context.Context) (int, error)

	// AddTags 記事にタグを追加
	AddTags(ctx context.Context, postID int64, tagIDs []int64) error

	// RemoveTags 記事からタグを削除
	RemoveTags(ctx context.Context, postID int64, tagIDs []int64) error

	// GetTags 記事のタグを取得
	GetTags(ctx context.Context, postID int64) ([]*entity.Tag, error)
}
