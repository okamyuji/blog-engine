package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"my-blog-engine/internal/interface/presenter"
	"my-blog-engine/internal/usecase"
)

// TagHandler タグハンドラー
type TagHandler struct {
	tagUseCase usecase.TagUseCase
}

// NewTagHandler 新しいTagHandlerを作成
func NewTagHandler(tagUseCase usecase.TagUseCase) *TagHandler {
	return &TagHandler{
		tagUseCase: tagUseCase,
	}
}

// Create タグ作成ハンドラー
func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tag, err := h.tagUseCase.Create(r.Context(), &req)
	if err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to create tag")
		return
	}

	presenter.JSONResponse(w, http.StatusCreated, tag)
}

// Update タグ更新ハンドラー
func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	var req usecase.UpdateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tag, err := h.tagUseCase.Update(r.Context(), id, &req)
	if err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to update tag")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, tag)
}

// Delete タグ削除ハンドラー
func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	if err := h.tagUseCase.Delete(r.Context(), id); err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to delete tag")
		return
	}

	presenter.JSONSuccess(w, nil, "Tag deleted successfully")
}

// GetByID ID指定でタグ取得ハンドラー
func (h *TagHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	tag, err := h.tagUseCase.GetByID(r.Context(), id)
	if err != nil {
		presenter.JSONError(w, http.StatusNotFound, "Tag not found")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, tag)
}

// GetBySlug スラッグ指定でタグ取得ハンドラー
func (h *TagHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		presenter.JSONError(w, http.StatusBadRequest, "Slug required")
		return
	}

	tag, err := h.tagUseCase.GetBySlug(r.Context(), slug)
	if err != nil {
		presenter.JSONError(w, http.StatusNotFound, "Tag not found")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, tag)
}

// List タグ一覧ハンドラー
func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	tags, err := h.tagUseCase.List(r.Context())
	if err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Failed to list tags")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, tags)
}
