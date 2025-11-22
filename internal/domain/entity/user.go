package entity

import (
	"time"

	"github.com/uptrace/bun"
)

// UserRole ユーザーの役割を表す型
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleEditor UserRole = "editor"
	RoleViewer UserRole = "viewer"
)

// UserStatus ユーザーのステータスを表す型
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
)

// User ユーザーエンティティ
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           int64      `bun:"id,pk,autoincrement"`
	Username     string     `bun:"username,unique,notnull"`
	Email        string     `bun:"email,unique,notnull"`
	PasswordHash string     `bun:"password_hash,notnull"`
	Role         UserRole   `bun:"role,notnull,default:'viewer'"`
	Status       UserStatus `bun:"status,notnull,default:'active'"`
	CreatedAt    time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

// IsActive ユーザーがアクティブかどうかを判定
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// HasRole 指定されたロールを持っているかチェック
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// CanEdit 編集権限があるかチェック
func (u *User) CanEdit() bool {
	return u.Role == RoleAdmin || u.Role == RoleEditor
}

// IsAdmin 管理者権限があるかチェック
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
