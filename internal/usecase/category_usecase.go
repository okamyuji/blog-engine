package usecase

import (
	"context"
	"fmt"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/domain/repository"
)

// CategoryUseCase カテゴリユースケースのインターフェース
type CategoryUseCase interface {
	Create(ctx context.Context, req *CreateCategoryRequest) (*entity.Category, error)
	Update(ctx context.Context, id int64, req *UpdateCategoryRequest) (*entity.Category, error)
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	List(ctx context.Context) ([]*entity.Category, error)
}

// CreateCategoryRequest カテゴリ作成リクエスト
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// UpdateCategoryRequest カテゴリ更新リクエスト
type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
}

// categoryUseCase CategoryUseCaseの実装
type categoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryUseCase 新しいCategoryUseCaseを作成
func NewCategoryUseCase(categoryRepo repository.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{
		categoryRepo: categoryRepo,
	}
}

// Create 新しいカテゴリを作成
func (u *categoryUseCase) Create(ctx context.Context, req *CreateCategoryRequest) (*entity.Category, error) {
	if req.Name == "" || req.Slug == "" {
		return nil, fmt.Errorf("name and slug are required")
	}

	category := &entity.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := u.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// Update カテゴリを更新
func (u *categoryUseCase) Update(ctx context.Context, id int64, req *UpdateCategoryRequest) (*entity.Category, error) {
	category, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.Description != nil {
		category.Description = *req.Description
	}

	if err := u.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// Delete カテゴリを削除
func (u *categoryUseCase) Delete(ctx context.Context, id int64) error {
	if err := u.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

// GetByID IDでカテゴリを取得
func (u *categoryUseCase) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	category, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	return category, nil
}

// GetBySlug スラッグでカテゴリを取得
func (u *categoryUseCase) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	category, err := u.categoryRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	return category, nil
}

// List カテゴリ一覧を取得
func (u *categoryUseCase) List(ctx context.Context) ([]*entity.Category, error) {
	categories, err := u.categoryRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	return categories, nil
}
