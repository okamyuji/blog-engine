package repository

import (
	"context"
	"my-blog-engine/internal/domain/entity"
)

// CategoryRepository カテゴリリポジトリのインターフェース
type CategoryRepository interface {
	// Create 新しいカテゴリを作成
	Create(ctx context.Context, category *entity.Category) error

	// FindByID IDでカテゴリを検索
	FindByID(ctx context.Context, id int64) (*entity.Category, error)

	// FindBySlug スラッグでカテゴリを検索
	FindBySlug(ctx context.Context, slug string) (*entity.Category, error)

	// Update カテゴリを更新
	Update(ctx context.Context, category *entity.Category) error

	// Delete カテゴリを削除
	Delete(ctx context.Context, id int64) error

	// List カテゴリ一覧を取得
	List(ctx context.Context) ([]*entity.Category, error)

	// Count カテゴリ数を取得
	Count(ctx context.Context) (int, error)
}
