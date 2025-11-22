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

// categoryRepositoryImpl CategoryRepositoryの実装
type categoryRepositoryImpl struct {
	db *bun.DB
}

// NewCategoryRepository 新しいCategoryRepositoryを作成
func NewCategoryRepository(db *bun.DB) repository.CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

// Create 新しいカテゴリを作成
func (r *categoryRepositoryImpl) Create(ctx context.Context, category *entity.Category) error {
	_, err := r.db.NewInsert().
		Model(category).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

// FindByID IDでカテゴリを検索
func (r *categoryRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.Category, error) {
	category := new(entity.Category)
	err := r.db.NewSelect().
		Model(category).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	return category, nil
}

// FindBySlug スラッグでカテゴリを検索
func (r *categoryRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	category := new(entity.Category)
	err := r.db.NewSelect().
		Model(category).
		Where("slug = ?", slug).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	return category, nil
}

// Update カテゴリを更新
func (r *categoryRepositoryImpl) Update(ctx context.Context, category *entity.Category) error {
	_, err := r.db.NewUpdate().
		Model(category).
		OmitZero().
		Column("name", "slug", "description", "updated_at").
		WherePK().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	return nil
}

// Delete カテゴリを削除
func (r *categoryRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().
		Model((*entity.Category)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// List カテゴリ一覧を取得
func (r *categoryRepositoryImpl) List(ctx context.Context) ([]*entity.Category, error) {
	categories := make([]*entity.Category, 0)
	err := r.db.NewSelect().
		Model(&categories).
		Order("name ASC").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

// Count カテゴリ数を取得
func (r *categoryRepositoryImpl) Count(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*entity.Category)(nil)).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count categories: %w", err)
	}

	return count, nil
}
