package config

import (
	"log"
	"os"
	"time"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Pgxpool postgres connection pool config values
func DatabaseConfig() *pgxpool.Config {
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL %v\n", err)
	}

	config.MaxConns = 50
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = time.Minute

	return config
}