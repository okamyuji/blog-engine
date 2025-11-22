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

func TestTagUseCase_Create(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	tagRepo := persistence.NewTagRepository(db)
	tagUseCase := usecase.NewTagUseCase(tagRepo)

	ctx := context.Background()

	req := &usecase.CreateTagRequest{
		Name: "Golang",
		Slug: "golang",
	}

	tag, err := tagUseCase.Create(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, req.Name, tag.Name)
	assert.Equal(t, req.Slug, tag.Slug)
}

func TestTagUseCase_Update(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	tagRepo := persistence.NewTagRepository(db)
	tagUseCase := usecase.NewTagUseCase(tagRepo)

	ctx := context.Background()

	// タグ作成
	createReq := &usecase.CreateTagRequest{
		Name: "Original",
		Slug: "original",
	}
	tag, err := tagUseCase.Create(ctx, createReq)
	require.NoError(t, err)

	// 更新
	newName := "Updated"
	updateReq := &usecase.UpdateTagRequest{
		Name: &newName,
	}

	updated, err := tagUseCase.Update(ctx, tag.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, tag.Slug, updated.Slug)
}

func TestTagUseCase_Delete(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	tagRepo := persistence.NewTagRepository(db)
	tagUseCase := usecase.NewTagUseCase(tagRepo)

	ctx := context.Background()

	// タグ作成
	createReq := &usecase.CreateTagRequest{
		Name: "ToDelete",
		Slug: "to-delete",
	}
	tag, err := tagUseCase.Create(ctx, createReq)
	require.NoError(t, err)

	// 削除
	err = tagUseCase.Delete(ctx, tag.ID)
	assert.NoError(t, err)

	// 取得失敗確認
	_, err = tagUseCase.GetByID(ctx, tag.ID)
	assert.Error(t, err)
}

func TestTagUseCase_GetBySlug(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	tagRepo := persistence.NewTagRepository(db)
	tagUseCase := usecase.NewTagUseCase(tagRepo)

	ctx := context.Background()

	// タグ作成
	createReq := &usecase.CreateTagRequest{
		Name: "Test",
		Slug: "test-slug",
	}
	created, err := tagUseCase.Create(ctx, createReq)
	require.NoError(t, err)

	// スラッグで取得
	tag, err := tagUseCase.GetBySlug(ctx, "test-slug")
	require.NoError(t, err)
	assert.Equal(t, created.ID, tag.ID)
	assert.Equal(t, created.Name, tag.Name)
}

func TestTagUseCase_List(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	tagRepo := persistence.NewTagRepository(db)
	tagUseCase := usecase.NewTagUseCase(tagRepo)

	ctx := context.Background()

	// タグ作成
	tags := []string{"Tag1", "Tag2", "Tag3"}
	for _, name := range tags {
		req := &usecase.CreateTagRequest{
			Name: name,
			Slug: name,
		}
		_, err := tagUseCase.Create(ctx, req)
		require.NoError(t, err)
	}

	// 一覧取得
	list, err := tagUseCase.List(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 3)
}
