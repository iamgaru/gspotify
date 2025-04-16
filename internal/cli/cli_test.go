package cli

import (
	"context"
	"os"
	"testing"

	"github.com/iamgaru/gspotty/internal/testutils"
	"github.com/iamgaru/gspotty/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
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
	client := utils.GetSpotifyClient(context.Background())
	assert.NotNil(t, client)

	// Test case 2: Missing variables
	os.Setenv("SPOTIFY_ID", "")
	os.Setenv("SPOTIFY_SECRET", "")
	assert.Panics(t, func() {
		utils.CheckEnvironmentVariables()
	})
}

func TestCLI(t *testing.T) {
	ctx := context.Background()

	t.Run("Environment Variables", func(t *testing.T) {
		// Set required environment variables
		os.Setenv("SPOTIFY_CLIENT_ID", "test_client_id")
		os.Setenv("SPOTIFY_CLIENT_SECRET", "test_client_secret")
		os.Setenv("SPOTIFY_REDIRECT_URI", "http://localhost:8080/callback")

		assert.NotPanics(t, func() {
			utils.CheckEnvironmentVariables()
		})

		client := utils.GetSpotifyClient(ctx)
		assert.NotNil(t, client)
	})

	t.Run("Client Creation", func(t *testing.T) {
		// Set required environment variables
		os.Setenv("SPOTIFY_CLIENT_ID", "test_client_id")
		os.Setenv("SPOTIFY_CLIENT_SECRET", "test_client_secret")
		os.Setenv("SPOTIFY_REDIRECT_URI", "http://localhost:8080/callback")

		client := utils.GetSpotifyClient(ctx)
		assert.NotNil(t, client)
	})

	t.Run("Search Functionality", func(t *testing.T) {
		mockClient := &testutils.MockSpotifyClient{}

		// Test track search
		t.Run("Track Search", func(t *testing.T) {
			query := "test track"
			_, err := mockClient.Search(context.Background(), query, spotify.SearchTypeTrack)
			assert.NoError(t, err)
			assert.True(t, mockClient.SearchCalled)
		})

		// Test album search
		t.Run("Album Search", func(t *testing.T) {
			query := "test album"
			_, err := mockClient.Search(context.Background(), query, spotify.SearchTypeAlbum)
			assert.NoError(t, err)
			assert.True(t, mockClient.SearchCalled)
		})

		// Test playlist search
		t.Run("Playlist Search", func(t *testing.T) {
			query := "test playlist"
			_, err := mockClient.Search(context.Background(), query, spotify.SearchTypePlaylist)
			assert.NoError(t, err)
			assert.True(t, mockClient.SearchCalled)
		})
	})
}

func TestCLIWithMockContext(t *testing.T) {
	// Set required environment variables
	os.Setenv("SPOTIFY_CLIENT_ID", "test_client_id")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "test_client_secret")
	os.Setenv("SPOTIFY_REDIRECT_URI", "http://localhost:8080/callback")

	ctx := context.Background()
	assert.NotPanics(t, func() {
		utils.CheckEnvironmentVariables()
	})

	client := utils.GetSpotifyClient(ctx)
	assert.NotNil(t, client)
	assert.NotNil(t, ctx)
}
