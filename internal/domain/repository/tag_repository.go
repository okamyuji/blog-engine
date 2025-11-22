package repository

import (
	"context"
	"my-blog-engine/internal/domain/entity"
)

// TagRepository タグリポジトリのインターフェース
type TagRepository interface {
	// Create 新しいタグを作成
	Create(ctx context.Context, tag *entity.Tag) error

	// FindByID IDでタグを検索
	FindByID(ctx context.Context, id int64) (*entity.Tag, error)

	// FindBySlug スラッグでタグを検索
	FindBySlug(ctx context.Context, slug string) (*entity.Tag, error)

	// FindByIDs 複数のIDでタグを検索
	FindByIDs(ctx context.Context, ids []int64) ([]*entity.Tag, error)

	// Update タグを更新
	Update(ctx context.Context, tag *entity.Tag) error

	// Delete タグを削除
	Delete(ctx context.Context, id int64) error

	// List タグ一覧を取得
	List(ctx context.Context) ([]*entity.Tag, error)

	// Count タグ数を取得
	Count(ctx context.Context) (int, error)
}
