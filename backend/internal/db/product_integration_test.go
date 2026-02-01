//go:build integration

package db_test

import (
	"backend/internal/db"
	"backend/internal/types"
	"backend/pkg/test"
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *test.TestDBContainer

// Setup and teardown database for all tests in this package
func TestMain(m *testing.M) {
	var cleanup func()
	testDB, cleanup = test.SetupTestDatabaseForTestMain()
	defer cleanup()

	// Run all tests
	code := m.Run()

	os.Exit(code)
}

// Integration tests for InsertProductForUser SQL func
func TestInsertProductForUser(t *testing.T) {
	pool := testDB.Pool
	t.Run("Returns correct added product", func(t *testing.T) {
		test.CleanupTables(t, pool)

		userID := test.SeedUser(t, pool, "example@example.com")
		product, err := db.InsertProductForUser(userID, "product", pool)
		
		require.NoError(t, err)
		require.NotEmpty(t, product.ID, "product ID should be returned")
		assert.Equal(t, "product", product.Name, "return product name should match input")
		require.NotEmpty(t, product.URL, "product url should be returned")
		require.NotZero(t, product.CreatedAt, "product url should be returned")
		require.NotNil(t, product.Prices, "product url should be returned")
	})

	t.Run("adds existing product to new user's watchlist", func(t *testing.T) {
		test.CleanupTables(t, pool)

		var product types.Product
		var err error

		users := []string{"user1@example.com", "user2@example.com"}
		for _, user := range users {
			userID := test.SeedUser(t, pool, user)

			product, err = db.InsertProductForUser(userID, "product", pool)
			require.NoError(t, err)
		}


		require.NotEmpty(t, product.ID, "product ID should be returned")
		assert.Equal(t, "product", product.Name, "return product name should match input")
		require.NotEmpty(t, product.URL, "product url should be returned")
		require.NotZero(t, product.CreatedAt, "product url should be returned")
		require.NotNil(t, product.Prices, "product url should be returned")
	})


	t.Run("throws sql error when duplicate product is entered for user", func(t *testing.T) {
		test.CleanupTables(t, pool)

		userID := test.SeedUser(t, pool, "user@example.com")

		var err error
		var pgErr *pgconn.PgError

		for i := 0; i < 2; i++ {
			_, err = db.InsertProductForUser(userID, "product", pool)
		}

		require.ErrorAs(t, err, &pgErr, "Should throw error due to duplicate product for the same error")
		assert.Equal(t, "23505", pgErr.Code, "Should be unique violation code")
		assert.Contains(t, pgErr.Message, "user_watchlist") // verify correct table name
	})

	t.Run("transaction rolls back when second insert fails", func(t *testing.T) {
		test.CleanupTables(t, pool)

		nonExistantUserID := 99999
		nonExistantProduct := "nonexistant"
		insertProduct := "New Product"

		// database will have one product before transaction fails
		test.SeedProduct(t, pool, nonExistantProduct, "https://imgur.com/idk.jpg")

		// attempt to insert, will fail on second insert because the user doesn't exist
		_, err := db.InsertProductForUser(nonExistantUserID, insertProduct, pool)
		require.Error(t, err, "Should fail due to non-existant user")

		var pgErr *pgconn.PgError
		require.ErrorAs(t, err, &pgErr)
		assert.Equal(t, "23503", pgErr.Code, "Should be foreign key violation")

		// verify that transaction rolled back and product does not exist
		var countAfter int
		insertErr := pool.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM products WHERE product_name = $1", 
			insertProduct,
		).Scan(&countAfter)

		require.NoError(t, insertErr)
		assert.Equal(t, 0, countAfter, "Product should not exist after rollback")
	})
}

