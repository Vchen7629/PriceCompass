//go:build unit

package handler_test

import (
	"backend/internal/handler"
	"backend/pkg/test"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Unit tests for the UserSignUp Handler function
func TestUserSignUpHandler(t *testing.T) {
	t.Run("Empty request body", func(t *testing.T) {
		mockPool := test.NewMockPool()
		mockHandler := handler.NewHandler(mockPool)

		// Create empty request body
		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockHandler.UserSignUp(w, req)

		// Verify we get a bad request status
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Invalid json payload handling", func(t *testing.T) {
		mockPool := test.NewMockPool()
		mockHandler := handler.NewHandler(mockPool)
		
		// Create payload with invalid data
		invalidPayload := map[string]interface{any}{
			"username":      "a",
			"password": "a",
		}
		body, _ := json.Marshal(invalidPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockHandler.UserSignUp(w, req)

		// Verify we get a bad request status
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// Unit tests for the UserLogin Handler function
func TestUserLoginHandler(t *testing.T) {
	t.Run("Empty request body", func(t *testing.T) {
		mockPool := test.NewMockPool()
		mockHandler := handler.NewHandler(mockPool)

		// Create empty request body
		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockHandler.UserLogin(w, req)

		// Verify we get a bad request status
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Invalid json payload handling", func(t *testing.T) {
		mockPool := test.NewMockPool()
		mockHandler := handler.NewHandler(mockPool)
		
		// Create payload with invalid data
		invalidPayload := map[string]interface{any}{
			"username":      "a",
			"password": "a",
		}
		body, _ := json.Marshal(invalidPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockHandler.UserLogin(w, req)

		// Verify we get a bad request status
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
