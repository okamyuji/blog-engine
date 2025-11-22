package renderer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// MermaidRenderer Mermaidレンダラーインターフェース
type MermaidRenderer interface {
	RenderToSVG(mermaidCode string) (string, error)
}

// mermaidCLIRenderermermaid-cliを使ったMermaidRendererの実装
type mermaidCLIRenderer struct {
	tmpDir string
}

// NewMermaidRenderer 新しいMermaidRendererを作成
func NewMermaidRenderer() MermaidRenderer {
	tmpDir := os.TempDir()
	return &mermaidCLIRenderer{
		tmpDir: tmpDir,
	}
}

// RenderToSVG MermaidコードをSVGにレンダリング
func (r *mermaidCLIRenderer) RenderToSVG(mermaidCode string) (string, error) {
	if mermaidCode == "" {
		return "", fmt.Errorf("mermaid code cannot be empty")
	}

	// 一時ファイルを作成
	inputFile := filepath.Join(r.tmpDir, fmt.Sprintf("mermaid-%d.mmd", os.Getpid()))
	outputFile := filepath.Join(r.tmpDir, fmt.Sprintf("mermaid-%d.svg", os.Getpid()))

	// 入力ファイルにMermaidコードを書き込み
	if err := os.WriteFile(inputFile, []byte(mermaidCode), 0644); err != nil {
		return "", fmt.Errorf("failed to write input file: %w", err)
	}
	defer func() {
		_ = os.Remove(inputFile)
	}()

	// mmdc コマンド実行
	cmd := exec.Command("mmdc", "-i", inputFile, "-o", outputFile, "-b", "transparent")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to execute mmdc: %w, stderr: %s", err, stderr.String())
	}
	defer func() {
		_ = os.Remove(outputFile)
	}()

	// SVGファイルを読み込み
	svg, err := os.ReadFile(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read output file: %w", err)
	}

	return string(svg), nil
}

// mockMermaidRenderer テスト用のモックMermaidRenderer
type mockMermaidRenderer struct{}

// NewMockMermaidRenderer 新しいモックMermaidRendererを作成
func NewMockMermaidRenderer() MermaidRenderer {
	return &mockMermaidRenderer{}
}

// RenderToSVG モックのSVGを返す
func (r *mockMermaidRenderer) RenderToSVG(mermaidCode string) (string, error) {
	// テスト用の簡単なSVGを返す
	return `<svg><text>Mermaid diagram</text></svg>`, nil
}
