package renderer

import (
	"bytes"
	"fmt"
	htmllib "html"
	"log"
	"regexp"
	"strings"

	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

// mermaidBlockRe Mermaidコードブロックマッチ用の正規表現（パフォーマンス最適化のため事前コンパイル）
var mermaidBlockRe = regexp.MustCompile("(?s)```mermaid\\s*\\n(.*?)```")

// MarkdownRenderer Markdownレンダラーインターフェース
type MarkdownRenderer interface {
	Render(source string) (string, error)
}

// markdownRenderer MarkdownRendererの実装
type markdownRenderer struct {
	parser          parser.Parser
	renderer        html.Renderer
	mermaidRenderer MermaidRenderer
}

// NewMarkdownRenderer 新しいMarkdownRendererを作成
func NewMarkdownRenderer(mermaidRenderer MermaidRenderer) MarkdownRenderer {
	p := parser.New(
		parser.WithExtensions(
			extension.GFMParser,           // GitHub Flavored Markdown
			extension.TableParser,         // テーブル
			extension.StrikethroughParser, // 取り消し線
			extension.LinkifyParser,       // 自動リンク化
			extension.TaskListItemParser,  // タスクリスト
		),
		parser.WithAutoHeadingID(), // 見出しに自動ID付与
	)
	r := html.New(
		html.WithExtensions(
			extension.GFMHTMLRenderer,
			extension.TableHTMLRenderer,
			extension.StrikethroughHTMLRenderer,
			extension.TaskListItemHTMLRenderer,
		),
		html.WithHardWraps(), // 改行を<br>に変換
		html.WithXHTML(),     // XHTML互換
		// html.WithUnsafe()は使用しない（XSS対策）
	)

	return &markdownRenderer{
		parser:          p,
		renderer:        r,
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
	src := []byte(processedSource)
	if err := r.renderer.Render(&buf, src, r.parser.Parse(src)); err != nil {
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
	svgMap := make(map[string]string)
	counter := 0

	result := mermaidBlockRe.ReplaceAllStringFunc(source, func(match string) string {
		// Mermaidコードを抽出（パッケージレベルの正規表現を使用）
		matches := mermaidBlockRe.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}

		mermaidCode := matches[1]

		// SVGに変換
		svg, err := r.mermaidRenderer.RenderToSVG(mermaidCode)
		if err != nil {
			// エラーをログに記録（本番環境でのデバッグ用）
			log.Printf("Mermaid rendering failed: %v (code preview: %.50s...)", err, mermaidCode)
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
