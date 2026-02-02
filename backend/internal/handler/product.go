package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	"backend/internal/db"
	pkgdb "backend/pkg/db"
)

type Handler struct {
	pool 		*pgxpool.Pool
	validate 	*validator.Validate
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: 		pool,
		validate: 	validator.New(),
	}
}

// POST route to add a new product for a user
func (h *Handler) AddProductName(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserId		int `json:"user_id" validate:"required,gte=0"`
		ProductName string `json:"product_name" validate:"required,min=2"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	product, dbErr := db.InsertProductForUser(payload.UserId, payload.ProductName, h.pool)
	if dbErr != nil {
		log.Println(dbErr)
		if pkgdb.HandleDatabaseErrors(w, dbErr) {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode(product)
	if encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GET route to fetch a list of the user's tracked products with product metadata like
// name, lowest price, lowest source, available from the database
func (h *Handler) GetUserTrackedProducts(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("id")
	if user_id == "" {
		http.Error(w, "Missing required user_id param", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(user_id)
	if err != nil {
		http.Error(w, "Invalid user_id: must be a string", http.StatusBadRequest)
		return
	}

	productList, dbErr := db.FetchUserTrackedProducts(userID, h.pool)
	if dbErr != nil {
		log.Println(dbErr)
		if pkgdb.HandleDatabaseErrors(w, dbErr) {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode(productList)
	if encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DELETE route to delete a product to be tracked using
// the product id in the database
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserId		int `json:"user_id" validate:"required,gte=0"`
		ProductID 	int `json:"product_id" validate:"required,gte=0"`
	}


	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}


	dbErr := db.DeleteProductForUser(payload.UserId, payload.ProductID, h.pool)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode("Successfully deleted product")
	if encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}