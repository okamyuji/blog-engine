package repository

import (
	"context"
	"my-blog-engine/internal/domain/entity"
)

// UserRepository ユーザーリポジトリのインターフェース
type UserRepository interface {
	// Create 新しいユーザーを作成
	Create(ctx context.Context, user *entity.User) error

	// FindByID IDでユーザーを検索
	FindByID(ctx context.Context, id int64) (*entity.User, error)

	// FindByUsername ユーザー名でユーザーを検索
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// FindByEmail メールアドレスでユーザーを検索
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update ユーザー情報を更新
	Update(ctx context.Context, user *entity.User) error

	// Delete ユーザーを削除
	Delete(ctx context.Context, id int64) error

	// List 全ユーザーを取得
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)

	// Count ユーザー数を取得
	Count(ctx context.Context) (int, error)
}
