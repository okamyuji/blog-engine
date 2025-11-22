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

// tagRepositoryImpl TagRepositoryの実装
type tagRepositoryImpl struct {
	db *bun.DB
}

// NewTagRepository 新しいTagRepositoryを作成
func NewTagRepository(db *bun.DB) repository.TagRepository {
	return &tagRepositoryImpl{db: db}
}

// Create 新しいタグを作成
func (r *tagRepositoryImpl) Create(ctx context.Context, tag *entity.Tag) error {
	_, err := r.db.NewInsert().
		Model(tag).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

// FindByID IDでタグを検索
func (r *tagRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.Tag, error) {
	tag := new(entity.Tag)
	err := r.db.NewSelect().
		Model(tag).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("tag not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}

	return tag, nil
}

// FindBySlug スラッグでタグを検索
func (r *tagRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	tag := new(entity.Tag)
	err := r.db.NewSelect().
		Model(tag).
		Where("slug = ?", slug).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("tag not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}

	return tag, nil
}

// FindByIDs 複数のIDでタグを検索
func (r *tagRepositoryImpl) FindByIDs(ctx context.Context, ids []int64) ([]*entity.Tag, error) {
	if len(ids) == 0 {
		return []*entity.Tag{}, nil
	}

	tags := make([]*entity.Tag, 0)
	err := r.db.NewSelect().
		Model(&tags).
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to find tags: %w", err)
	}

	return tags, nil
}

// Update タグを更新
func (r *tagRepositoryImpl) Update(ctx context.Context, tag *entity.Tag) error {
	_, err := r.db.NewUpdate().
		Model(tag).
		OmitZero().
		Column("name", "slug", "updated_at").
		WherePK().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	return nil
}

// Delete タグを削除
func (r *tagRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().
		Model((*entity.Tag)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// List タグ一覧を取得
func (r *tagRepositoryImpl) List(ctx context.Context) ([]*entity.Tag, error) {
	tags := make([]*entity.Tag, 0)
	err := r.db.NewSelect().
		Model(&tags).
		Order("name ASC").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return tags, nil
}

// Count タグ数を取得
func (r *tagRepositoryImpl) Count(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*entity.Tag)(nil)).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count tags: %w", err)
	}

	return count, nil
}
