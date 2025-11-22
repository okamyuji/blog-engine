package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/domain/repository"

	"github.com/uptrace/bun"
)

// userRepositoryImpl UserRepositoryの実装
type userRepositoryImpl struct {
	db *bun.DB
}

// NewUserRepository 新しいUserRepositoryを作成
func NewUserRepository(db *bun.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create 新しいユーザーを作成
func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByID IDでユーザーを検索
func (r *userRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	user := new(entity.User)
	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByUsername ユーザー名でユーザーを検索
func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := new(entity.User)
	err := r.db.NewSelect().
		Model(user).
		Where("username = ?", username).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByEmail メールアドレスでユーザーを検索
func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := new(entity.User)
	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// Update ユーザー情報を更新
func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		OmitZero().
		Column("username", "email", "password_hash", "role", "status", "updated_at").
		WherePK().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete ユーザーを削除
func (r *userRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().
		Model((*entity.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List 全ユーザーを取得
func (r *userRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	users := make([]*entity.User, 0)
	err := r.db.NewSelect().
		Model(&users).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// Count ユーザー数を取得
func (r *userRepositoryImpl) Count(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*entity.User)(nil)).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
