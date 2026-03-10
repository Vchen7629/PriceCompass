package service

import (
	"backend/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type SearchRes struct {
	ItemSummaries 	[]Item `json:"itemSummaries"`
	Total			int `json:"total"`
}

type Item struct {
	Name 		string 	`json:"title"`
	ItemID 		string 	`json:"itemId"`
	Price		Price	`json:"price"`	
	Image	 	Image  	`json:"image"`
	ItemURL  	string 	`json:"itemWebUrl"`
	Condition 	string 	`json:"condition"`
}

type Image struct {
	ImageURL string `json:"imageUrl"`
}

type Price struct {
	Value 		string `json:"value"`
	Currency 	string `json:"currency"`
}

// SearchProductEbay finds all products on eBay matching the query
func SearchProductEbay(client *OauthClient, query string) ([]types.SearchResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	additionalHeaders := map[string]string{"X-EBAY-C-MARKETPLACE-ID": "EBAY_US"}
	searchURL := fmt.Sprintf("https://api.ebay.com/buy/browse/v1/item_summary/search?q=%s&limit=50", url.QueryEscape(query))

	body, err := MakeAuthenticatedRequest(ctx, client, searchURL, additionalHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to make authenticated request: %w", err)
	}

	var searchResp SearchRes
	jsonErr := json.Unmarshal(body, &searchResp)

	if jsonErr != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	if len(searchResp.ItemSummaries) == 0 {
		return nil, fmt.Errorf("No results found for query: %s", query)
	}

	results := make([]types.SearchResult, 0, len(searchResp.ItemSummaries))
	for _, item := range searchResp.ItemSummaries {
		results = append(results, types.SearchResult{
			ProductID:		item.ItemID,
			ProductName: 	item.Name,
			Price: 			item.Price.Value,
			Url:			item.ItemURL,
			ImageUrl: 		item.Image.ImageURL,
			Source: 		"ebay",
		})
	}

	return results, nil
}