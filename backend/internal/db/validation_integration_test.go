//go:build integration

package db_test

import (
	"backend/internal/db"
	"backend/internal/test"
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)



func TestValidateExistsUserTable(t *testing.T) {
	pool := testDB.Pool
	ctx := context.Background()

	t.Run("Should not return an error if the value in column doesnt exist", func(t *testing.T) {
		test.CleanupTables(t, pool)
		
		test.SeedUser(t, pool, "username1", "idk@gmail.com")

		_ = db.WithTransaction(ctx, pool, func(tx pgx.Tx) error {
			err := db.ValidateExistsUserTable(ctx, tx, "username", "username10")

			assert.Equal(t, nil, err, "Should not return an error if value doesnt exist")

			return nil
		})
	})

	t.Run("Should return an error if the value in column exists", func(t *testing.T) {
		test.CleanupTables(t, pool)
		
		test.SeedUser(t, pool, "username1", "idk@gmail.com")

		_ = db.WithTransaction(ctx, pool, func(tx pgx.Tx) error {
			err := db.ValidateExistsUserTable(ctx, tx, "username", "username1")

			assert.NotEqual(t, nil, err, "Should return an error if value exists")

			return nil
		})
	})
}

func TestValidateSession(t *testing.T) {
	pool := testDB.Pool
	repo := db.NewRepository(pool)

	t.Run("Should return an error if nonexistant session token is provided", func(t *testing.T) {
		test.CleanupTables(t, pool)

		repo.InsertNewUser("user1", "idk@gmail.com", "pass123")
		_, err := repo.LoginUser("user1", "pass123")

		assert.Equal(t, nil, err)

		userId, username, err := repo.ValidateSession("23udweu")

		assert.NotEqual(t, nil, err, "Should return a non nil error")
		assert.Equal(t, 0, userId, "validate session returns 0 userId on error")
		assert.Equal(t, "", username, "returns an empty string on error")
	})

	t.Run("Should return correct userId and username if correct session token is provided", func(t *testing.T) {
		test.CleanupTables(t, pool)

		repo.InsertNewUser("user1", "idk@gmail.com", "pass123")
		sessionToken, err := repo.LoginUser("user1", "pass123")

		assert.Equal(t, nil, err)

		userId, username, err := repo.ValidateSession(sessionToken)

		assert.Equal(t, nil, err, "Should return a nil error")
		assert.Equal(t, 1, userId)
		assert.Equal(t, "user1", username)
	})
}

