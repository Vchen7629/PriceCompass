//go:build unit

package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"backend/internal/handler"
	"backend/pkg/test"
)

// Unit tests for the AddProductName Handler function
func TestAddProductNameHandler(t *testing.T) {
	t.Run("Empty request body", func(t *testing.T) {
		mockPool := test.NewMockPool()
		mockHandler := handler.NewAPI(mockPool)

		// Create empty request body
		req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockHandler.AddProductName(w, req)

		// Verify we get a bad request status
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Invalid json payload handling", func(t *testing.T) {
		mockPool := test.NewMockPool()
		mockHandler := handler.NewAPI(mockPool)
		
		// Create payload with invalid data
		invalidPayload := map[string]interface{any}{
			"user_id":      -1,
			"product_name": "a",
		}
		body, _ := json.Marshal(invalidPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockHandler.AddProductName(w, req)

		// Verify we get a bad request status
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// Unit tests for the GetUserTrackedProducts Handler function
func TestGetUserTrackedProductsHandler(t *testing.T) {
	t.Run("Empty User ID", func(t *testing.T) {
		mockPool := test.NewMockPool()
		mockHandler := handler.NewAPI(mockPool)

		req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockHandler.GetUserTrackedProducts(w, req)

		// Verify we get a bad request status
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
	t.Run("Non integer UserID", func(t *testing.T) {
		mockPool := test.NewMockPool()
		mockHandler := handler.NewAPI(mockPool)

		req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
		req.SetPathValue("id", "invalid")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockHandler.GetUserTrackedProducts(w, req)

		// Verify we get a bad request status
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// Unit tests for the DeleteProductHandler function
func TestDeleteProductHandler(t *testing.T) {
	t.Run("Empty request body", func(t *testing.T) {
		mockPool := test.NewMockPool()
		mockHandler := handler.NewAPI(mockPool)

		// Create empty request body
		req := httptest.NewRequest(http.MethodDelete, "/api/products", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockHandler.DeleteProduct(w, req)

		// Verify we get a bad request status
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
	t.Run("Invalid json payload handling", func(t *testing.T) {
		mockPool := test.NewMockPool()
		mockHandler := handler.NewAPI(mockPool)
		
		// Create payload with invalid data
		invalidPayload := map[string]interface{any}{
			"user_id":      -1,
			"product_id": 	-1,
		}
		body, _ := json.Marshal(invalidPayload)

		req := httptest.NewRequest(http.MethodDelete, "/api/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockHandler.DeleteProduct(w, req)

		// Verify we get a bad request status
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

