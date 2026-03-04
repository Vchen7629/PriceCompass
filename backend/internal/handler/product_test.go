//go:build unit

package handler_test

import (
	"backend/internal/handler"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Unit tests for the AddProductName Handler function
func TestAddProductNameHandler(t *testing.T) {
	t.Run("Empty request body", func(t *testing.T) {
		mockHandler := handler.NewProductHandler(&handler.MockProductStore{})

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
		mockHandler := handler.NewProductHandler(&handler.MockProductStore{})
		
		// Create payload with invalid data
		invalidPayload := map[string]interface{}{
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

	t.Run("db error returns 500", func(t *testing.T) {
		mock := &handler.MockProductStore{InsertProductErr: errors.New("db error")}
		mockHandler := handler.NewProductHandler(mock)

		payload := map[string]interface{}{"user_id": 1, "product_name": "gpu"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockHandler.AddProductName(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "It should return http status 500")
	})
}

// Unit tests for the GetUserTrackedProducts Handler function
func TestGetUserTrackedProductsHandler(t *testing.T) {
	t.Run("Empty User ID", func(t *testing.T) {
		mockHandler := handler.NewProductHandler(&handler.MockProductStore{})

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
		mockHandler := handler.NewProductHandler(&handler.MockProductStore{})

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

	t.Run("db error returns 500", func(t *testing.T) {
		mock := &handler.MockProductStore{FetchProductsErr: errors.New("db error")}
		mockHandler := handler.NewProductHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/get/2", nil)
		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/products/get/{id...}", mockHandler.GetUserTrackedProducts)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "It should return http status 500")
	})
}

// Unit tests for the DeleteProductHandler function
func TestDeleteProductHandler(t *testing.T) {
	t.Run("Empty request body", func(t *testing.T) {
		mockHandler := handler.NewProductHandler(&handler.MockProductStore{})

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
		mockHandler := handler.NewProductHandler(&handler.MockProductStore{})

		// Create payload with invalid data
		invalidPayload := map[string]interface{}{
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

	t.Run("item to be deleted isnt found returns 404", func(t *testing.T) {
		mock := &handler.MockProductStore{DeleteProductErr: errors.New("product not found in user's watchlist")}
		mockHandler := handler.NewProductHandler(mock)

		payload := map[string]interface{}{"user_id": 1, "product_id": 2}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products/delete", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockHandler.DeleteProduct(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "It should return http status 404")
	})

	t.Run("generic db error returns 500", func(t *testing.T) {
		mock := &handler.MockProductStore{DeleteProductErr: errors.New("db error")}
		mockHandler := handler.NewProductHandler(mock)

		payload := map[string]interface{}{"user_id": 1, "product_id": 2}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products/delete", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockHandler.DeleteProduct(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "It should return http status 500")
	})
}

