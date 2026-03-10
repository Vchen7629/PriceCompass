//go:build unit

package service_test

import (
	"backend/internal/service"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Unit tests for NewPlatformClients function
func TestNewPlatformClients(t *testing.T) {

	t.Run("throws error on missing ebay env variables", func(t *testing.T) {
		os.Setenv("EBAY_CLIENT_ID", "")
		os.Setenv("EBAY_CLIENT_SECRET", "")

		_, err := service.NewPlatformClients()

		assert.Equal(t, "missing required env vars: EBAY_CLIENT_ID and/or EBAY_CLIENT_SECRET", err.Error())
	})

	t.Run("throws error on missing best buy env variables", func(t *testing.T) {
		os.Setenv("EBAY_CLIENT_ID", "hi")
		os.Setenv("EBAY_CLIENT_SECRET", "hi")
		os.Setenv("BESTBUY_API_KEY", "")

		_, err := service.NewPlatformClients()

		assert.Equal(t, "missing required env variable: BESTBUY_API_KEY", err.Error())
	})

	t.Run("throws error on missing all env variables", func(t *testing.T) {
		os.Setenv("EBAY_CLIENT_ID", "")
		os.Setenv("EBAY_CLIENT_SECRET", "")
		os.Setenv("BESTBUY_API_KEY", "")

		_, err := service.NewPlatformClients()

		assert.Equal(t, "missing required env vars: EBAY_CLIENT_ID and/or EBAY_CLIENT_SECRET", err.Error())
	})
}