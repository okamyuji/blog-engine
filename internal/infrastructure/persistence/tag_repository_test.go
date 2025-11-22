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

func TestTagRepository_Create(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTagRepository(db)
	ctx := context.Background()

	tag := &entity.Tag{
		Name: "Golang",
		Slug: "golang",
	}

	err := repo.Create(ctx, tag)
	require.NoError(t, err)
	assert.NotZero(t, tag.ID)
}

func TestTagRepository_FindByID(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTagRepository(db)
	ctx := context.Background()

	tag := &entity.Tag{
		Name: "Golang",
		Slug: "golang",
	}
	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, tag.ID)
	require.NoError(t, err)
	assert.Equal(t, tag.Name, found.Name)
	assert.Equal(t, tag.Slug, found.Slug)
}

func TestTagRepository_FindBySlug(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTagRepository(db)
	ctx := context.Background()

	tag := &entity.Tag{
		Name: "Golang",
		Slug: "golang-slug",
	}
	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	found, err := repo.FindBySlug(ctx, "golang-slug")
	require.NoError(t, err)
	assert.Equal(t, tag.ID, found.ID)
	assert.Equal(t, tag.Name, found.Name)
}

func TestTagRepository_FindByIDs(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTagRepository(db)
	ctx := context.Background()

	// タグ作成
	tag1 := &entity.Tag{Name: "Tag1", Slug: "tag1"}
	tag2 := &entity.Tag{Name: "Tag2", Slug: "tag2"}
	err := repo.Create(ctx, tag1)
	require.NoError(t, err)
	err = repo.Create(ctx, tag2)
	require.NoError(t, err)

	// 複数IDで検索
	found, err := repo.FindByIDs(ctx, []int64{tag1.ID, tag2.ID})
	require.NoError(t, err)
	assert.Len(t, found, 2)
}

func TestTagRepository_FindByIDs_Empty(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTagRepository(db)
	ctx := context.Background()

	// 空のIDリスト
	found, err := repo.FindByIDs(ctx, []int64{})
	require.NoError(t, err)
	assert.Empty(t, found)
}

func TestTagRepository_Update(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTagRepository(db)
	ctx := context.Background()

	tag := &entity.Tag{
		Name: "Original",
		Slug: "original",
	}
	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	tag.Name = "Updated"
	err = repo.Update(ctx, tag)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, tag.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", found.Name)
}

func TestTagRepository_Delete(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTagRepository(db)
	ctx := context.Background()

	tag := &entity.Tag{
		Name: "ToDelete",
		Slug: "to-delete",
	}
	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	err = repo.Delete(ctx, tag.ID)
	require.NoError(t, err)

	_, err = repo.FindByID(ctx, tag.ID)
	assert.Error(t, err)
}

func TestTagRepository_List(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTagRepository(db)
	ctx := context.Background()

	// タグ作成
	for i := 1; i <= 3; i++ {
		tag := &entity.Tag{
			Name: fmt.Sprintf("Tag%d", i),
			Slug: fmt.Sprintf("tag-%d", i),
		}
		err := repo.Create(ctx, tag)
		require.NoError(t, err)
	}

	tags, err := repo.List(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(tags), 3)
}

func TestTagRepository_Count(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewTagRepository(db)
	ctx := context.Background()

	// タグ作成
	for i := 1; i <= 5; i++ {
		tag := &entity.Tag{
			Name: fmt.Sprintf("Tag%d", i),
			Slug: fmt.Sprintf("tag-%d", i),
		}
		err := repo.Create(ctx, tag)
		require.NoError(t, err)
	}

	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 5)
}
