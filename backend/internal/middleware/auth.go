package middleware

import (
	"backend/internal/store"
	"context"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "user"

type UserContext struct {
	UserId 		int
	Username 	string
}

type Handler struct {
	handler		store.ValidationStore
}

func NewMiddlewareHandler(store store.ValidationStore) *Handler {
	return &Handler{handler: store}
}

// validates session by checking the session token in the cookie
// for every request to fetch protected user data
func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized: no session token", http.StatusUnauthorized)
			return
		}

		userId, username, err := h.handler.ValidateSession(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)
			return
		}

		userContext := UserContext{
			UserId: 	userId,
			Username: 	username,
		}

		ctx := context.WithValue(r.Context(), userContextKey, userContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}