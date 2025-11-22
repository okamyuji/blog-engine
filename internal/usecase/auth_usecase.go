package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/domain/repository"
	"my-blog-engine/internal/infrastructure/auth"
)

// AuthUseCase 認証ユースケースのインターフェース
type AuthUseCase interface {
	Login(ctx context.Context, username, password string) (*LoginResponse, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, token string) (*entity.User, error)
}

// LoginResponse ログインレスポンス
type LoginResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresIn    int64        `json:"expiresIn"`
	User         *entity.User `json:"user"`
}

// RefreshTokenResponse リフレッシュトークンレスポンス
type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

// authUseCase AuthUseCaseの実装
type authUseCase struct {
	userRepo       repository.UserRepository
	tokenRepo      repository.TokenRepository
	jwtManager     auth.JWTManager
	passwordHasher auth.PasswordHasher
	accessExpiry   time.Duration
}

// NewAuthUseCase 新しいAuthUseCaseを作成
func NewAuthUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtManager auth.JWTManager,
	passwordHasher auth.PasswordHasher,
	accessExpiry time.Duration,
) AuthUseCase {
	return &authUseCase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		jwtManager:     jwtManager,
		passwordHasher: passwordHasher,
		accessExpiry:   accessExpiry,
	}
}

// Login ユーザーログイン
func (u *authUseCase) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	// ユーザー検索
	user, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// ステータスチェック
	if !user.IsActive() {
		return nil, fmt.Errorf("user account is inactive")
	}

	// パスワード検証
	if err := u.passwordHasher.Verify(user.PasswordHash, password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// トークン生成
	accessToken, err := u.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := u.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// パスワードハッシュをレスポンスから除外
	user.PasswordHash = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(u.accessExpiry.Seconds()),
		User:         user,
	}, nil
}

// Logout ユーザーログアウト
func (u *authUseCase) Logout(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}

	// トークン検証
	claims, err := u.jwtManager.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// ブラックリストに追加
	if err := u.tokenRepo.Add(ctx, claims.ID, claims.ExpiresAt.Time); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

// RefreshToken トークンをリフレッシュ
func (u *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	// トークン検証
	claims, err := u.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// ブラックリストチェック
	exists, err := u.tokenRepo.Exists(ctx, claims.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("token has been revoked")
	}

	// ユーザー検索
	user, err := u.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// ステータスチェック
	if !user.IsActive() {
		return nil, fmt.Errorf("user account is inactive")
	}

	// 新しいアクセストークン生成
	accessToken, err := u.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(u.accessExpiry.Seconds()),
	}, nil
}

// ValidateToken トークンを検証してユーザーを返す
func (u *authUseCase) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// トークン検証
	claims, err := u.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// ブラックリストチェック
	exists, err := u.tokenRepo.Exists(ctx, claims.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("token has been revoked")
	}

	// ユーザー検索
	user, err := u.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// ステータスチェック
	if !user.IsActive() {
		return nil, fmt.Errorf("user account is inactive")
	}

	// パスワードハッシュを除外
	user.PasswordHash = ""

	return user, nil
}
