package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"backend/internal/store"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	users		store.UserStore
	validate 	*validator.Validate
}

func NewUserHandler(users store.UserStore) *UserHandler {
	return &UserHandler{
		users: 		users,
		validate: 	validator.New(),
	}
}

// Create a new user account
func (h *UserHandler) UserSignUp(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username" validate:"required,min=2"`
		Email string `json:"email" validate:"required,min=4"`
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

	insertErr := h.users.InsertNewUser(payload.Username, payload.Email, payload.Password)
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
func (h *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
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

	sessionToken, err := h.users.LoginUser(payload.Username, payload.Password)
	if err != nil {
		if err.Error() == "passwords don't match, can't login" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
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
	
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
