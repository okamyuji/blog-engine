package handler

import (
	"net/http"

	"my-blog-engine/internal/interface/presenter"

	"github.com/uptrace/bun"
)

// HealthHandler ヘルスチェックハンドラー
type HealthHandler struct {
	db *bun.DB
}

// NewHealthHandler 新しいHealthHandlerを作成
func NewHealthHandler(db *bun.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Check ヘルスチェック
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	// データベース接続チェック
	if err := h.db.Ping(); err != nil {
		presenter.JSONError(w, http.StatusServiceUnavailable, "Database connection failed")
		return
	}

	response := map[string]string{
		"status":   "healthy",
		"database": "connected",
	}

	presenter.JSONResponse(w, http.StatusOK, response)
}
