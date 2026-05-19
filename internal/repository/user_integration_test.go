package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"weather-app/internal/repository"
)

func TestUserRepo_CreateAndFetch_Integration(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	ctx := context.Background()
	db, err := pgxpool.New(ctx, url)
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, db.Ping(ctx))

	repo := repository.NewUserRepo(db)
	email := fmt.Sprintf("itest-%d@example.com", time.Now().UnixNano())

	// сохранение данных
	created, err := repo.CreateWithPassword(ctx, "ITest", email, "hash123")
	require.NoError(t, err)
	require.NotZero(t, created.ID)
	assert.Equal(t, email, created.Email)

	// чтение данных + работа SQL
	got, hash, err := repo.GetByEmail(ctx, email)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "hash123", hash)

	byID, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, email, byID.Email)

	// cleanup
	_ = repo.SoftDelete(ctx, created.ID)
}