// Integration tests for FetchUserTrackedProducts SQL func
func TestFetchUserTrackedProducts(t *testing.T) {
	pool := testDB.Pool

	t.Run("returns empty list for user with no tracked products", func(t *testing.T) {
		test.CleanupTables(t, pool)

		userID := test.SeedUser(t, pool, "empty@example.com")

		products, err := db.FetchUserTrackedProducts(userID, pool)

		require.NoError(t, err)
		assert.Empty(t, products, "Expected no products for user with empty watchlist")
	})

	t.Run("returns single product with multiple sources and latest prices", func(t *testing.T) {
		test.CleanupTables(t, pool)

		userID := test.SeedUser(t, pool, "empty@example.com")
		productID := test.SeedProduct(t, pool, "Mechanical Keyboard", "https://example.com/keyboard.jpg")
		test.AddProductToWatchlist(t, pool, userID, productID)
		
		amazonSourceID := test.SeedProductSource(t, pool, test.ProductSourceConfig{
			ProductID: 			productID,
			Platform:  			"amazon",
			PlatformProductID: 	"BZAHSDH123",
			URL: 				"https://amazon.com/db/BZAHSDH123",
		})

		// this should not be returned because its an older price (not latest price snapshot)
		oldTime := time.Now().Add(-24 * time.Hour)
		test.SeedPriceSnapshot(t, pool, test.PriceSnapshotConfig{
			ProductSourceID: 	amazonSourceID,
			Price: 				99.99,
			InStock: 			true,
			CheckedAt: 			&oldTime,		
		})

		// since this is the latest price snapshot this should be returned
		test.SeedPriceSnapshot(t, pool, test.PriceSnapshotConfig{
			ProductSourceID: 	amazonSourceID,
			Price: 				76.99,
			InStock: 			true,
		})

		bestbuySourceID := test.SeedProductSource(t, pool, test.ProductSourceConfig{
			ProductID: 			productID,
			Platform:  			"bestbuy",
			PlatformProductID: 	"SKU456",
			URL: 				"https://bestbuy.com/products/SKU456",
		})

		// since this is the latest price snapshot this should be returned
		test.SeedPriceSnapshot(t, pool, test.PriceSnapshotConfig{
			ProductSourceID: 	bestbuySourceID,
			Price: 				39.99,
			InStock: 			true,
		})

		products, err := db.FetchUserTrackedProducts(userID, pool)

		require.NoError(t, err)
		require.Len(t, products, 1, "Expected 1 product")

		product := products[0]
		assert.Equal(t, productID, product.ProductID, "return product id should match input")
		assert.Equal(t, "Mechanical Keyboard", product.ProductName, "return product name should match input")
		assert.Equal(t, "https://example.com/keyboard.jpg", product.ImageUrl, "return product image should match input")

		assert.Equal(t, 39.99, product.LowestPrice, "Should return lowest price across all sources")
		assert.Equal(t, "bestbuy", product.LowestSource, "Best Buy has the lowest price")
		assert.True(t, product.InStock, "Lowest priced source is in stock")
	})

	t.Run("returns multiple products for same user", func(t *testing.T) {
		test.CleanupTables(t, pool)
		
		userID := test.SeedUser(t, pool, "example@example.com")

		product1ID := test.SeedProduct(t, pool, "Laptop", "https://example.com/laptop.jpg")
		test.AddProductToWatchlist(t, pool, userID, product1ID)

		product2ID := test.SeedProduct(t, pool, "Monitor", "https://example.com/monitor.jpg")
		test.AddProductToWatchlist(t, pool, userID, product2ID)

		source1ID := test.SeedProductSource(t, pool, test.ProductSourceConfig{
				ProductID: 			product1ID,
				Platform: 			"amazon",
				PlatformProductID: 	"LAPTOP123",
				URL: 				"https://amazon.com/laptop",	
		})
		source2ID := test.SeedProductSource(t, pool, test.ProductSourceConfig{
				ProductID: 			product2ID,
				Platform: 			"newegg",
				PlatformProductID: 	"MONITOR123",
				URL: 				"https://newegg.com/monitor",	
		})

		test.SeedPriceSnapshot(t, pool, test.PriceSnapshotConfig{
				ProductSourceID: source1ID,
				Price: 			 1299.99,
				InStock: 		 true,	
		})
		test.SeedPriceSnapshot(t, pool, test.PriceSnapshotConfig{
				ProductSourceID: source2ID,
				Price: 			 399.99,
				InStock: 		 false,	
		})

		products, err := db.FetchUserTrackedProducts(userID, pool)

		require.NoError(t, err)
		require.Len(t,  products, 2, "Expected 2 products")

		productMap := make(map[string]*db.UserProduct)
		for i := range products {
			productMap[products[i].ProductName] = &products[i]
		}

		laptop := productMap["Laptop"]
		require.NotNil(t, laptop)
		assert.Equal(t, 1299.99, laptop.LowestPrice)
		assert.Equal(t, "amazon", laptop.LowestSource)
		assert.True(t, laptop.InStock)

		monitor := productMap["Monitor"]
		require.NotNil(t, monitor)
		assert.Equal(t, 399.99, monitor.LowestPrice)
		assert.Equal(t, "newegg", monitor.LowestSource)
		assert.False(t, monitor.InStock)
	})

	t.Run("Handles source with no price snapshots", func(t *testing.T) {
		test.CleanupTables(t, pool)
		
		userID := test.SeedUser(t, pool, "example@example.com")
		productID := test.SeedProduct(t, pool, "New Product", "https://example.com/new.jpg")
		test.AddProductToWatchlist(t, pool, userID, productID)

		test.SeedProductSource(t, pool, test.ProductSourceConfig{
				ProductID: 			productID,
				Platform: 			"amazon",
				PlatformProductID:  "new",
				URL: 				"https://amazon.com/new",
		})

		products, err := db.FetchUserTrackedProducts(userID, pool)

		require.NoError(t, err)
		require.Len(t, products, 1)
		// Product with no price snapshots should have empty lowest source and 0 price
		assert.Empty(t, products[0].LowestSource, "Product with no snapshots should have empty lowest source")
		assert.Equal(t, 0.0, products[0].LowestPrice, "Product with no snapshots should have 0 price")
	})

	t.Run("Only returns products for specified user", func(t *testing.T) {
		test.CleanupTables(t, pool)
		
		user1ID := test.SeedUser(t, pool, "user1@example.com")
		user2ID := test.SeedUser(t, pool, "user2@example.com")
		
		product1ID := test.SeedProduct(t, pool, "User1 Product", "https://example.com/p1.jpg")
		product2ID := test.SeedProduct(t, pool, "User2 Product", "https://example.com/p2.jpg")

		test.AddProductToWatchlist(t, pool, user1ID, product1ID)
		test.AddProductToWatchlist(t, pool, user2ID, product2ID)

		source1D := test.SeedProductSource(t, pool, test.ProductSourceConfig{
				ProductID: 			product1ID,
				Platform: 			"amazon",
				PlatformProductID:  "p1",
				URL: 				"https://amazon.com/p1",
		})
		source2D := test.SeedProductSource(t, pool, test.ProductSourceConfig{
				ProductID: 			product2ID,
				Platform: 			"amazon",
				PlatformProductID:  "p2",
				URL: 				"https://amazon.com/p2",
		})

		test.SeedPriceSnapshot(t, pool, test.PriceSnapshotConfig{
				ProductSourceID: 	source1D,
				Price: 				100.00,
				InStock: 			true,	
		})
		test.SeedPriceSnapshot(t, pool, test.PriceSnapshotConfig{
				ProductSourceID: 	source2D,
				Price: 				200.00,
				InStock: 			true,	
		})

		products, err := db.FetchUserTrackedProducts(user1ID, pool)

		require.NoError(t, err)
		require.Len(t, products, 1, "User 1 only has one product")
		assert.Equal(t, "User1 Product", products[0].ProductName)
		assert.Equal(t, 100.00, products[0].LowestPrice)
	})

	t.Run("Handles product with multiple prices in different currencies", func(t *testing.T) {
		test.CleanupTables(t, pool)
		
		userID := test.SeedUser(t, pool, "currency@example.com")		
		productID := test.SeedProduct(t, pool, "Global Product", "https://example.com/global.jpg")
		test.AddProductToWatchlist(t, pool, userID, productID)

		usSourceID := test.SeedProductSource(t, pool, test.ProductSourceConfig{
				ProductID: 			productID,
				Platform: 			"amazon_us",
				PlatformProductID:  "US123",
				URL: 				"https://amazon.com/product",
		})
		ukSourceID := test.SeedProductSource(t, pool, test.ProductSourceConfig{
				ProductID: 			productID,
				Platform: 			"amazon_uk",
				PlatformProductID:  "UK123",
				URL: 				"https://amazon.co.uk/product",
		})

		test.SeedPriceSnapshot(t, pool, test.PriceSnapshotConfig{
				ProductSourceID: 	usSourceID,
				Price: 				99.99,
				Currency: 			"USD",
				InStock: 			true,	
		})
		test.SeedPriceSnapshot(t, pool, test.PriceSnapshotConfig{
				ProductSourceID: 	ukSourceID,
				Price: 				79.99,
				Currency:	 		"GBP",	
				InStock: 			true,	
		})

		products, err := db.FetchUserTrackedProducts(userID, pool)

		require.NoError(t, err)
		require.Len(t, products, 1, "User has one product with multiple currency sources")

		// Note: This test shows a limitation - comparing prices across currencies without conversion
		// The query will return the lowest numeric value, which may not be the actual cheapest
		assert.Equal(t, "Global Product", products[0].ProductName)
		assert.Equal(t, 79.99, products[0].LowestPrice, "GBP 79.99 is lower numerically than USD 99.99")
		assert.Equal(t, "amazon_uk", products[0].LowestSource)
	})

	t.Run("Handles multiple sources with same lowest price", func(t *testing.T) {
		test.CleanupTables(t, pool)

		userID := test.SeedUser(t, pool, "sameprice@example.com")
		productID := test.SeedProduct(t, pool, "Same Price Product", "https://example.com/same.jpg")
		test.AddProductToWatchlist(t, pool, userID, productID)

		platforms := []string{"amazon", "bestbuy", "newegg"}
		for _, platform := range platforms {
			sourceID := test.SeedProductSource(t, pool, test.ProductSourceConfig{
				ProductID:		 	productID,
				Platform: 			platform,
				PlatformProductID:  platform + "_id",
				URL: 				"https://" + platform + ".com/product",	
			})
			test.SeedPriceSnapshot(t, pool, test.PriceSnapshotConfig{
				ProductSourceID: sourceID,
				Price: 			 49.99,
				InStock: 		 true,
			})
		}

		products, err := db.FetchUserTrackedProducts(userID, pool)

		require.NoError(t, err)
		require.Len(t, products, 1)
		assert.Equal(t, 49.99, products[0].LowestPrice)
		assert.Contains(t, platforms, products[0].LowestSource, "Should pick one of the platforms with lowest price")
	})

	t.Run("Handles non-existant user gracefully", func(t *testing.T) {
		test.CleanupTables(t, pool)

		products, err := db.FetchUserTrackedProducts(99999, pool)

		require.NoError(t, err)
		assert.Empty(t, products)
	})
}

// Integration tests for InsertDeleteProductForUser SQL func
func TestDeleteProductForUser(t *testing.T) {
	pool := testDB.Pool
	t.Run("Returns nil when delete successful", func(t *testing.T) {
		test.CleanupTables(t, pool)

		userID := test.SeedUser(t, pool, "example@example.com")
		productID := test.SeedProduct(t, pool, "product", "https://imgur.com/123")
		test.AddProductToWatchlist(t, pool, userID, productID)

		err := db.DeleteProductForUser(userID, productID, pool)

		require.NoError(t, err)
	})

	t.Run("returns error when user doesnt exist to delete for", func(t *testing.T) {
		test.CleanupTables(t, pool)
		
		err := db.DeleteProductForUser(1, 2, pool)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found in user's watchlist")
	})

	t.Run("returns error when product doesnt exist in user watchlist", func(t *testing.T) {
		test.CleanupTables(t, pool)

		userID := test.SeedUser(t, pool, "example@example.com")

		err := db.DeleteProductForUser(userID, 2, pool)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found in user's watchlist")
	})
}