package service

import (
	"fmt"
	"os"
)

type PlatformClients struct {
	Ebay		*OauthClient
	BestBuy 	*OauthClient
	Newegg 		*OauthClient
}

// Initialize all platform clients for search service
func NewPlatformClients() (*PlatformClients, error) {
	ebayClientID := os.Getenv("EBAY_CLIENT_ID")
	ebayClientSecret := os.Getenv("EBAY_CLIENT_SECRET")
	if ebayClientID == "" || ebayClientSecret == "" {
		return nil, fmt.Errorf("missing required env vars: EBAY_CLIENT_ID and/or EBAY_CLIENT_SECRET")
	}

	bestbuyAPIKEY := os.Getenv("BESTBUY_API_KEY")
	if bestbuyAPIKEY == "" {
		return nil, fmt.Errorf("missing required env variable: BESTBUY_API_KEY")
	}

	return &PlatformClients{
		Ebay: &OauthClient{
			clientID:     ebayClientID,
			clientSecret: ebayClientSecret,
			tokenURL:     "https://api.ebay.com/identity/v1/oauth2/token",
			scope:        "https://api.ebay.com/oauth/api_scope",
		},
		BestBuy: &OauthClient{
			clientID:     bestbuyAPIKEY,
			clientSecret: "",
			tokenURL:     "",
			scope:        "",
		},
		Newegg: nil,
	}, nil
}