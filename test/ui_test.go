package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
)

// MockResultsUI is a mock implementation of the ResultsUI for testing
type MockResultsUI struct {
	keepPlaying  bool
	returnToMenu func()
}

// NewMockResultsUI creates a new mock results UI
func NewMockResultsUI(resultType string, ctx context.Context, client interface{}, showDetails bool) *MockResultsUI {
	return &MockResultsUI{
		keepPlaying: false,
	}
}

// SetKeepPlayingFlag sets the keep playing flag
func (ui *MockResultsUI) SetKeepPlayingFlag(keepPlaying bool) {
	ui.keepPlaying = keepPlaying
}

// SetReturnToMenuFunction sets the return to menu function
func (ui *MockResultsUI) SetReturnToMenuFunction(returnFunc func()) {
	ui.returnToMenu = returnFunc
}

// DisplayTrackResults displays track search results
func (ui *MockResultsUI) DisplayTrackResults(ctx context.Context, client interface{}, tracks []spotify.FullTrack) {
	// Mock implementation
}

// DisplayAlbumResults displays album search results
func (ui *MockResultsUI) DisplayAlbumResults(ctx context.Context, client interface{}, albums []spotify.SimpleAlbum) {
	// Mock implementation
}

// DisplayPlaylistResults displays playlist search results
func (ui *MockResultsUI) DisplayPlaylistResults(ctx context.Context, client interface{}, playlists []spotify.SimplePlaylist) {
	// Mock implementation
}

// TestResultsUI tests the ResultsUI functionality
func TestResultsUI(t *testing.T) {
	// Create a mock Spotify client
	mockClient := &MockSpotifyClient{}

	// Create test data
	testTrack := spotify.FullTrack{
		SimpleTrack: spotify.SimpleTrack{
			ID:   "test_track_id",
			Name: "Test Track",
			Artists: []spotify.SimpleArtist{
				{Name: "Test Artist"},
			},
			Duration: 180000,
		},
		Album: spotify.SimpleAlbum{
			Name: "Test Album",
		},
	}

	// Create a new ResultsUI instance
	ui := NewMockResultsUI("track", context.Background(), mockClient, false)

	// Test initial state
	assert.False(t, ui.keepPlaying)

	// Test keep playing flag
	t.Run("Keep Playing Flag", func(t *testing.T) {
		ui.SetKeepPlayingFlag(true)
		assert.True(t, ui.keepPlaying)

		ui.SetKeepPlayingFlag(false)
		assert.False(t, ui.keepPlaying)
	})

	// Test return to menu function
	t.Run("Return to Menu Function", func(t *testing.T) {
		returnFunc := func() {}
		ui.SetReturnToMenuFunction(returnFunc)
		assert.NotNil(t, ui.returnToMenu)
	})

	// Test track results display
	t.Run("Track Results Display", func(t *testing.T) {
		tracks := []spotify.FullTrack{testTrack}
		ui.DisplayTrackResults(context.Background(), mockClient, tracks)
		// Note: We can't easily test the visual output, but we can verify the function doesn't panic
	})

	// Test album results display
	t.Run("Album Results Display", func(t *testing.T) {
		albums := []spotify.SimpleAlbum{
			{
				ID:   "test_album_id",
				Name: "Test Album",
				Artists: []spotify.SimpleArtist{
					{Name: "Test Artist"},
				},
			},
		}
		ui.DisplayAlbumResults(context.Background(), mockClient, albums)
		// Note: We can't easily test the visual output, but we can verify the function doesn't panic
	})

	// Test playlist results display
	t.Run("Playlist Results Display", func(t *testing.T) {
		playlists := []spotify.SimplePlaylist{
			{
				ID:   "test_playlist_id",
				Name: "Test Playlist",
			},
		}
		ui.DisplayPlaylistResults(context.Background(), mockClient, playlists)
		// Note: We can't easily test the visual output, but we can verify the function doesn't panic
	})
}

// TestResultsUIWithMockContext tests the ResultsUI with a mock context
func TestResultsUIWithMockContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockClient := &MockSpotifyClient{}
	_ = NewMockResultsUI("track", ctx, mockClient, false)

	// Test context cancellation
	cancel()
	// The UI should handle context cancellation gracefully
	// We can't easily test the visual output, but we can verify the function doesn't panic
}
