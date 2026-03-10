package handler

import (
	"backend/internal/store"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	products	store.ProductStore
	validate 	*validator.Validate
}

func NewProductHandler(products store.ProductStore) *ProductHandler {
	return &ProductHandler{
		products: 	products,
		validate: 	validator.New(),
	}
}

// POST route to add a new product for a user
func (h *ProductHandler) AddProductName(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserId		int `json:"user_id" validate:"required,gte=0"`
		ProductName string `json:"product_name" validate:"required,min=2"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = h.Validate.Struct(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	product, dbErr := h.products.InsertProductForUser(payload.UserId, payload.ProductName)
	if dbErr != nil {
		log.Println(dbErr)
		if db.HandleDatabaseErrors(w, dbErr) {
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode(product)
	if encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GET route to fetch a list of the user's tracked products with product metadata like
// name, lowest price, lowest source, available from the database
func (h *ProductHandler) GetUserTrackedProducts(w http.ResponseWriter, r *http.Request) {
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

	productList, dbErr := h.products.FetchUserTrackedProducts(userID)
	if dbErr != nil {
		log.Println(dbErr)
		if db.HandleDatabaseErrors(w, dbErr) {
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode(productList)
	if encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DELETE route to delete a product to be tracked using
// the product id in the database
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserId		int `json:"user_id" validate:"required,gte=0"`
		ProductID 	int `json:"product_id" validate:"required,gte=0"`
	}


	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = api.Validate.Struct(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}


	dbErr := h.products.DeleteProductForUser(payload.UserId, payload.ProductID)
	if dbErr != nil {
		if dbErr.Error() == "product not found in user's watchlist" {
			http.Error(w, dbErr.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode("Successfully deleted product")
	if encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}