package usecase

import (
	"context"
	"fmt"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/domain/repository"
	"my-blog-engine/internal/infrastructure/renderer"
)

// PostUseCase 記事ユースケースのインターフェース
type PostUseCase interface {
	Create(ctx context.Context, req *CreatePostRequest) (*entity.Post, error)
	Update(ctx context.Context, id int64, req *UpdatePostRequest) (*entity.Post, error)
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*entity.Post, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Post, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Post, int, error)
	ListPublished(ctx context.Context, limit, offset int) ([]*entity.Post, int, error)
	ListByCategory(ctx context.Context, categorySlug string, limit, offset int) ([]*entity.Post, int, error)
	ListByTag(ctx context.Context, tagSlug string, limit, offset int) ([]*entity.Post, int, error)
	Publish(ctx context.Context, id int64) error
	Unpublish(ctx context.Context, id int64) error
}

// CreatePostRequest 記事作成リクエスト
type CreatePostRequest struct {
	Title      string  `json:"title"`
	Slug       string  `json:"slug"`
	Content    string  `json:"content"`
	Status     string  `json:"status"`
	AuthorID   int64   `json:"authorId"`
	CategoryID *int64  `json:"categoryId"`
	TagIDs     []int64 `json:"tagIds"`
}

// UpdatePostRequest 記事更新リクエスト
type UpdatePostRequest struct {
	Title      *string `json:"title"`
	Slug       *string `json:"slug"`
	Content    *string `json:"content"`
	Status     *string `json:"status"`
	CategoryID *int64  `json:"categoryId"`
	TagIDs     []int64 `json:"tagIds"`
}

// postUseCase PostUseCaseの実装
type postUseCase struct {
	postRepo     repository.PostRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
	mdRenderer   renderer.MarkdownRenderer
}

// NewPostUseCase 新しいPostUseCaseを作成
func NewPostUseCase(
	postRepo repository.PostRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
	mdRenderer renderer.MarkdownRenderer,
) PostUseCase {
	return &postUseCase{
		postRepo:     postRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		mdRenderer:   mdRenderer,
	}
}

// Create 新しい記事を作成
func (u *postUseCase) Create(ctx context.Context, req *CreatePostRequest) (*entity.Post, error) {
	if req.Title == "" || req.Slug == "" || req.Content == "" {
		return nil, fmt.Errorf("title, slug, and content are required")
	}

	// Markdownレンダリング
	renderedHTML, err := u.mdRenderer.Render(req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to render markdown: %w", err)
	}

	// 記事作成
	post := &entity.Post{
		Title:        req.Title,
		Slug:         req.Slug,
		Content:      req.Content,
		RenderedHTML: renderedHTML,
		Status:       entity.PostStatus(req.Status),
		AuthorID:     req.AuthorID,
		CategoryID:   req.CategoryID,
	}

	if err := u.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// タグ追加
	if len(req.TagIDs) > 0 {
		if err := u.postRepo.AddTags(ctx, post.ID, req.TagIDs); err != nil {
			return nil, fmt.Errorf("failed to add tags: %w", err)
		}
	}

	// 作成した記事を取得(リレーション含む)
	return u.postRepo.FindByID(ctx, post.ID)
}

