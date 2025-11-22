package persistence_test

import (
	"context"
	"fmt"
	"testing"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/infrastructure/persistence"
	"my-blog-engine/tests/integration/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_Create(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewCategoryRepository(db)
	ctx := context.Background()

	category := &entity.Category{
		Name:        "Technology",
		Slug:        "technology",
		Description: "Tech articles",
	}

	err := repo.Create(ctx, category)
	require.NoError(t, err)
	assert.NotZero(t, category.ID)
}

func TestCategoryRepository_FindByID(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewCategoryRepository(db)
	ctx := context.Background()

	category := &entity.Category{
		Name: "Technology",
		Slug: "technology",
	}
	err := repo.Create(ctx, category)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, category.ID)
	require.NoError(t, err)
	assert.Equal(t, category.Name, found.Name)
	assert.Equal(t, category.Slug, found.Slug)
}

func TestCategoryRepository_FindBySlug(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewCategoryRepository(db)
	ctx := context.Background()

	category := &entity.Category{
		Name: "Technology",
		Slug: "tech-slug",
	}
	err := repo.Create(ctx, category)
	require.NoError(t, err)

	found, err := repo.FindBySlug(ctx, "tech-slug")
	require.NoError(t, err)
	assert.Equal(t, category.ID, found.ID)
	assert.Equal(t, category.Name, found.Name)
}

func TestCategoryRepository_Update(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewCategoryRepository(db)
	ctx := context.Background()

	category := &entity.Category{
		Name: "Original",
		Slug: "original",
	}
	err := repo.Create(ctx, category)
	require.NoError(t, err)

	category.Name = "Updated"
	err = repo.Update(ctx, category)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, category.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", found.Name)
}

func TestCategoryRepository_Delete(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewCategoryRepository(db)
	ctx := context.Background()

	category := &entity.Category{
		Name: "ToDelete",
		Slug: "to-delete",
	}
	err := repo.Create(ctx, category)
	require.NoError(t, err)

	err = repo.Delete(ctx, category.ID)
	require.NoError(t, err)

	_, err = repo.FindByID(ctx, category.ID)
	assert.Error(t, err)
}

func TestCategoryRepository_List(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewCategoryRepository(db)
	ctx := context.Background()

	// カテゴリ作成
	for i := 1; i <= 3; i++ {
		category := &entity.Category{
			Name: fmt.Sprintf("Category%d", i),
			Slug: fmt.Sprintf("category-%d", i),
		}
		err := repo.Create(ctx, category)
		require.NoError(t, err)
	}

	categories, err := repo.List(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(categories), 3)
}

func TestCategoryRepository_Count(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewCategoryRepository(db)
	ctx := context.Background()

	// カテゴリ作成
	for i := 1; i <= 5; i++ {
		category := &entity.Category{
			Name: fmt.Sprintf("Category%d", i),
			Slug: fmt.Sprintf("category-%d", i),
		}
		err := repo.Create(ctx, category)
		require.NoError(t, err)
	}

	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 5)
}
