package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"my-blog-engine/internal/interface/presenter"
	"my-blog-engine/internal/usecase"
)

// CategoryHandler カテゴリハンドラー
type CategoryHandler struct {
	categoryUseCase usecase.CategoryUseCase
}

// NewCategoryHandler 新しいCategoryHandlerを作成
func NewCategoryHandler(categoryUseCase usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}

// Create カテゴリ作成ハンドラー
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category, err := h.categoryUseCase.Create(r.Context(), &req)
	if err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	presenter.JSONResponse(w, http.StatusCreated, category)
}

// Update カテゴリ更新ハンドラー
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var req usecase.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category, err := h.categoryUseCase.Update(r.Context(), id, &req)
	if err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, category)
}

// Delete カテゴリ削除ハンドラー
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.categoryUseCase.Delete(r.Context(), id); err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	presenter.JSONSuccess(w, nil, "Category deleted successfully")
}

// GetByID ID指定でカテゴリ取得ハンドラー
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.categoryUseCase.GetByID(r.Context(), id)
	if err != nil {
		presenter.JSONError(w, http.StatusNotFound, "Category not found")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, category)
}

// GetBySlug スラッグ指定でカテゴリ取得ハンドラー
func (h *CategoryHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		presenter.JSONError(w, http.StatusBadRequest, "Slug required")
		return
	}

	category, err := h.categoryUseCase.GetBySlug(r.Context(), slug)
	if err != nil {
		presenter.JSONError(w, http.StatusNotFound, "Category not found")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, category)
}

// List カテゴリ一覧ハンドラー
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryUseCase.List(r.Context())
	if err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to list categories")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, categories)
}
