package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Handles OAuth token fetching and HTTP req execution
func MakeAuthenticatedRequest(
	ctx context.Context, 
	authClient *OauthClient,
	url string,
	headers map[string]string,
) ([]byte, error) {
	token, err := authClient.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	// set additional headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("req failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}