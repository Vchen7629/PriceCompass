package db

import "github.com/jackc/pgx/v4/pgxpool"

type Repository struct {
	pool *pgxpool.Pool
}

// allows us to initialize the db pool in main.go
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}