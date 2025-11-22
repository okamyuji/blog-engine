package auth_test

import (
	"testing"
	"time"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/infrastructure/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	cfg := auth.JWTConfig{
		SecretKey:     "test-secret-key-min-32-chars-long",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 168 * time.Hour,
	}

	manager, err := auth.NewJWTManager(cfg)
	require.NoError(t, err)

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Role:     entity.RoleAdmin,
	}

	token, err := manager.GenerateAccessToken(user)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	cfg := auth.JWTConfig{
		SecretKey:     "test-secret-key-min-32-chars-long",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 168 * time.Hour,
	}

	manager, err := auth.NewJWTManager(cfg)
	require.NoError(t, err)

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Role:     entity.RoleAdmin,
	}

	token, err := manager.GenerateRefreshToken(user)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTManager_ValidateToken(t *testing.T) {
	cfg := auth.JWTConfig{
		SecretKey:     "test-secret-key-min-32-chars-long",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 168 * time.Hour,
	}

	manager, err := auth.NewJWTManager(cfg)
	require.NoError(t, err)

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Role:     entity.RoleEditor,
	}

	token, err := manager.GenerateAccessToken(user)
	require.NoError(t, err)

	// トークン検証成功
	claims, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, user.Role, claims.Role)
}

func TestJWTManager_ValidateToken_InvalidToken(t *testing.T) {
	cfg := auth.JWTConfig{
		SecretKey:     "test-secret-key-min-32-chars-long",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 168 * time.Hour,
	}

	manager, err := auth.NewJWTManager(cfg)
	require.NoError(t, err)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "invalid token format",
			token: "invalid.token.format",
		},
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "random string",
			token: "randomstring",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manager.ValidateToken(tt.token)
			assert.Error(t, err)
		})
	}
}

func TestJWTManager_GetJTI(t *testing.T) {
	cfg := auth.JWTConfig{
		SecretKey:     "test-secret-key-min-32-chars-long",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 168 * time.Hour,
	}

	manager, err := auth.NewJWTManager(cfg)
	require.NoError(t, err)

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Role:     entity.RoleAdmin,
	}

	token, err := manager.GenerateAccessToken(user)
	require.NoError(t, err)

	jti, err := manager.GetJTI(token)
	require.NoError(t, err)
	assert.NotEmpty(t, jti)
}

func TestNewJWTManager_InvalidConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  auth.JWTConfig
	}{
		{
			name: "empty secret key",
			cfg: auth.JWTConfig{
				SecretKey:     "",
				AccessExpiry:  15 * time.Minute,
				RefreshExpiry: 168 * time.Hour,
			},
		},
		{
			name: "short secret key",
			cfg: auth.JWTConfig{
				SecretKey:     "short",
				AccessExpiry:  15 * time.Minute,
				RefreshExpiry: 168 * time.Hour,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := auth.NewJWTManager(tt.cfg)
			assert.Error(t, err)
		})
	}
}
