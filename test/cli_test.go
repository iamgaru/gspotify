package test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnvironmentVariables tests the environment variable checking
func TestEnvironmentVariables(t *testing.T) {
	// Save original environment variables
	originalID := os.Getenv("SPOTIFY_ID")
	originalSecret := os.Getenv("SPOTIFY_SECRET")

	// Clean up after test
	defer func() {
		os.Setenv("SPOTIFY_ID", originalID)
		os.Setenv("SPOTIFY_SECRET", originalSecret)
	}()

	// Test case 1: Both variables set
	os.Setenv("SPOTIFY_ID", "test_id")
	os.Setenv("SPOTIFY_SECRET", "test_secret")
	client := getSpotifyClient(context.Background())
	assert.NotNil(t, client)

	// Test case 2: Missing variables
	os.Setenv("SPOTIFY_ID", "")
	os.Setenv("SPOTIFY_SECRET", "")
	assert.Panics(t, func() {
		checkEnvironmentVariables()
	})
}
