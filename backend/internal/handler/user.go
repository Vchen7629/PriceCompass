package handler

import (
	"backend/internal/db"
	"encoding/json"
	"fmt"
	"net/http"
)

// Create a new user account
func (h *Handler) UserSignUp(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username" validate:"required,min=2"`
		Password string `json:"password" validate:"required,min=2"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid json payload", http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insertErr := db.InsertNewUser(payload.Username, payload.Password, h.pool)
	if insertErr != nil {
		http.Error(w, insertErr.Error(), http.StatusInternalServerError)
		return
	}

	res := fmt.Sprintf("Created new user: %s", payload.Username)
	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode(res)
	if encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// login to existing user account
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username" validate:"required,min=2"`
		Password string `json:"password" validate:"required,min=2"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid json payload", http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionToken, err := db.LoginUser(payload.Username, payload.Password, h.pool)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   86400, // 24 hours in seconds
		HttpOnly: true,  
		Secure:   true,  
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})
}
