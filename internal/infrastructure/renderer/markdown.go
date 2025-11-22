package renderer

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownRenderer Markdownレンダラーインターフェース
type MarkdownRenderer interface {
	Render(source string) (string, error)
}

// markdownRenderer MarkdownRendererの実装
type markdownRenderer struct {
	md              goldmark.Markdown
	mermaidRenderer MermaidRenderer
}

// NewMarkdownRenderer 新しいMarkdownRendererを作成
func NewMarkdownRenderer(mermaidRenderer MermaidRenderer) MarkdownRenderer {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,           // GitHub Flavored Markdown
			extension.Table,         // テーブル
			extension.Strikethrough, // 取り消し線
			extension.Linkify,       // 自動リンク化
			extension.TaskList,      // タスクリスト
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // 見出しに自動ID付与
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // 改行を<br>に変換
			html.WithXHTML(),     // XHTML互換
		),
	)

	return &markdownRenderer{
		md:              md,
		mermaidRenderer: mermaidRenderer,
	}
}

// Render MarkdownをHTMLにレンダリング
func (r *markdownRenderer) Render(source string) (string, error) {
	if source == "" {
		return "", nil
	}

	// Mermaidコードブロックを抽出してSVGに変換
	processedSource, err := r.processMermaidBlocks(source)
	if err != nil {
		return "", fmt.Errorf("failed to process mermaid blocks: %w", err)
	}

	var buf bytes.Buffer
	if err := r.md.Convert([]byte(processedSource), &buf); err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	return buf.String(), nil
}

// processMermaidBlocks MermaidコードブロックをSVGに変換
func (r *markdownRenderer) processMermaidBlocks(source string) (string, error) {
	// ```mermaid ... ``` のパターンをマッチ
	re := regexp.MustCompile("(?s)```mermaid\\s*\\n(.*?)```")

	result := re.ReplaceAllStringFunc(source, func(match string) string {
		// Mermaidコードを抽出
		codeRe := regexp.MustCompile("(?s)```mermaid\\s*\\n(.*?)```")
		matches := codeRe.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}

		mermaidCode := matches[1]

		// SVGに変換
		svg, err := r.mermaidRenderer.RenderToSVG(mermaidCode)
		if err != nil {
			// エラー時は元のコードブロックを返す
			return match
		}

		return svg
	})

	return result, nil
}
