package handler

import (
	"backend/internal/service"
	"backend/internal/types"
	"encoding/json"
	"log"
	"net/http"
	"github.com/go-playground/validator/v10"
)

type SearchHandler struct {
	Validate 	*validator.Validate
	PlatformClients	*service.PlatformClients
}

// handler for validating api input params
func NewSearchHandler(clients *service.PlatformClients) *SearchHandler {
	return &SearchHandler{
		Validate: 			validator.New(),
		PlatformClients: 	clients,
	}
}

// GET route to search for a product using the name across all 4 platforms
// (amazon, ebay, newegg, bestbuy) and display it on the frontend
func (search *SearchHandler) SearchProductByName(w http.ResponseWriter, r *http.Request) {
	productName := r.PathValue("name")
	if productName == "" {
		http.Error(w, "No product name provided", http.StatusBadRequest)
		return
	}

	if search.PlatformClients == nil {
		http.Error(w, "Search service not available", http.StatusServiceUnavailable)
		return
	}

	// Goroutine to check all 4 platforms in parallel
	//amazonCh := make(chan []searchResult)
	ebayCh   := make(chan []types.SearchResult)
	bestbuyCh:= make(chan []types.SearchResult)
	neweggCh := make(chan []types.SearchResult)

	go func() {
		res, err := service.SearchProductEbay(search.PlatformClients.Ebay, productName)
		if err != nil {
			log.Printf("ebay search error: %v", err)
			res = []types.SearchResult{}
		}
		ebayCh <- res
	}()

	go func() { bestbuyCh <- []types.SearchResult{} }()
	go func() { neweggCh <- []types.SearchResult{} }()

	var allResults []types.SearchResult
	allResults = append(allResults, <-ebayCh...)
	allResults = append(allResults, <-bestbuyCh...)
	allResults = append(allResults, <-neweggCh...)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(allResults)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}