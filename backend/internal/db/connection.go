package db

import (
	"context"
	"log"
	"github.com/jackc/pgx/v4/pgxpool"
	"backend/internal/config"
)

// set up pgxpool db connection pool
func ConnectionPool() *pgxpool.Pool {
	pool, err := pgxpool.ConnectConfig(context.Background(), config.DatabaseConfig())
	if err != nil {
		log.Fatalf("Unable to connect to database %v\n", err)
	}

	return pool
}