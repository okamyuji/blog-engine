package testhelper

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"my-blog-engine/internal/domain/entity"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

var (
	sharedContainer *mysql.MySQLContainer
	sharedDB        *bun.DB
	containerMutex  sync.Mutex
	initOnce        sync.Once
)

// SetupTestDB テスト用のデータベースをセットアップ
func SetupTestDB(t *testing.T) (*bun.DB, func()) {
	t.Helper()

	initOnce.Do(func() {
		if err := initSharedContainer(); err != nil {
			t.Fatalf("Failed to initialize shared container: %v", err)
		}
	})

	// テストごとにトランザクションを使用してデータを分離
	ctx := context.Background()

	// テーブルをクリーンアップ
	cleanupTables(ctx, t, sharedDB)

	cleanup := func() {
		// テスト後のクリーンアップ
		cleanupTables(ctx, t, sharedDB)
	}

	return sharedDB, cleanup
}

// initSharedContainer 共有コンテナを初期化
func initSharedContainer() error {
	containerMutex.Lock()
	defer containerMutex.Unlock()

	if sharedContainer != nil {
		return nil
	}

	ctx := context.Background()

	// MySQLコンテナを起動
	container, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
	)
	if err != nil {
		return fmt.Errorf("failed to start MySQL container: %w", err)
	}

	sharedContainer = container

	// 接続情報を取得
	host, err := container.Host(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		return fmt.Errorf("failed to get container port: %w", err)
	}

	// DSN作成
	dsn := fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb?parseTime=true&loc=UTC&multiStatements=true",
		host, port.Port())

	// 接続
	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	db := bun.NewDB(sqldb, mysqldialect.New())

	// BUNモデル登録
	db.RegisterModel((*entity.PostTag)(nil))

	// 接続待機
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		if i == 29 {
			return fmt.Errorf("failed to connect to database after 30 attempts")
		}
		time.Sleep(time.Second)
	}

	// マイグレーション実行
	if err := runMigrationsFromFile(ctx, db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	sharedDB = db

	return nil
}

// runMigrationsFromFile SQLファイルからマイグレーションを実行
func runMigrationsFromFile(ctx context.Context, db *bun.DB) error {
	// プロジェクトルートからの相対パス
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	migrationFile := filepath.Join(projectRoot, "migrations", "001_initial_schema.up.sql")

	// SQLファイルを読み込み
	sqlContent, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// マイグレーション実行
	_, err = db.ExecContext(ctx, string(sqlContent))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

// findProjectRoot プロジェクトルートディレクトリを検索
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// go.modの存在をチェック
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found")
		}
		dir = parent
	}
}

// cleanupTables テスト用にテーブルをクリーンアップ
func cleanupTables(ctx context.Context, t *testing.T, db *bun.DB) {
	t.Helper()

	// 外部キー制約を一時的に無効化
	_, err := db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		t.Logf("Warning: failed to disable foreign key checks: %v", err)
	}

	// 各テーブルをトランケート
	tables := []string{
		"post_tags",
		"posts",
		"tags",
		"categories",
		"token_blacklist",
		"users",
	}

	for _, table := range tables {
		_, err := db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s", table))
		if err != nil {
			t.Logf("Warning: failed to truncate table %s: %v", table, err)
		}
	}

	// 外部キー制約を再度有効化
	_, err = db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		t.Logf("Warning: failed to enable foreign key checks: %v", err)
	}
}

// CleanupSharedContainer 共有コンテナをクリーンアップ（テスト全体の終了時に呼び出す）
func CleanupSharedContainer() error {
	containerMutex.Lock()
	defer containerMutex.Unlock()

	if sharedDB != nil {
		if err := sharedDB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		sharedDB = nil
	}

	if sharedContainer != nil {
		ctx := context.Background()
		if err := sharedContainer.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate container: %w", err)
		}
		sharedContainer = nil
	}

	return nil
}