// Update 記事を更新
func (u *postUseCase) Update(ctx context.Context, id int64, req *UpdatePostRequest) (*entity.Post, error) {
	// 既存の記事を取得
	post, err := u.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find post: %w", err)
	}

	// 更新
	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Slug != nil {
		post.Slug = *req.Slug
	}
	if req.Content != nil {
		post.Content = *req.Content
		// Markdownレンダリング
		renderedHTML, err := u.mdRenderer.Render(post.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to render markdown: %w", err)
		}
		post.RenderedHTML = renderedHTML
	}
	if req.Status != nil {
		post.Status = entity.PostStatus(*req.Status)
	}
	if req.CategoryID != nil {
		post.CategoryID = req.CategoryID
	}

	if err := u.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// タグ更新
	if req.TagIDs != nil {
		// 既存のタグを取得
		existingTags, err := u.postRepo.GetTags(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing tags: %w", err)
		}

		// 既存のタグIDを抽出
		existingTagIDs := make([]int64, len(existingTags))
		for i, tag := range existingTags {
			existingTagIDs[i] = tag.ID
		}

		// 既存のタグを削除
		if len(existingTagIDs) > 0 {
			if err := u.postRepo.RemoveTags(ctx, id, existingTagIDs); err != nil {
				return nil, fmt.Errorf("failed to remove tags: %w", err)
			}
		}

		// 新しいタグを追加
		if len(req.TagIDs) > 0 {
			if err := u.postRepo.AddTags(ctx, id, req.TagIDs); err != nil {
				return nil, fmt.Errorf("failed to add tags: %w", err)
			}
		}
	}

	// 更新した記事を取得
	return u.postRepo.FindByID(ctx, id)
}

// Delete 記事を削除
func (u *postUseCase) Delete(ctx context.Context, id int64) error {
	if err := u.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

// GetByID IDで記事を取得
func (u *postUseCase) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	post, err := u.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find post: %w", err)
	}
	return post, nil
}

// GetBySlug スラッグで記事を取得
func (u *postUseCase) GetBySlug(ctx context.Context, slug string) (*entity.Post, error) {
	post, err := u.postRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to find post: %w", err)
	}
	return post, nil
}

// List 記事一覧を取得
func (u *postUseCase) List(ctx context.Context, limit, offset int) ([]*entity.Post, int, error) {
	posts, err := u.postRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list posts: %w", err)
	}

	count, err := u.postRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	return posts, count, nil
}

// ListPublished 公開済み記事一覧を取得
func (u *postUseCase) ListPublished(ctx context.Context, limit, offset int) ([]*entity.Post, int, error) {
	posts, err := u.postRepo.ListPublished(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list published posts: %w", err)
	}

	count, err := u.postRepo.CountPublished(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count published posts: %w", err)
	}

	return posts, count, nil
}

// ListByCategory カテゴリ別記事一覧を取得
func (u *postUseCase) ListByCategory(ctx context.Context, categorySlug string, limit, offset int) ([]*entity.Post, int, error) {
	category, err := u.categoryRepo.FindBySlug(ctx, categorySlug)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find category: %w", err)
	}

	posts, err := u.postRepo.ListByCategory(ctx, category.ID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list posts by category: %w", err)
	}

	// TODO: カテゴリ別の記事数カウントメソッド追加
	return posts, len(posts), nil
}

// ListByTag タグ別記事一覧を取得
func (u *postUseCase) ListByTag(ctx context.Context, tagSlug string, limit, offset int) ([]*entity.Post, int, error) {
	tag, err := u.tagRepo.FindBySlug(ctx, tagSlug)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find tag: %w", err)
	}

	posts, err := u.postRepo.ListByTag(ctx, tag.ID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list posts by tag: %w", err)
	}

	// TODO: タグ別の記事数カウントメソッド追加
	return posts, len(posts), nil
}

// Publish 記事を公開
func (u *postUseCase) Publish(ctx context.Context, id int64) error {
	post, err := u.postRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find post: %w", err)
	}

	if post.Status == entity.StatusPublished {
		return nil // すでに公開済み
	}

	post.Publish()
	if err := u.postRepo.Update(ctx, post); err != nil {
		return fmt.Errorf("failed to publish post: %w", err)
	}

	return nil
}

// Unpublish 記事を非公開にする
func (u *postUseCase) Unpublish(ctx context.Context, id int64) error {
	post, err := u.postRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find post: %w", err)
	}

	if post.Status == entity.StatusDraft {
		return nil // すでに下書き
	}

	post.Unpublish()
	if err := u.postRepo.Update(ctx, post); err != nil {
		return fmt.Errorf("failed to unpublish post: %w", err)
	}

	return nil
}
