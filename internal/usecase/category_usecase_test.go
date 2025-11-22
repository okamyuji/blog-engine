package usecase_test

import (
	"context"
	"testing"

	"my-blog-engine/internal/infrastructure/persistence"
	"my-blog-engine/internal/usecase"
	"my-blog-engine/tests/integration/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryUseCase_Create(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	categoryRepo := persistence.NewCategoryRepository(db)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)

	ctx := context.Background()

	req := &usecase.CreateCategoryRequest{
		Name:        "Tech",
		Slug:        "tech",
		Description: "Technology articles",
	}

	category, err := categoryUseCase.Create(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, req.Name, category.Name)
	assert.Equal(t, req.Slug, category.Slug)
	assert.Equal(t, req.Description, category.Description)
}

func TestCategoryUseCase_Update(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	categoryRepo := persistence.NewCategoryRepository(db)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)

	ctx := context.Background()

	// カテゴリ作成
	createReq := &usecase.CreateCategoryRequest{
		Name:        "Original",
		Slug:        "original",
		Description: "Original description",
	}
	category, err := categoryUseCase.Create(ctx, createReq)
	require.NoError(t, err)

	// 更新
	newName := "Updated"
	updateReq := &usecase.UpdateCategoryRequest{
		Name: &newName,
	}

	updated, err := categoryUseCase.Update(ctx, category.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, category.Slug, updated.Slug)
}

func TestCategoryUseCase_Delete(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	categoryRepo := persistence.NewCategoryRepository(db)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)

	ctx := context.Background()

	// カテゴリ作成
	createReq := &usecase.CreateCategoryRequest{
		Name: "ToDelete",
		Slug: "to-delete",
	}
	category, err := categoryUseCase.Create(ctx, createReq)
	require.NoError(t, err)

	// 削除
	err = categoryUseCase.Delete(ctx, category.ID)
	assert.NoError(t, err)

	// 取得失敗確認
	_, err = categoryUseCase.GetByID(ctx, category.ID)
	assert.Error(t, err)
}

func TestCategoryUseCase_GetBySlug(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	categoryRepo := persistence.NewCategoryRepository(db)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)

	ctx := context.Background()

	// カテゴリ作成
	createReq := &usecase.CreateCategoryRequest{
		Name: "Test",
		Slug: "test-slug",
	}
	created, err := categoryUseCase.Create(ctx, createReq)
	require.NoError(t, err)

	// スラッグで取得
	category, err := categoryUseCase.GetBySlug(ctx, "test-slug")
	require.NoError(t, err)
	assert.Equal(t, created.ID, category.ID)
	assert.Equal(t, created.Name, category.Name)
}

func TestCategoryUseCase_List(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	categoryRepo := persistence.NewCategoryRepository(db)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)

	ctx := context.Background()

	// カテゴリ作成
	categories := []string{"Cat1", "Cat2", "Cat3"}
	for _, name := range categories {
		req := &usecase.CreateCategoryRequest{
			Name: name,
			Slug: name,
		}
		_, err := categoryUseCase.Create(ctx, req)
		require.NoError(t, err)
	}

	// 一覧取得
	list, err := categoryUseCase.List(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 3)
}
