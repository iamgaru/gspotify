package profile

import (
	"context"
	"os"
	"testing"

	"github.com/iamgaru/gspotty/internal/testutils"
	"github.com/stretchr/testify/assert"
)

// MockGetProfile is a test helper that uses a mock client
func MockGetProfile(t *testing.T) {
	// Save original environment variables
	originalID := os.Getenv("SPOTIFY_ID")
	originalSecret := os.Getenv("SPOTIFY_SECRET")

	// Clean up after test
	defer func() {
		os.Setenv("SPOTIFY_ID", originalID)
		os.Setenv("SPOTIFY_SECRET", originalSecret)
		*userID = "" // Reset the flag
	}()

	t.Run("Missing User ID", func(t *testing.T) {
		*userID = ""
		assert.NotPanics(t, func() {
			GetProfile()
		})
	})

	t.Run("With Valid Configuration", func(t *testing.T) {
		// Set required environment variables
		os.Setenv("SPOTIFY_ID", "test_id")
		os.Setenv("SPOTIFY_SECRET", "test_secret")
		*userID = "test_user"

		assert.NotPanics(t, func() {
			GetProfile()
		})
	})
}

func TestProfile(t *testing.T) {
	ctx := context.Background()
	mockClient := &testutils.MockSpotifyClient{}

	// Save the original function and restore it after the test
	originalGetClient := getSpotifyClientFunc
	defer func() {
		getSpotifyClientFunc = originalGetClient
	}()

	// Replace with mock function
	getSpotifyClientFunc = func(ctx context.Context) SpotifyClient {
		return mockClient
	}

	t.Run("Display Profile", func(t *testing.T) {
		displayProfile(ctx, mockClient, "test_user")
		// Test passes if no panic occurs
	})

	t.Run("GetProfile Function", func(t *testing.T) {
		// Save original environment variables
		originalID := os.Getenv("SPOTIFY_ID")
		originalSecret := os.Getenv("SPOTIFY_SECRET")

		// Clean up after test
		defer func() {
			os.Setenv("SPOTIFY_ID", originalID)
			os.Setenv("SPOTIFY_SECRET", originalSecret)
			*userID = "" // Reset the flag
		}()

		t.Run("Missing User ID", func(t *testing.T) {
			*userID = ""
			assert.NotPanics(t, func() {
				GetProfile()
			})
		})

		t.Run("With Valid Configuration", func(t *testing.T) {
			*userID = "test_user"
			assert.NotPanics(t, func() {
				GetProfile()
			})
		})
	})
}
