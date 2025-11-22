package entity_test

import (
	"testing"

	"my-blog-engine/internal/domain/entity"

	"github.com/stretchr/testify/assert"
)

func TestUser_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   entity.UserStatus
		expected bool
	}{
		{
			name:     "active user",
			status:   entity.StatusActive,
			expected: true,
		},
		{
			name:     "inactive user",
			status:   entity.StatusInactive,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &entity.User{Status: tt.status}
			assert.Equal(t, tt.expected, user.IsActive())
		})
	}
}

func TestUser_HasRole(t *testing.T) {
	user := &entity.User{Role: entity.RoleEditor}

	assert.True(t, user.HasRole(entity.RoleEditor))
	assert.False(t, user.HasRole(entity.RoleAdmin))
	assert.False(t, user.HasRole(entity.RoleViewer))
}

func TestUser_CanEdit(t *testing.T) {
	tests := []struct {
		name     string
		role     entity.UserRole
		expected bool
	}{
		{
			name:     "admin can edit",
			role:     entity.RoleAdmin,
			expected: true,
		},
		{
			name:     "editor can edit",
			role:     entity.RoleEditor,
			expected: true,
		},
		{
			name:     "viewer cannot edit",
			role:     entity.RoleViewer,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &entity.User{Role: tt.role}
			assert.Equal(t, tt.expected, user.CanEdit())
		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     entity.UserRole
		expected bool
	}{
		{
			name:     "admin user",
			role:     entity.RoleAdmin,
			expected: true,
		},
		{
			name:     "editor user",
			role:     entity.RoleEditor,
			expected: false,
		},
		{
			name:     "viewer user",
			role:     entity.RoleViewer,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &entity.User{Role: tt.role}
			assert.Equal(t, tt.expected, user.IsAdmin())
		})
	}
}
