package persistence_test

import (
	"context"
	"fmt"
	"testing"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/infrastructure/persistence"
	"my-blog-engine/tests/integration/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         entity.RoleEditor,
		Status:       entity.StatusActive,
	}

	err := repo.Create(ctx, user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_FindByID(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewUserRepository(db)
	ctx := context.Background()

	// テストデータ作成
	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         entity.RoleEditor,
		Status:       entity.StatusActive,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// 検索
	found, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewUserRepository(db)
	ctx := context.Background()

	// テストデータ作成
	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         entity.RoleEditor,
		Status:       entity.StatusActive,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// 検索
	found, err := repo.FindByUsername(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepository_Update(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewUserRepository(db)
	ctx := context.Background()

	// テストデータ作成
	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         entity.RoleEditor,
		Status:       entity.StatusActive,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// 更新
	user.Email = "updated@example.com"
	err = repo.Update(ctx, user)
	require.NoError(t, err)

	// 検証
	found, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated@example.com", found.Email)
}

func TestUserRepository_Delete(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewUserRepository(db)
	ctx := context.Background()

	// テストデータ作成
	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         entity.RoleEditor,
		Status:       entity.StatusActive,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// 削除
	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	// 検証
	_, err = repo.FindByID(ctx, user.ID)
	assert.Error(t, err)
}

func TestUserRepository_List(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewUserRepository(db)
	ctx := context.Background()

	// テストデータ作成
	for i := 1; i <= 3; i++ {
		user := &entity.User{
			Username:     fmt.Sprintf("testuser%d", i),
			Email:        fmt.Sprintf("test%d@example.com", i),
			PasswordHash: "hashedpassword",
			Role:         entity.RoleEditor,
			Status:       entity.StatusActive,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	// 一覧取得
	users, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 3)
}

func TestUserRepository_Count(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	repo := persistence.NewUserRepository(db)
	ctx := context.Background()

	// テストデータ作成
	for i := 1; i <= 5; i++ {
		user := &entity.User{
			Username:     fmt.Sprintf("testuser%d", i),
			Email:        fmt.Sprintf("test%d@example.com", i),
			PasswordHash: "hashedpassword",
			Role:         entity.RoleEditor,
			Status:       entity.StatusActive,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	// カウント
	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 5)
}
