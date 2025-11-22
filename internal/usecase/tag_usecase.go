package usecase

import (
	"context"
	"fmt"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/domain/repository"
)

// TagUseCase タグユースケースのインターフェース
type TagUseCase interface {
	Create(ctx context.Context, req *CreateTagRequest) (*entity.Tag, error)
	Update(ctx context.Context, id int64, req *UpdateTagRequest) (*entity.Tag, error)
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*entity.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Tag, error)
	List(ctx context.Context) ([]*entity.Tag, error)
}

// CreateTagRequest タグ作成リクエスト
type CreateTagRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// UpdateTagRequest タグ更新リクエスト
type UpdateTagRequest struct {
	Name *string `json:"name"`
	Slug *string `json:"slug"`
}

// tagUseCase TagUseCaseの実装
type tagUseCase struct {
	tagRepo repository.TagRepository
}

// NewTagUseCase 新しいTagUseCaseを作成
func NewTagUseCase(tagRepo repository.TagRepository) TagUseCase {
	return &tagUseCase{
		tagRepo: tagRepo,
	}
}

// Create 新しいタグを作成
func (u *tagUseCase) Create(ctx context.Context, req *CreateTagRequest) (*entity.Tag, error) {
	if req.Name == "" || req.Slug == "" {
		return nil, fmt.Errorf("name and slug are required")
	}

	tag := &entity.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := u.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

// Update タグを更新
func (u *tagUseCase) Update(ctx context.Context, id int64, req *UpdateTagRequest) (*entity.Tag, error) {
	tag, err := u.tagRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}

	if req.Name != nil {
		tag.Name = *req.Name
	}
	if req.Slug != nil {
		tag.Slug = *req.Slug
	}

	if err := u.tagRepo.Update(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return tag, nil
}

// Delete タグを削除
func (u *tagUseCase) Delete(ctx context.Context, id int64) error {
	if err := u.tagRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

// GetByID IDでタグを取得
func (u *tagUseCase) GetByID(ctx context.Context, id int64) (*entity.Tag, error) {
	tag, err := u.tagRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}
	return tag, nil
}

// GetBySlug スラッグでタグを取得
func (u *tagUseCase) GetBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	tag, err := u.tagRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}
	return tag, nil
}

// List タグ一覧を取得
func (u *tagUseCase) List(ctx context.Context) ([]*entity.Tag, error) {
	tags, err := u.tagRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	return tags, nil
}
