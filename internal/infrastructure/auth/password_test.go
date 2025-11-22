package auth_test

import (
	"testing"

	"my-blog-engine/internal/infrastructure/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHasher_Hash(t *testing.T) {
	hasher := auth.NewPasswordHasher()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "long password",
			password: "this-is-a-very-long-password-with-many-characters-to-test-bcrypt-limits",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := hasher.Hash(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash)
			}
		})
	}
}

func TestPasswordHasher_Verify(t *testing.T) {
	hasher := auth.NewPasswordHasher()

	password := "password123"
	hash, err := hasher.Hash(password)
	require.NoError(t, err)

	tests := []struct {
		name     string
		hash     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			hash:     hash,
			password: password,
			wantErr:  false,
		},
		{
			name:     "invalid password",
			hash:     hash,
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "empty password",
			hash:     hash,
			password: "",
			wantErr:  true,
		},
		{
			name:     "empty hash",
			hash:     "",
			password: password,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.Verify(tt.hash, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordHasher_HashAndVerify(t *testing.T) {
	hasher := auth.NewPasswordHasher()

	password := "my-secure-password"

	// ハッシュ化
	hash, err := hasher.Hash(password)
	require.NoError(t, err)

	// 検証成功
	err = hasher.Verify(hash, password)
	assert.NoError(t, err)

	// 検証失敗
	err = hasher.Verify(hash, "wrong-password")
	assert.Error(t, err)
}
