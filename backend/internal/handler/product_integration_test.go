//go:build integration

package handler_test

import (
	"backend/internal/handler"
	"backend/internal/test"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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

// Integration tests for InsertProductForUser route handler
func TestAddProductNameHandlerIntegration(t *testing.T) {
	pool := testDB.Pool
	t.Run("successful product encoded properly", func(t *testing.T) {
		test.CleanupTables(t, pool)

		userID := test.SeedUser(t, pool, "example@example.com")

		payload := map[string]interface{any}{
			"user_id":      userID,
			"product_name": "gpu",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		h := handler.NewAPI(pool, nil)
		h.AddProductName(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Success should return 200")

		var response map[string]interface{any}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err, "Response should be valid JSON")

		// Assert response contains expected fields
		assert.Equal(t, "gpu", response["product_name"], "Product name should match")
		assert.NotNil(t, response["product_id"], "Response should contain product ID")
	})
}

// Integration tests for GetUserTrackedProducts route handler
func TestGetUserTrackedProductsIntegration(t *testing.T) {
	pool := testDB.Pool
	t.Run("successful get products response encoded properly", func(t *testing.T) {
		test.CleanupTables(t, pool)

		userID := test.SeedUser(t, pool, "example@example.com")

		products := []string{"product_1", "product_2", "product_3"}

		for _, product := range products {
			productID := test.SeedProduct(t, pool, product, "http://imgur.com/idk")

			test.AddProductToWatchlist(t, pool, userID, productID)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
		req.SetPathValue("id", strconv.Itoa(userID))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		h := handler.NewAPI(pool, nil)
		h.GetUserTrackedProducts(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Success should return 200")

		// Decode response as array of products
		var response []map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err, "Response should be valid JSON")

		assert.Equal(t, 3, len(response), "Should return 3 products")

		// Check each product has required fields
		expectedNames := map[string]bool{"product_1": true, "product_2": true, "product_3": true}
		for _, product := range response {
			assert.NotNil(t, product["product_id"], "Product should have ID")
			assert.NotNil(t, product["product_name"], "Product should have name")
			assert.NotNil(t, product["image_url"], "Product should have image URL")

			productName := product["product_name"].(string)
			assert.True(t, expectedNames[productName], "Product name should be one of the seeded products")
		}
	})
}

// Integration tests for DeleteProduct route handler
func TestDeleteProductrIntegration(t *testing.T) {
	pool := testDB.Pool
	t.Run("successful delete encoded properly", func(t *testing.T) {
		test.CleanupTables(t, pool)

		userID := test.SeedUser(t, pool, "example@example.com")
		productID := test.SeedProduct(t, pool, "product", "image.jpg")
		test.AddProductToWatchlist(t, pool, userID, productID)

		payload := map[string]interface{any}{
			"user_id":      userID,
			"product_id": 	productID,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodDelete, "/api/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		h := handler.NewAPI(pool, nil)
		h.DeleteProduct(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Success should return 200")

		// Decode response as a string (not an object)
		var response string
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err, "Response should be valid JSON")

		// Assert response is the success message
		assert.Equal(t, "Successfully deleted product", response, "Should return success message")
	})
}