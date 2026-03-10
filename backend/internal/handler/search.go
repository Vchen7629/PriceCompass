package handler

import (
	"backend/internal/service"
	"backend/internal/types"
	"encoding/json"
	"log"
	"net/http"
)

type SearchHandler struct {
	EbayClient *service.OauthClient
}

// handler for validating api input params
func NewSearchHandler(ebayClient *service.OauthClient) *SearchHandler {
	return &SearchHandler{
		EbayClient: ebayClient,
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

	if search.EbayClient == nil {
		http.Error(w, "Search service not available", http.StatusServiceUnavailable)
		return
	}

	//amazonCh := make(chan []types.SearchResult)
	ebayCh := make(chan []types.SearchResult)

	go func() {
		res, err := service.SearchProductEbay(search.EbayClient, productName)
		if err != nil {
			log.Printf("ebay search error: %v", err)
			res = []types.SearchResult{}
		}
		ebayCh <- res
	}()

	var allResults []types.SearchResult
	allResults = append(allResults, <-ebayCh...)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(allResults)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}