package db

import (
	"context"
	"errors"
	"time"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"backend/internal/types"
	"backend/pkg/db"
)

type UserProduct struct {
	ProductID 		int 		`json:"product_id"`
	ProductName		string		`json:"product_name"`
	ImageUrl		string		`json:"image_url"`
	AddedAt			time.Time	`json:"added_at"`
	LastCheckedAt 	time.Time	`json:"last_checked_at"`
	LowestPrice		float64		`json:"lowest_price"`
	LowestSource	string 		`json:"lowest_source"`
	InStock			bool		`json:"in_stock"`
}

// Insert a product into the products table for the user with the name and timestamps
// additionally returns product metadata so frontend can immediately show the product
func InsertProductForUser(userID int, productName string, pool *pgxpool.Pool) (types.Product, error) {
	ctx := context.Background()
	var product types.Product

	createdAt := time.Now()

	// One transaction for both queries so it can rollback on errors
	err := db.WithTransaction(ctx, pool, func(tx pgx.Tx) error {
		// Try to insert, but if conflict, fetch the existing product instead
		productQuery := `
			INSERT INTO products (product_name, created_at, last_checked_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (product_name) DO UPDATE
				SET product_name = EXCLUDED.product_name
			RETURNING id, product_name, created_at`

		err := tx.QueryRow(ctx, productQuery, productName, createdAt, nil).Scan(
			&product.ID,
			&product.Name,
			&product.CreatedAt,
		)
		if err != nil {
			return err
		}

		// Add this product for the user
		userProductQuery := `
			INSERT INTO user_watchlist (user_id, product_id, added_at)
			VALUES ($1, $2, $3)`

		_, err = tx.Exec(ctx, userProductQuery, userID, product.ID, product.CreatedAt)
		return err
	})

	if err != nil {
		return types.Product{}, err
	}

	product.Prices = []types.PriceData{}

	// todo: create a helper function to trigger price scraping, call the api for that ig
	// should prob just send a job request to a job queeue and have the api for that consume from it

	return types.Product{
		ID: product.ID,
		Name: product.Name,
		URL: "idk man",
		CreatedAt: createdAt,
		Prices: []types.PriceData{},
	}, nil
}

// Fetch all tracked products for the user, returns a list of products with the
// name, added at timestamp, lowest price, and an availablity flag
func FetchUserTrackedProducts(userID int, pool *pgxpool.Pool) ([]UserProduct, error) {
	ctx := context.Background()
	var productList []UserProduct
	
	err := db.WithTransaction(ctx, pool, func(pgx.Tx) error {
		query := `
			WITH latest_prices AS (
				SELECT DISTINCT ON (pso.product_id, pso.platform)
					pso.product_id, pso.platform,
					psnap.price, psnap.in_stock, psnap.checked_at
				FROM product_sources pso
				LEFT JOIN price_snapshots psnap ON pso.id = psnap.product_source_id
				ORDER BY pso.product_id, pso.platform, psnap.checked_at DESC NULLS LAST
			),
			lowest_prices AS (
				SELECT
					product_id,
					MIN(price) FILTER (WHERE price IS NOT NULL) as lowest_price,
					(ARRAY_AGG(platform ORDER BY price ASC) FILTER (WHERE price IS NOT NULL))[1] as lowest_source,
					(ARRAY_AGG(in_stock ORDER BY price ASC) FILTER (WHERE price IS NOT NULL))[1] as in_stock
				FROM latest_prices
				GROUP BY product_id
			)
			SELECT 
				p.id, p.product_name, p.image_url, p.last_checked_at,
				uw.added_at, 
				COALESCE(lp.lowest_price, 0) as lowest_price,
				COALESCE(lp.lowest_source, '') as lowest_source,
				COALESCE(lp.in_stock, false) as in_stock
			FROM user_watchlist uw
			INNER JOIN products p ON uw.product_id = p.id
			LEFT JOIN lowest_prices lp ON p.id = lp.product_id
			WHERE uw.user_id = $1
			ORDER BY uw.added_at DESC`
	
		rows, err := pool.Query(
			context.Background(), 
			query, 
			userID,
		)

		if err != nil {
			return err
		}

		for rows.Next() {
			var product UserProduct

			err := rows.Scan(
				&product.ProductID,
				&product.ProductName,
				&product.ImageUrl,
				&product.LastCheckedAt,
				&product.AddedAt,
				&product.LowestPrice,
				&product.LowestSource,
				&product.InStock,
			)
			if err != nil {
				return err
			}

			productList = append(productList, product)
		}

		return nil
	})

	if err != nil {
		return []UserProduct{}, err
	}

	return productList, nil
}

// Delete the specified product for the user
func DeleteProductForUser(userID, productID int, pool *pgxpool.Pool) error {
	ctx := context.Background()

	err := db.WithTransaction(ctx, pool, func(tx pgx.Tx) error {
		query := `
			DELETE FROM user_watchlist
			WHERE user_id = $1
			AND product_id = $2
		`

		cmdTag, err := tx.Exec(ctx, query, userID, productID)
		if err != nil {
			return err
		}

		if cmdTag.RowsAffected() == 0 {
			return errors.New("product not found in user's watchlist")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}