package auth

import (
	"fmt"
	"time"

	"my-blog-engine/internal/domain/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims JWTクレーム
type Claims struct {
	UserID   int64           `json:"sub"`
	Username string          `json:"username"`
	Role     entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// JWTManagerJWT 管理インターフェース
type JWTManager interface {
	GenerateAccessToken(user *entity.User) (string, error)
	GenerateRefreshToken(user *entity.User) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	GetJTI(tokenString string) (string, error)
}

// jwtManager JWTManagerの実装
type jwtManager struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	signingMethod jwt.SigningMethod
}

// JWTConfigJWT 設定
type JWTConfig struct {
	SecretKey     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

// NewJWTManager 新しいJWTManagerを作成
func NewJWTManager(cfg JWTConfig) (JWTManager, error) {
	if cfg.SecretKey == "" {
		return nil, fmt.Errorf("secret key cannot be empty")
	}

	if len(cfg.SecretKey) < 32 {
		return nil, fmt.Errorf("secret key must be at least 32 characters")
	}

	return &jwtManager{
		secretKey:     []byte(cfg.SecretKey),
		accessExpiry:  cfg.AccessExpiry,
		refreshExpiry: cfg.RefreshExpiry,
		signingMethod: jwt.SigningMethodHS256,
	}, nil
}

// GenerateAccessToken アクセストークンを生成
func (m *jwtManager) GenerateAccessToken(user *entity.User) (string, error) {
	return m.generateToken(user, m.accessExpiry)
}

// GenerateRefreshToken リフレッシュトークンを生成
func (m *jwtManager) GenerateRefreshToken(user *entity.User) (string, error) {
	return m.generateToken(user, m.refreshExpiry)
}

// generateToken トークンを生成
func (m *jwtManager) generateToken(user *entity.User, expiry time.Duration) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user cannot be nil")
	}

	now := time.Now()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(m.signingMethod, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken トークンを検証
func (m *jwtManager) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名方法の検証
		if token.Method != m.signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GetJTI トークンからJTIを取得
func (m *jwtManager) GetJTI(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	if claims.ID == "" {
		return "", fmt.Errorf("token does not contain JTI")
	}

	return claims.ID, nil
}
