package handler

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"my-blog-engine/internal/domain/entity"
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

// PostView テンプレート用の投稿ビュー
type PostView struct {
	*entity.Post
	SafeHTML template.HTML
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

	// 投稿をテンプレート用に変換
	// RenderedHTMLをtemplate.HTMLに変換することは安全です。
	// なぜなら、RenderedHTMLはmarkdownレンダラー（goldmark）によって
	// すでにサニタイズされており、HTMLエスケープとプレースホルダーベースの
	// SVG挿入によってXSS攻撃から保護されているためです。
	postViews := make([]PostView, len(posts))
	for i, post := range posts {
		postViews[i] = PostView{
			Post:     post,
			SafeHTML: template.HTML(post.RenderedHTML),
		}
	}

	data := map[string]interface{}{
		"Title":      "Blog Home",
		"Posts":      postViews,
		"Categories": categories,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
