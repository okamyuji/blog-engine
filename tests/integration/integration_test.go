package integration

import (
	"os"
	"testing"

	"my-blog-engine/tests/integration/testhelper"
)

// TestMain テスト全体のセットアップとクリーンアップ
func TestMain(m *testing.M) {
	// テスト実行
	code := m.Run()

	// 共有コンテナのクリーンアップ
	if err := testhelper.CleanupSharedContainer(); err != nil {
		// エラーログを出力するが、テスト結果には影響させない
		println("Warning: failed to cleanup shared container:", err.Error())
	}

	os.Exit(code)
}
