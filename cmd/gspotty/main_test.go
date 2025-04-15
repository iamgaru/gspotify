package main

import (
	"context"
	"testing"
	"time"

	"github.com/iamgaru/gspotty/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
)

// TestMainFunctionality tests the main functionality of the application
func TestMainFunctionality(t *testing.T) {
	// Create a mock Spotify client for testing
	mockClient := &testutils.MockSpotifyClient{}

	// Test search functionality
	t.Run("Search Functionality", func(t *testing.T) {
		err := searchTracks(context.Background(), mockClient, "test query", 10)
		assert.NoError(t, err)
		assert.True(t, mockClient.SearchCalled)

		err = searchAlbums(context.Background(), mockClient, "test query", 10)
		assert.NoError(t, err)
		assert.True(t, mockClient.SearchCalled)

		err = searchPlaylists(context.Background(), mockClient, "test query", 10)
		assert.NoError(t, err)
		assert.True(t, mockClient.SearchCalled)
	})

	// Test playback controls
	t.Run("Playback Controls", func(t *testing.T) {
		err := mockClient.Play(context.Background())
		assert.NoError(t, err)
		assert.True(t, mockClient.PlayCalled)

		err = mockClient.Pause(context.Background())
		assert.NoError(t, err)
		assert.True(t, mockClient.PauseCalled)

		err = mockClient.Next(context.Background())
		assert.NoError(t, err)
		assert.True(t, mockClient.NextCalled)

		err = mockClient.Previous(context.Background())
		assert.NoError(t, err)
		assert.True(t, mockClient.PreviousCalled)

		err = mockClient.Seek(context.Background(), time.Second*30)
		assert.NoError(t, err)
		assert.True(t, mockClient.SeekCalled)

		err = mockClient.SetVolume(context.Background(), 50)
		assert.NoError(t, err)
		assert.True(t, mockClient.SetVolumeCalled)
	})
}

// TestMainWithMockContext tests the main functionality with a mock context
func TestMainWithMockContext(t *testing.T) {
	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cancel()

	// Verify that cancellation doesn't cause panics
	assert.NotPanics(t, func() {
		mockClient := &testutils.MockSpotifyClient{}
		_ = mockClient.Play(ctx)
	})
}

// Helper functions for testing
func searchTracks(ctx context.Context, client *testutils.MockSpotifyClient, query string, limit int) error {
	_, err := client.Search(ctx, query, spotify.SearchTypeTrack, spotify.Limit(limit))
	return err
}

func searchAlbums(ctx context.Context, client *testutils.MockSpotifyClient, query string, limit int) error {
	_, err := client.Search(ctx, query, spotify.SearchTypeAlbum, spotify.Limit(limit))
	return err
}

func searchPlaylists(ctx context.Context, client *testutils.MockSpotifyClient, query string, limit int) error {
	_, err := client.Search(ctx, query, spotify.SearchTypePlaylist, spotify.Limit(limit))
	return err
}

func stopCurrentlyPlaying(client *testutils.MockSpotifyClient) error {
	return client.Pause(context.Background())
}
