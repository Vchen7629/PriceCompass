package store

import (
	"backend/internal/types"
	"net/http"
)

// Repository interface for storing the interfaces
// used by user handlers
type UserStore interface {
	InsertNewUser(username, email, password string) error
	LoginUser(username, password string) (string, error)
}

type ProductStore interface {
	InsertProductForUser(userID int, productName string) (types.Product, error)
	FetchUserTrackedProducts(userID int) ([]types.UserProduct, error)
	DeleteProductForUser(userID, productID int) error
}

type MiddlewareStore interface {
	Logging(next http.Handler) http.Handler
	AuthMiddleware(next http.HandlerFunc) http.HandlerFunc
}

type ValidationStore interface {
	ValidateSession(sessionToken string) (int, string, error)
}
