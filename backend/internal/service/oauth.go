package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TokenProvider interface for OAuth token management
// This allows for easy mocking in tests
type TokenProvider interface {
	GetAccessToken(ctx context.Context) (string, error)
}

type OauthClient struct {
	clientID 		string
	clientSecret 	string
	tokenURL 		string
	scope 			string
	// cached token data
	mu				sync.RWMutex
	accessToken		string
	expiresAt		time.Time
}

type TokenResp struct {
	AccessToken	string `json:"access_token"`
	ExpiresIn 	int    `json:"expires_in"` // seconds
	TokenType   string `json:"token_type"`
}

// returns a valid access token and fetches a new one if needed
func (c *OauthClient) GetAccessToken(ctx context.Context) (string, error) {
	// first check if we have a valid cached token
	c.mu.RLock()
	if c.accessToken != "" && time.Now().Before(c.expiresAt) {
		token := c.accessToken
		c.mu.RUnlock()
		return token, nil
	}
	c.mu.RUnlock()

	// fetch a new token if no cached
	return c.fetchNewToken(ctx)
}

// this is called if we dont have a cached token, requests a new access token from the api
func (c *OauthClient) fetchNewToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// first double-check if another goroutine didn't just fetch it
	if c.accessToken != "" && time.Now().Before(c.expiresAt) {
		return c.accessToken, nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", c.scope)

	req, err := http.NewRequestWithContext(ctx, "POST", c.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	// ebay requires base64 encoded
	credentials := base64.StdEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch token: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResp
	jsonErr := json.NewDecoder(resp.Body).Decode(&tokenResp)
	if jsonErr != nil {
		return "", fmt.Errorf("failed to decode token response: %w", jsonErr)
	}

	// cache the fetched token with a 5 min buffer before actual expiry
	c.accessToken = tokenResp.AccessToken
	c.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn - 300) * time.Second)

	return c.accessToken, nil
}
