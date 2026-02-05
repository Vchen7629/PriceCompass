package db

import (
	"backend/pkg/db"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Create a new user account in the database with the provided
// username and password strings
func InsertNewUser(username, password string, pool *pgxpool.Pool) error {
	ctx := context.Background()
	createdAt := time.Now()
	
	err := db.WithTransaction(ctx, pool, func(tx pgx.Tx) error {
		var exists bool
		userExistQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

		// need to update all my other queries to use tx instead of pool for atomic transactions
		err := tx.QueryRow(ctx, userExistQuery, username).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("user already exists")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}

		userQuery := `
			INSERT INTO users (username, password_hash, created_at)
			VALUES ($1, $2, $3)
		`

		_, err = tx.Exec(ctx, userQuery, username, hashedPassword, createdAt)
		if err != nil {
			return fmt.Errorf("failed to create new user: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// attempt to login to an user account using the username and password provided
// return the session token if successful
func LoginUser(username, password string, pool *pgxpool.Pool) (string, error) {
	ctx := context.Background()
	var sessionToken string

	err := db.WithTransaction(ctx, pool, func(tx pgx.Tx) error {
		var hashedPassword string

		fetchPassQuery := `SELECT password FROM users WHERE username = $1`

		err := tx.QueryRow(ctx, fetchPassQuery, username).Scan(&hashedPassword)
		if err != nil {
			return fmt.Errorf("error fetching password: %w", err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			return fmt.Errorf("Passwords don't match, can't login")
		}

		// generating session token using random crypto
		tokenBytes := make([]byte, 32)
		_, err = rand.Read(tokenBytes)
		if err != nil {
			return fmt.Errorf("error generating token: %w", err)
		}

		sessionToken = hex.EncodeToString(tokenBytes)

		insertSessionToken := `
			INSERT INTO sessions (username, token, expires_at)
			VALUES ($1, $2, $3)
		`
		expiresAt := time.Now().Add(time.Hour)
		_, err = tx.Exec(ctx, insertSessionToken, username, sessionToken, expiresAt)
		if err != nil {
			return fmt.Errorf("error creating session: %w", err)
		}

		return nil
	})

	return sessionToken, err
}