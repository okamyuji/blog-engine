package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

// PasswordHasher パスワードハッシュ化インターフェース
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) error
}

// bcryptHasher bcryptを使ったPasswordHasherの実装
type bcryptHasher struct{}

// NewPasswordHasher 新しいPasswordHasherを作成
func NewPasswordHasher() PasswordHasher {
	return &bcryptHasher{}
}

// Hash パスワードをハッシュ化
func (h *bcryptHasher) Hash(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

// Verify パスワードを検証
func (h *bcryptHasher) Verify(hashedPassword, password string) error {
	if hashedPassword == "" || password == "" {
		return fmt.Errorf("password and hash cannot be empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}
