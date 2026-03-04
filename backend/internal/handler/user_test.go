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

// Unit tests for the UserSignUp Handler function
func TestUserSignUpHandler(t *testing.T) {
	t.Run("Empty request body", func(t *testing.T) {
		mockHandler := handler.NewUserHandler(&handler.MockUserStore{})

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
		mockHandler := handler.NewUserHandler(&handler.MockUserStore{})
		
		// Create payload with invalid data
		invalidPayload := map[string]interface{}{
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

	t.Run("db error returns 500", func(t *testing.T) {
		mock := &handler.MockUserStore{InsertUserErr: errors.New("db error")}
		mockHandler := handler.NewUserHandler(mock)

		payload := map[string]interface{}{
			"username": "idk", 
			"email": "idk@gmail.com", 
			"password": "pass123",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockHandler.UserSignUp(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "It should return http status 500")
	})
}

// Unit tests for the UserLogin Handler function
func TestUserLoginHandler(t *testing.T) {
	t.Run("Empty request body", func(t *testing.T) {
		mockHandler := handler.NewUserHandler(&handler.MockUserStore{})

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
		mockHandler := handler.NewUserHandler(&handler.MockUserStore{})
		
		// Create payload with invalid data
		invalidPayload := map[string]interface{}{
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

	t.Run("invalid password returns 401", func(t *testing.T) {
		mock := &handler.MockUserStore{LoginUserErr: errors.New("passwords don't match, can't login")}
		mockHandler := handler.NewUserHandler(mock)

		payload := map[string]interface{}{"username": "idk", "password": "pass123"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockHandler.UserLogin(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "It should return http status 401")
	})

	t.Run("generic db error returns 500", func(t *testing.T) {
		mock := &handler.MockUserStore{LoginUserErr: errors.New("db error")}
		mockHandler := handler.NewUserHandler(mock)

		payload := map[string]interface{}{"username": "idk", "password": "pass123"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockHandler.UserLogin(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "It should return http status 500")
	})
}
