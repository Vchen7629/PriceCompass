package service

import (
	"fmt"
	"os"
)

// Ebay client constructor, creates the client with api key, permissions and scope
func NewEbayClient() (*OauthClient, error) {
	ebayClientID := os.Getenv("EBAY_CLIENT_ID")
	ebayClientSecret := os.Getenv("EBAY_CLIENT_SECRET")
	if ebayClientID == "" || ebayClientSecret == "" {
		return nil, fmt.Errorf("missing required env vars: EBAY_CLIENT_ID and/or EBAY_CLIENT_SECRET")
	}

	return &OauthClient{
		clientID:     ebayClientID,
		clientSecret: ebayClientSecret,
		tokenURL:     "https://api.ebay.com/identity/v1/oauth2/token",
		scope:        "https://api.ebay.com/oauth/api_scope",
	}, nil
}