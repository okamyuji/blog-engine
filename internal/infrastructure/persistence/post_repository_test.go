package persistence_test

import (
	"context"
	"testing"
	"time"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/domain/repository"
	"my-blog-engine/internal/infrastructure/persistence"
	"my-blog-engine/tests/integration/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPostTest(t *testing.T) (repository.PostRepository, *entity.User, func()) {
	t.Helper()

	db, cleanup := testhelper.SetupTestDB(t)

	postRepo := persistence.NewPostRepository(db)
	userRepo := persistence.NewUserRepository(db)

	ctx := context.Background()

	// テストユーザー作成
	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         entity.RoleEditor,
		Status:       entity.StatusActive,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	return postRepo, user, cleanup
}

func TestPostRepository_Create(t *testing.T) {
	repo, user, cleanup := setupPostTest(t)
	defer cleanup()

	ctx := context.Background()

	post := &entity.Post{
		Title:        "Test Post",
		Slug:         "test-post",
		Content:      "Test content",
		RenderedHTML: "<p>Test content</p>",
		Status:       entity.StatusDraft,
		AuthorID:     user.ID,
	}

	err := repo.Create(ctx, post)
	require.NoError(t, err)
	assert.NotZero(t, post.ID)
}

func TestPostRepository_FindByID(t *testing.T) {
	repo, user, cleanup := setupPostTest(t)
	defer cleanup()

	ctx := context.Background()

	post := &entity.Post{
		Title:        "Test Post",
		Slug:         "test-post",
		Content:      "Test content",
		RenderedHTML: "<p>Test content</p>",
		Status:       entity.StatusDraft,
		AuthorID:     user.ID,
	}

	err := repo.Create(ctx, post)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, post.ID)
	require.NoError(t, err)
	assert.Equal(t, post.Title, found.Title)
	assert.Equal(t, post.Slug, found.Slug)
	assert.NotNil(t, found.Author)
	assert.Equal(t, user.Username, found.Author.Username)
}

func TestPostRepository_FindBySlug(t *testing.T) {
	repo, user, cleanup := setupPostTest(t)
	defer cleanup()

	ctx := context.Background()

	post := &entity.Post{
		Title:        "Test Post",
		Slug:         "unique-slug",
		Content:      "Test content",
		RenderedHTML: "<p>Test content</p>",
		Status:       entity.StatusDraft,
		AuthorID:     user.ID,
	}

	err := repo.Create(ctx, post)
	require.NoError(t, err)

	found, err := repo.FindBySlug(ctx, "unique-slug")
	require.NoError(t, err)
	assert.Equal(t, post.ID, found.ID)
	assert.Equal(t, post.Title, found.Title)
}

func TestPostRepository_Update(t *testing.T) {
	repo, user, cleanup := setupPostTest(t)
	defer cleanup()

	ctx := context.Background()

	post := &entity.Post{
		Title:        "Original Title",
		Slug:         "original-slug",
		Content:      "Original content",
		RenderedHTML: "<p>Original content</p>",
		Status:       entity.StatusDraft,
		AuthorID:     user.ID,
	}

	err := repo.Create(ctx, post)
	require.NoError(t, err)

	post.Title = "Updated Title"
	post.Content = "Updated content"
	err = repo.Update(ctx, post)
	require.NoError(t, err)

	updated, err := repo.FindByID(ctx, post.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, "Updated content", updated.Content)
}

func TestPostRepository_Delete(t *testing.T) {
	repo, user, cleanup := setupPostTest(t)
	defer cleanup()

	ctx := context.Background()

	post := &entity.Post{
		Title:        "Test Post",
		Slug:         "test-post",
		Content:      "Test content",
		RenderedHTML: "<p>Test content</p>",
		Status:       entity.StatusDraft,
		AuthorID:     user.ID,
	}

	err := repo.Create(ctx, post)
	require.NoError(t, err)

	err = repo.Delete(ctx, post.ID)
	require.NoError(t, err)

	_, err = repo.FindByID(ctx, post.ID)
	assert.Error(t, err)
}

func TestPostRepository_ListPublished(t *testing.T) {
	repo, user, cleanup := setupPostTest(t)
	defer cleanup()

	ctx := context.Background()

	// ドラフト作成
	draft := &entity.Post{
		Title:        "Draft Post",
		Slug:         "draft-post",
		Content:      "Draft content",
		RenderedHTML: "<p>Draft content</p>",
		Status:       entity.StatusDraft,
		AuthorID:     user.ID,
	}
	err := repo.Create(ctx, draft)
	require.NoError(t, err)

	// 公開記事作成
	now := time.Now()
	published := &entity.Post{
		Title:        "Published Post",
		Slug:         "published-post",
		Content:      "Published content",
		RenderedHTML: "<p>Published content</p>",
		Status:       entity.StatusPublished,
		AuthorID:     user.ID,
		PublishedAt:  &now,
	}
	err = repo.Create(ctx, published)
	require.NoError(t, err)

	posts, err := repo.ListPublished(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(posts), 1)

	// すべて公開済みか確認
	for _, p := range posts {
		assert.Equal(t, entity.StatusPublished, p.Status)
	}
}

func TestPostRepository_AddTags(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	userRepo := persistence.NewUserRepository(db)
	postRepo := persistence.NewPostRepository(db)
	tagRepo := persistence.NewTagRepository(db)

	// ユーザー作成
	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         entity.RoleEditor,
		Status:       entity.StatusActive,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// 記事作成
	post := &entity.Post{
		Title:        "Test Post",
		Slug:         "test-post",
		Content:      "Test content",
		RenderedHTML: "<p>Test content</p>",
		Status:       entity.StatusDraft,
		AuthorID:     user.ID,
	}
	err = postRepo.Create(ctx, post)
	require.NoError(t, err)

	// タグ作成
	tag1 := &entity.Tag{Name: "Go", Slug: "go"}
	tag2 := &entity.Tag{Name: "Web", Slug: "web"}
	err = tagRepo.Create(ctx, tag1)
	require.NoError(t, err)
	err = tagRepo.Create(ctx, tag2)
	require.NoError(t, err)

	// タグ追加
	err = postRepo.AddTags(ctx, post.ID, []int64{tag1.ID, tag2.ID})
	require.NoError(t, err)

	// 確認
	tags, err := postRepo.GetTags(ctx, post.ID)
	require.NoError(t, err)
	assert.Len(t, tags, 2)
}

func TestPostRepository_RemoveTags(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	userRepo := persistence.NewUserRepository(db)
	postRepo := persistence.NewPostRepository(db)
	tagRepo := persistence.NewTagRepository(db)

	// ユーザー作成
	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         entity.RoleEditor,
		Status:       entity.StatusActive,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// 記事作成
	post := &entity.Post{
		Title:        "Test Post",
		Slug:         "test-post",
		Content:      "Test content",
		RenderedHTML: "<p>Test content</p>",
		Status:       entity.StatusDraft,
		AuthorID:     user.ID,
	}
	err = postRepo.Create(ctx, post)
	require.NoError(t, err)

	// タグ作成と追加
	tag := &entity.Tag{Name: "Go", Slug: "go"}
	err = tagRepo.Create(ctx, tag)
	require.NoError(t, err)

	err = postRepo.AddTags(ctx, post.ID, []int64{tag.ID})
	require.NoError(t, err)

	// タグ削除
	err = postRepo.RemoveTags(ctx, post.ID, []int64{tag.ID})
	require.NoError(t, err)

	// 確認
	tags, err := postRepo.GetTags(ctx, post.ID)
	require.NoError(t, err)
	assert.Len(t, tags, 0)
}
