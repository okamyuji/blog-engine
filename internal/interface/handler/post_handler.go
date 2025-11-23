package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/interface/middleware"
	"my-blog-engine/internal/interface/presenter"
	"my-blog-engine/internal/usecase"
)

// PostHandler 記事ハンドラー
type PostHandler struct {
	postUseCase usecase.PostUseCase
}

// NewPostHandler 新しいPostHandlerを作成
func NewPostHandler(postUseCase usecase.PostUseCase) *PostHandler {
	return &PostHandler{
		postUseCase: postUseCase,
	}
}

// Create 記事作成ハンドラー
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		presenter.JSONError(w, http.StatusUnauthorized, "User not found")
		return
	}

	var req usecase.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.AuthorID = user.ID

	post, err := h.postUseCase.Create(r.Context(), &req)
	if err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	presenter.JSONResponse(w, http.StatusCreated, post)
}

// Update 記事更新ハンドラー
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var req usecase.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	post, err := h.postUseCase.Update(r.Context(), id, &req)
	if err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to update post")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, post)
}

// Delete 記事削除ハンドラー
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	if err := h.postUseCase.Delete(r.Context(), id); err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	presenter.JSONSuccess(w, nil, "Post deleted successfully")
}

// GetByID ID指定で記事取得ハンドラー（公開記事のみ）
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	post, err := h.postUseCase.GetByID(r.Context(), id)
	if err != nil {
		presenter.JSONError(w, http.StatusNotFound, "Post not found")
		return
	}

	// 公開記事のみ返却（下書きは認証が必要）
	if post.Status != entity.StatusPublished {
		presenter.JSONError(w, http.StatusNotFound, "Post not found")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, post)
}

// GetBySlug スラッグ指定で記事取得ハンドラー（公開記事のみ）
func (h *PostHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		presenter.JSONError(w, http.StatusBadRequest, "Slug is required")
		return
	}

	post, err := h.postUseCase.GetBySlug(r.Context(), slug)
	if err != nil {
		presenter.JSONError(w, http.StatusNotFound, "Post not found")
		return
	}

	// 公開記事のみ返却（下書きは認証が必要）
	if post.Status != entity.StatusPublished {
		presenter.JSONError(w, http.StatusNotFound, "Post not found")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, post)
}

// List 記事一覧ハンドラー
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, count, err := h.postUseCase.List(r.Context(), limit, offset)
	if err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to list posts")
		return
	}

	response := map[string]interface{}{
		"posts":  posts,
		"total":  count,
		"limit":  limit,
		"offset": offset,
	}

	presenter.JSONResponse(w, http.StatusOK, response)
}

// ListPublished 公開記事一覧ハンドラー
func (h *PostHandler) ListPublished(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, count, err := h.postUseCase.ListPublished(r.Context(), limit, offset)
	if err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to list published posts")
		return
	}

	response := map[string]interface{}{
		"posts":  posts,
		"total":  count,
		"limit":  limit,
		"offset": offset,
	}

	presenter.JSONResponse(w, http.StatusOK, response)
}

// Publish 記事公開ハンドラー
func (h *PostHandler) Publish(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	if err := h.postUseCase.Publish(r.Context(), id); err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to publish post")
		return
	}

	presenter.JSONSuccess(w, nil, "Post published successfully")
}

// Unpublish 記事非公開ハンドラー
func (h *PostHandler) Unpublish(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	if err := h.postUseCase.Unpublish(r.Context(), id); err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to unpublish post")
		return
	}

	presenter.JSONSuccess(w, nil, "Post unpublished successfully")
}
