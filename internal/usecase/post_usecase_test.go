package usecase_test

import (
	"context"
	"testing"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/infrastructure/persistence"
	"my-blog-engine/internal/infrastructure/renderer"
	"my-blog-engine/internal/usecase"
	"my-blog-engine/tests/integration/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPostUseCase(t *testing.T) (usecase.PostUseCase, *entity.User, func()) {
	t.Helper()

	db, cleanup := testhelper.SetupTestDB(t)

	userRepo := persistence.NewUserRepository(db)
	postRepo := persistence.NewPostRepository(db)
	categoryRepo := persistence.NewCategoryRepository(db)
	tagRepo := persistence.NewTagRepository(db)

	mermaidRenderer := renderer.NewMockMermaidRenderer()
	mdRenderer := renderer.NewMarkdownRenderer(mermaidRenderer)

	postUseCase := usecase.NewPostUseCase(
		postRepo,
		categoryRepo,
		tagRepo,
		mdRenderer,
	)

	// テストユーザー作成
	ctx := context.Background()
	user := &entity.User{
		Username:     "testauthor",
		Email:        "author@example.com",
		PasswordHash: "hash",
		Role:         entity.RoleEditor,
		Status:       entity.StatusActive,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	return postUseCase, user, cleanup
}

func TestPostUseCase_Create(t *testing.T) {
	postUseCase, user, cleanup := setupPostUseCase(t)
	defer cleanup()

	ctx := context.Background()

	req := &usecase.CreatePostRequest{
		Title:    "Test Post",
		Slug:     "test-post",
		Content:  "# Test\n\nThis is a test post.",
		Status:   "draft",
		AuthorID: user.ID,
	}

	post, err := postUseCase.Create(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, req.Title, post.Title)
	assert.Equal(t, req.Slug, post.Slug)
	assert.NotEmpty(t, post.RenderedHTML)
	assert.Contains(t, post.RenderedHTML, "<h1")
	assert.Contains(t, post.RenderedHTML, "Test")
}

func TestPostUseCase_Update(t *testing.T) {
	postUseCase, user, cleanup := setupPostUseCase(t)
	defer cleanup()

	ctx := context.Background()

	// 記事作成
	req := &usecase.CreatePostRequest{
		Title:    "Original Title",
		Slug:     "original-slug",
		Content:  "Original content",
		Status:   "draft",
		AuthorID: user.ID,
	}

	post, err := postUseCase.Create(ctx, req)
	require.NoError(t, err)

	// 更新
	newTitle := "Updated Title"
	updateReq := &usecase.UpdatePostRequest{
		Title: &newTitle,
	}

	updated, err := postUseCase.Update(ctx, post.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, newTitle, updated.Title)
	assert.Equal(t, post.Slug, updated.Slug) // 変更されていない
}

func TestPostUseCase_Delete(t *testing.T) {
	postUseCase, user, cleanup := setupPostUseCase(t)
	defer cleanup()

	ctx := context.Background()

	// 記事作成
	req := &usecase.CreatePostRequest{
		Title:    "Test Post",
		Slug:     "test-post",
		Content:  "Test content",
		Status:   "draft",
		AuthorID: user.ID,
	}

	post, err := postUseCase.Create(ctx, req)
	require.NoError(t, err)

	// 削除
	err = postUseCase.Delete(ctx, post.ID)
	assert.NoError(t, err)

	// 取得失敗確認
	_, err = postUseCase.GetByID(ctx, post.ID)
	assert.Error(t, err)
}

func TestPostUseCase_GetBySlug(t *testing.T) {
	postUseCase, user, cleanup := setupPostUseCase(t)
	defer cleanup()

	ctx := context.Background()

	req := &usecase.CreatePostRequest{
		Title:    "Test Post",
		Slug:     "unique-test-slug",
		Content:  "Test content",
		Status:   "draft",
		AuthorID: user.ID,
	}

	created, err := postUseCase.Create(ctx, req)
	require.NoError(t, err)

	// スラッグで取得
	post, err := postUseCase.GetBySlug(ctx, "unique-test-slug")
	require.NoError(t, err)
	assert.Equal(t, created.ID, post.ID)
	assert.Equal(t, created.Title, post.Title)
}

func TestPostUseCase_Publish(t *testing.T) {
	postUseCase, user, cleanup := setupPostUseCase(t)
	defer cleanup()

	ctx := context.Background()

	req := &usecase.CreatePostRequest{
		Title:    "Test Post",
		Slug:     "test-post",
		Content:  "Test content",
		Status:   "draft",
		AuthorID: user.ID,
	}

	post, err := postUseCase.Create(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusDraft, post.Status)

	// 公開
	err = postUseCase.Publish(ctx, post.ID)
	require.NoError(t, err)

	// 確認
	published, err := postUseCase.GetByID(ctx, post.ID)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusPublished, published.Status)
	assert.NotNil(t, published.PublishedAt)
}

func TestPostUseCase_Unpublish(t *testing.T) {
	postUseCase, user, cleanup := setupPostUseCase(t)
	defer cleanup()

	ctx := context.Background()

	req := &usecase.CreatePostRequest{
		Title:    "Test Post",
		Slug:     "test-post",
		Content:  "Test content",
		Status:   "published",
		AuthorID: user.ID,
	}

	post, err := postUseCase.Create(ctx, req)
	require.NoError(t, err)

	// 公開して非公開に
	err = postUseCase.Publish(ctx, post.ID)
	require.NoError(t, err)

	err = postUseCase.Unpublish(ctx, post.ID)
	require.NoError(t, err)

	// 確認
	unpublished, err := postUseCase.GetByID(ctx, post.ID)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusDraft, unpublished.Status)
}
