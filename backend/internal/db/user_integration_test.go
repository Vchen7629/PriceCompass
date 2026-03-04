//go:build integration

package db_test

import (
	"backend/internal/db"
	"backend/pkg/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Integration tests for InsertProductForUser SQL func
func TestInsertNewUser(t *testing.T) {
	pool := testDB.Pool
	repo := db.NewRepository(pool)

	t.Run("Trying to insert a new username that already exists", func(t *testing.T) {
		test.CleanupTables(t, pool)
	
		test.SeedUser(t, pool, "username1", "idk@gmail.com")
		err := repo.InsertNewUser("username1", "123@gmail.com", "password")

		assert.Equal(t, "username already exists", err.Error(), "should return error")
	})

	t.Run("Trying to insert a new email that already exists", func(t *testing.T) {
		test.CleanupTables(t, pool)

		test.SeedUser(t, pool, "username1", "idk@gmail.com")
		err := repo.InsertNewUser("username2", "idk@gmail.com", "password")

		assert.Equal(t, "email already exists", err.Error(), "should return error")
	})

	t.Run("Trying to insert new user should succeed", func(t *testing.T) {
		test.CleanupTables(t, pool)

		test.SeedUser(t, pool, "username2", "idk@gmail.com")
		err := repo.InsertNewUser("username1", "123@gmail.com", "password")

		assert.Equal(t, nil, err, "should return nil or no error")
	})

	t.Run("Inserting actually stores the data in db", func(t *testing.T) {
		test.CleanupTables(t, pool)
		var exists bool

		err := repo.InsertNewUser("username1", "123@gmail.com", "password")

		assert.Equal(t, nil, err, "should return nil or no error")

		userExistQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
		err = pool.QueryRow(context.Background(), userExistQuery, "username1").Scan(&exists)
		assert.Equal(t, nil, err, "querying should not error")

		assert.Equal(t, true, exists, "Querying should be true meaning the username is in the database")
	})

	t.Run("Checking that password is not stored as plaintext", func(t *testing.T) {
		test.CleanupTables(t, pool)
		var password string

		err := repo.InsertNewUser("username1", "123@gmail.com", "password")

		assert.Equal(t, nil, err, "should return nil or no error")

		userExistQuery := `SELECT password_hash FROM users WHERE username = $1`
		err = pool.QueryRow(context.Background(), userExistQuery, "username1").Scan(&password)
		assert.Equal(t, nil, err, "querying should not error")

		assert.NotEqual(t, "password", password, "Password fetched from the database should not be plaintext")
	})
}

// Integration tests for LoginUser SQL func
func TestLoginUser(t *testing.T) {
	pool := testDB.Pool
	repo := db.NewRepository(pool)

	t.Run("Invalid password should raise an error", func(t *testing.T) {
		test.CleanupTables(t, pool)
		repo.InsertNewUser("user1", "123@gmail.com", "password123")

		sessionToken, err := repo.LoginUser("user1", "invalidpass")

		assert.NotEqual(t, nil, err, "should return an error")
		assert.Equal(t, "", sessionToken, "session token should be empty")
	})

	t.Run("Invalid username should raise an error", func(t *testing.T) {
		test.CleanupTables(t, pool)
		repo.InsertNewUser("user1", "123@gmail.com", "password123")

		sessionToken, err := repo.LoginUser("user123", "password123")

		assert.NotEqual(t, nil, err, "should return an error")
		assert.Equal(t, "", sessionToken, "session token should be empty")
	})

	t.Run("Correct username and password should create a new session in database and return it", func(t *testing.T) {
		test.CleanupTables(t, pool)
		var sessionExists bool
		repo.InsertNewUser("user1", "123@gmail.com", "password123")

		sessionToken, err := repo.LoginUser("user1", "password123")

		assert.Equal(t, nil, err, "should not return an error")

		sessionExistQuery := `SELECT EXISTS(SELECT 1 FROM sessions WHERE username = $1)`
		err = pool.QueryRow(context.Background(), sessionExistQuery, "user1").Scan(&sessionExists)

		assert.Equal(t, true, sessionExists, "Session token should exist for the user in the table")
		assert.NotEqual(t, "", sessionToken, "session token should not be empty")
	})
}
