package renderer

import (
	"bytes"
	"fmt"
	htmllib "html"
	"regexp"
	"strings"

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
			// html.WithUnsafe()は使用しない（XSS対策）
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

	// Mermaidコードブロックを一時プレースホルダーに置換してSVGを保存
	processedSource, svgMap, err := r.extractMermaidBlocks(source)
	if err != nil {
		return "", fmt.Errorf("failed to process mermaid blocks: %w", err)
	}

	// Markdownをレンダリング（HTMLエスケープされる）
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(processedSource), &buf); err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	// プレースホルダーをSVGに置き換え
	// プレースホルダーはgoldmarkによってHTMLエスケープされるため、
	// エスケープされた形式で置換する必要があります
	result := buf.String()
	for placeholder, svg := range svgMap {
		escapedPlaceholder := htmllib.EscapeString(placeholder)
		result = strings.ReplaceAll(result, escapedPlaceholder, svg)
	}

	return result, nil
}

// extractMermaidBlocks MermaidコードブロックをSVGに変換してプレースホルダーに置換
// この関数は並行呼び出しに対して安全です。
// counterとsvgMapは各呼び出しごとにローカル変数として生成されるため、
// 複数のgoroutineから同時に呼び出されても競合状態は発生しません。
func (r *markdownRenderer) extractMermaidBlocks(source string) (string, map[string]string, error) {
	// ```mermaid ... ``` のパターンをマッチ
	re := regexp.MustCompile("(?s)```mermaid\\s*\\n(.*?)```")

	svgMap := make(map[string]string)
	counter := 0

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
			// Mermaidレンダリングエラー時は元のコードブロックを返す
			// これによりユーザーはMarkdown内でエラーを確認でき、
			// 記事全体のレンダリングは継続されます
			return match
		}

		// プレースホルダーを生成（一意性を保証）
		// ユーザーコンテンツとの衝突を防ぐため、特殊な接頭辞 + カウンター + SVG長を使用
		// MERMAIDSVGPLACEHOLDER形式はMarkdown記法と衝突せず、
		// 通常のコンテンツに含まれる可能性が極めて低いです
		placeholder := fmt.Sprintf("MERMAIDSVGPLACEHOLDER%dLEN%d", counter, len(svg))
		counter++
		svgMap[placeholder] = svg

		return placeholder
	})

	return result, svgMap, nil
}
