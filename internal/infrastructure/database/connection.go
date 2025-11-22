package database

import (
	"database/sql"
	"fmt"
	"time"

	"my-blog-engine/internal/domain/entity"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

// Config データベース接続設定
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewConnection 新しいデータベース接続を作成
func NewConnection(cfg Config) (*bun.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続プール設定
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetConnMaxLifetime(5 * time.Minute)
	sqldb.SetConnMaxIdleTime(5 * time.Minute)

	// 接続テスト
	if err := sqldb.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := bun.NewDB(sqldb, mysqldialect.New())

	// BUNモデル登録(多対多関係)
	db.RegisterModel((*entity.PostTag)(nil))

	return db, nil
}

// Close データベース接続をクローズ
func Close(db *bun.DB) error {
	return db.Close()
}
