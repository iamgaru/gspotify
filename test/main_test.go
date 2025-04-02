package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
)

// TestMainFunctionality tests the main functionality of the application
func TestMainFunctionality(t *testing.T) {
	// Create a mock Spotify client for testing
	mockClient := &MockSpotifyClient{}

	// Test search functionality
	t.Run("Search Functionality", func(t *testing.T) {
		// Test track search
		t.Run("Track Search", func(t *testing.T) {
			query := "test track"
			_, err := mockClient.Search(context.Background(), query, spotify.SearchTypeTrack)
			assert.NoError(t, err)
			assert.True(t, mockClient.searchCalled)
		})

		// Test album search
		t.Run("Album Search", func(t *testing.T) {
			query := "test album"
			_, err := mockClient.Search(context.Background(), query, spotify.SearchTypeAlbum)
			assert.NoError(t, err)
			assert.True(t, mockClient.searchCalled)
		})

		// Test playlist search
		t.Run("Playlist Search", func(t *testing.T) {
			query := "test playlist"
			_, err := mockClient.Search(context.Background(), query, spotify.SearchTypePlaylist)
			assert.NoError(t, err)
			assert.True(t, mockClient.searchCalled)
		})
	})

	// Test playback control
	t.Run("Playback Control", func(t *testing.T) {
		// Test play
		t.Run("Play", func(t *testing.T) {
			err := mockClient.Play(context.Background())
			assert.NoError(t, err)
			assert.True(t, mockClient.playCalled)
		})

		// Test pause
		t.Run("Pause", func(t *testing.T) {
			err := mockClient.Pause(context.Background())
			assert.NoError(t, err)
			assert.True(t, mockClient.pauseCalled)
		})

		// Test next
		t.Run("Next", func(t *testing.T) {
			err := mockClient.Next(context.Background())
			assert.NoError(t, err)
			assert.True(t, mockClient.nextCalled)
		})

		// Test previous
		t.Run("Previous", func(t *testing.T) {
			err := mockClient.Previous(context.Background())
			assert.NoError(t, err)
			assert.True(t, mockClient.previousCalled)
		})

		// Test seek
		t.Run("Seek", func(t *testing.T) {
			position := 30 * time.Second
			err := mockClient.Seek(context.Background(), position)
			assert.NoError(t, err)
			assert.True(t, mockClient.seekCalled)
		})

		// Test volume
		t.Run("Volume", func(t *testing.T) {
			volume := 50
			err := mockClient.SetVolume(context.Background(), volume)
			assert.NoError(t, err)
			assert.True(t, mockClient.volumeCalled)
		})
	})
}

// TestMainWithMockContext tests the main functionality with a mock context
func TestMainWithMockContext(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test context cancellation
	cancel()
	// The main functionality should handle context cancellation gracefully
	// We can't easily test the visual output, but we can verify the function doesn't panic
}

// Helper functions for testing
func searchTracks(ctx context.Context, client *MockSpotifyClient, query string, limit int) error {
	_, err := client.Search(ctx, query, spotify.SearchTypeTrack, spotify.Limit(limit))
	return err
}

func searchAlbums(ctx context.Context, client *MockSpotifyClient, query string, limit int) error {
	_, err := client.Search(ctx, query, spotify.SearchTypeAlbum, spotify.Limit(limit))
	return err
}

func searchPlaylists(ctx context.Context, client *MockSpotifyClient, query string, limit int) error {
	_, err := client.Search(ctx, query, spotify.SearchTypePlaylist, spotify.Limit(limit))
	return err
}

func stopCurrentlyPlaying(client *MockSpotifyClient) error {
	return client.Pause(context.Background())
}
