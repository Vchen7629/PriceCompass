package test

import "github.com/jackc/pgx/v4/pgxpool"

type MockPool struct{}

// mocking the database pool for unit tests
func NewMockPool() *pgxpool.Pool {
	return nil
}