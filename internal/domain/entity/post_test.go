package entity_test

import (
	"testing"
	"time"

	"my-blog-engine/internal/domain/entity"

	"github.com/stretchr/testify/assert"
)

func TestPost_IsPublished(t *testing.T) {
	tests := []struct {
		name     string
		status   entity.PostStatus
		expected bool
	}{
		{
			name:     "published post",
			status:   entity.StatusPublished,
			expected: true,
		},
		{
			name:     "draft post",
			status:   entity.StatusDraft,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &entity.Post{Status: tt.status}
			assert.Equal(t, tt.expected, post.IsPublished())
		})
	}
}

func TestPost_Publish(t *testing.T) {
	post := &entity.Post{Status: entity.StatusDraft}

	// 初回公開
	post.Publish()
	assert.Equal(t, entity.StatusPublished, post.Status)
	assert.NotNil(t, post.PublishedAt)

	firstPublishedAt := post.PublishedAt

	// 2回目の公開（PublishedAtは変更されない）
	time.Sleep(10 * time.Millisecond)
	post.Publish()
	assert.Equal(t, entity.StatusPublished, post.Status)
	assert.Equal(t, firstPublishedAt, post.PublishedAt)
}

func TestPost_Unpublish(t *testing.T) {
	now := time.Now()
	post := &entity.Post{
		Status:      entity.StatusPublished,
		PublishedAt: &now,
	}

	post.Unpublish()
	assert.Equal(t, entity.StatusDraft, post.Status)
	// PublishedAtはそのまま残る
	assert.NotNil(t, post.PublishedAt)
}
