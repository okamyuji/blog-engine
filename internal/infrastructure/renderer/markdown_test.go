package renderer_test

import (
	"testing"

	"my-blog-engine/internal/infrastructure/renderer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownRenderer_Render(t *testing.T) {
	mermaidRenderer := renderer.NewMockMermaidRenderer()
	mdRenderer := renderer.NewMarkdownRenderer(mermaidRenderer)

	tests := []struct {
		name     string
		source   string
		contains []string
	}{
		{
			name:     "simple markdown",
			source:   "# Hello World\n\nThis is a test.",
			contains: []string{"<h1", "Hello World", "<p>", "This is a test"},
		},
		{
			name:     "bold and italic",
			source:   "**bold** and *italic*",
			contains: []string{"<strong>", "bold", "<em>", "italic"},
		},
		{
			name:     "code block",
			source:   "```go\nfunc main() {}\n```",
			contains: []string{"<pre>", "<code", "func main"},
		},
		{
			name:     "link",
			source:   "[Google](https://google.com)",
			contains: []string{"<a", "href=\"https://google.com\"", "Google"},
		},
		{
			name:     "list",
			source:   "- Item 1\n- Item 2\n- Item 3",
			contains: []string{"<ul>", "<li>", "Item 1", "Item 2", "Item 3"},
		},
		{
			name:     "table",
			source:   "| Header 1 | Header 2 |\n|----------|----------|\n| Cell 1   | Cell 2   |",
			contains: []string{"<table>", "<th>", "Header 1", "Header 2", "<td>", "Cell 1", "Cell 2"},
		},
		{
			name:     "empty content",
			source:   "",
			contains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mdRenderer.Render(tt.source)
			require.NoError(t, err)

			for _, expected := range tt.contains {
				assert.Contains(t, result, expected)
			}
		})
	}
}

func TestMarkdownRenderer_RenderWithMermaid(t *testing.T) {
	mermaidRenderer := renderer.NewMockMermaidRenderer()
	mdRenderer := renderer.NewMarkdownRenderer(mermaidRenderer)

	source := "# Diagram\n\n```mermaid\ngraph TD\n  A-->B\n```\n\nEnd."

	result, err := mdRenderer.Render(source)
	require.NoError(t, err)

	// Mermaidがモックで変換されているか確認
	assert.Contains(t, result, "Mermaid diagram")
	assert.Contains(t, result, "<h1")
	assert.Contains(t, result, "Diagram")
	assert.Contains(t, result, "<p>")
	assert.Contains(t, result, "End")
}
