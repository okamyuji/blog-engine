package handler

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"my-blog-engine/internal/usecase"
)

// PublicHandler 公開ページのハンドラー
type PublicHandler struct {
	postUseCase     usecase.PostUseCase
	categoryUseCase usecase.CategoryUseCase
	templates       *template.Template
}

// NewPublicHandler PublicHandlerのコンストラクタ
func NewPublicHandler(
	postUseCase usecase.PostUseCase,
	categoryUseCase usecase.CategoryUseCase,
) *PublicHandler {
	// テンプレートファイルを個別にパース
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Printf("Warning: Failed to parse templates: %v", err)
		tmpl = template.New("fallback")
	}

	return &PublicHandler{
		postUseCase:     postUseCase,
		categoryUseCase: categoryUseCase,
		templates:       tmpl,
	}
}

// Home ホームページ表示
func (h *PublicHandler) Home(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// 公開投稿一覧取得
	posts, _, err := h.postUseCase.ListPublished(ctx, 10, 0)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	// カテゴリ一覧取得
	categories, err := h.categoryUseCase.List(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":      "Blog Home",
		"Posts":      posts,
		"Categories": categories,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
