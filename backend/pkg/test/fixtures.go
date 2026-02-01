package test

import (
	"context"
	"testing"
	"time"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductSourceConfig struct {                                                                        
	ProductID         int                                                                                  
	Platform          string                                                                               
	PlatformProductID string                                                                               
	URL               string                                                                               
} 

type PriceSnapshotConfig struct {
	ProductSourceID int
	Price			float64
	Currency 		string
	InStock			bool
	CheckedAt		*time.Time
}

// Create a test user and return userID in database
func SeedUser(t *testing.T, pool *pgxpool.Pool, email string) int {
	t.Helper()

	ctx := context.Background()
	var userId int

	err := pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, created_at)
		 VALUES ($1, $2, NOW())
		 RETURNING id`,
		 email, "hashed_password_placeholder",
	).Scan(&userId)

	if err != nil {
		t.Fatalf("Failed to seed user: %v", err)
	}

	return userId
}

// Create a product and return the product ID
func SeedProduct(t *testing.T, pool *pgxpool.Pool, name, imageURL string) int {
	t.Helper()

	ctx := context.Background()
	var productId int

	err := pool.QueryRow(ctx,
		`INSERT INTO products (product_name, image_url, created_at, last_checked_at, check_priority)
		 VALUES ($1, $2, NOW(), NOW(), 0)
		 RETURNING id`,
		 name, imageURL,
	).Scan(&productId)

	if err != nil {
		t.Fatalf("Failed to seed product: %v", err)
	}

	return productId
}

// Add a product to user's watchlist
func AddProductToWatchlist(t *testing.T, pool *pgxpool.Pool, userId, productID int) {
	t.Helper()

	ctx := context.Background()

	_, err := pool.Exec(ctx,
		`INSERT INTO user_watchlist (user_id, product_id, added_at)
		 VALUES ($1, $2, NOW())`,
		 userId, productID,
	)

	if err != nil {
		t.Fatalf("Failed to add product to watchlist: %v", err)
	}
}

// creates a product source and returns the sourceId
func SeedProductSource(t *testing.T, pool *pgxpool.Pool, config ProductSourceConfig) int {
	t.Helper()

	ctx := context.Background()
	var sourceID int

	err := pool.QueryRow(ctx,
		`INSERT INTO product_sources (product_id, platform, platform_product_id, product_url)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		 config.ProductID, config.Platform, config.PlatformProductID, config.URL,
	).Scan(&sourceID)

	if err != nil {
		t.Fatalf("Failed to seed product source: %v", err)
	}

	return sourceID
}

// creates a price snapshot
func SeedPriceSnapshot(t *testing.T, pool *pgxpool.Pool, config PriceSnapshotConfig) {
	t.Helper()

	ctx := context.Background()

	// default to current time if not specified
	checkedAt := time.Now()
	if config.CheckedAt != nil {
		checkedAt = *config.CheckedAt
	}

	// Default currency if not specified
	currency := "USD"
	if config.Currency != "" {
		currency = config.Currency
	}

	_, err := pool.Exec(ctx,
		`INSERT INTO price_snapshots (product_source_id, price, currency, in_stock, checked_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		 config.ProductSourceID, config.Price, currency, config.InStock, checkedAt,
	)

	if err != nil {
		t.Fatalf("Failed to seed product snapshot: %v", err)
	}
}