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

// postRepositoryImpl PostRepositoryの実装
type postRepositoryImpl struct {
	db *bun.DB
}

// NewPostRepository 新しいPostRepositoryを作成
func NewPostRepository(db *bun.DB) repository.PostRepository {
	return &postRepositoryImpl{db: db}
}

// Create 新しい記事を作成
func (r *postRepositoryImpl) Create(ctx context.Context, post *entity.Post) error {
	_, err := r.db.NewInsert().
		Model(post).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// FindByID IDで記事を検索
func (r *postRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.Post, error) {
	post := new(entity.Post)
	err := r.db.NewSelect().
		Model(post).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Where("p.id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find post: %w", err)
	}

	return post, nil
}

// FindBySlug スラッグで記事を検索
func (r *postRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entity.Post, error) {
	post := new(entity.Post)
	err := r.db.NewSelect().
		Model(post).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Where("p.slug = ?", slug).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find post: %w", err)
	}

	return post, nil
}

// Update 記事を更新
func (r *postRepositoryImpl) Update(ctx context.Context, post *entity.Post) error {
	_, err := r.db.NewUpdate().
		Model(post).
		OmitZero().
		Column("title", "slug", "content", "category_id", "author_id", "status", "published_at", "updated_at").
		WherePK().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

// Delete 記事を削除
func (r *postRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().
		Model((*entity.Post)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// List 記事一覧を取得
func (r *postRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Post, error) {
	posts := make([]*entity.Post, 0)
	err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, nil
}

// ListPublished 公開済み記事一覧を取得
func (r *postRepositoryImpl) ListPublished(ctx context.Context, limit, offset int) ([]*entity.Post, error) {
	posts := make([]*entity.Post, 0)
	err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Where("p.status = ?", entity.StatusPublished).
		Order("published_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list published posts: %w", err)
	}

	return posts, nil
}

// ListByCategory カテゴリ別記事一覧を取得
func (r *postRepositoryImpl) ListByCategory(ctx context.Context, categoryID int64, limit, offset int) ([]*entity.Post, error) {
	posts := make([]*entity.Post, 0)
	err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Where("p.category_id = ?", categoryID).
		Where("p.status = ?", entity.StatusPublished).
		Order("published_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list posts by category: %w", err)
	}

	return posts, nil
}

// ListByTag タグ別記事一覧を取得
func (r *postRepositoryImpl) ListByTag(ctx context.Context, tagID int64, limit, offset int) ([]*entity.Post, error) {
	posts := make([]*entity.Post, 0)
	err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Join("JOIN post_tags AS pt ON pt.post_id = p.id").
		Where("pt.tag_id = ?", tagID).
		Where("p.status = ?", entity.StatusPublished).
		Order("published_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list posts by tag: %w", err)
	}

	return posts, nil
}

// ListByAuthor 著者別記事一覧を取得
func (r *postRepositoryImpl) ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*entity.Post, error) {
	posts := make([]*entity.Post, 0)
	err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Where("p.author_id = ?", authorID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list posts by author: %w", err)
	}

	return posts, nil
}

// Count 記事数を取得
func (r *postRepositoryImpl) Count(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*entity.Post)(nil)).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %w", err)
	}

	return count, nil
}

// CountPublished 公開済み記事数を取得
func (r *postRepositoryImpl) CountPublished(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*entity.Post)(nil)).
		Where("status = ?", entity.StatusPublished).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count published posts: %w", err)
	}

	return count, nil
}

// AddTags 記事にタグを追加
func (r *postRepositoryImpl) AddTags(ctx context.Context, postID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}

	postTags := make([]*entity.PostTag, len(tagIDs))
	for i, tagID := range tagIDs {
		postTags[i] = &entity.PostTag{
			PostID: postID,
			TagID:  tagID,
		}
	}

	_, err := r.db.NewInsert().
		Model(&postTags).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to add tags: %w", err)
	}

	return nil
}

// RemoveTags 記事からタグを削除
func (r *postRepositoryImpl) RemoveTags(ctx context.Context, postID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}

	_, err := r.db.NewDelete().
		Model((*entity.PostTag)(nil)).
		Where("post_id = ?", postID).
		Where("tag_id IN (?)", bun.In(tagIDs)).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to remove tags: %w", err)
	}

	return nil
}

// GetTags 記事のタグを取得
func (r *postRepositoryImpl) GetTags(ctx context.Context, postID int64) ([]*entity.Tag, error) {
	tags := make([]*entity.Tag, 0)
	err := r.db.NewSelect().
		Model(&tags).
		Join("JOIN post_tags AS pt ON pt.tag_id = t.id").
		Where("pt.post_id = ?", postID).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, nil
}
