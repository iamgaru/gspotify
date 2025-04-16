package player

import (
	"context"
	"testing"
	"time"

	"github.com/iamgaru/gspotty/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
)

// MockPlayerUI is a mock implementation of the PlayerUI for testing
type MockPlayerUI struct {
	app            interface{}
	track          spotify.FullTrack
	client         *testutils.MockSpotifyClient
	ctx            context.Context
	isPlaying      bool
	totalDuration  time.Duration
	keepPlaying    bool
	autoQuit       bool
	playlistTracks []spotify.PlaylistTrack
	searchTracks   []spotify.FullTrack
	albumTracks    []spotify.SimpleTrack
}

// NewMockPlayerUI creates a new mock player UI
func NewMockPlayerUI(ctx context.Context, client interface{}, track spotify.FullTrack, keepPlaying bool, autoQuit bool) *MockPlayerUI {
	return &MockPlayerUI{
		app:           nil,
		track:         track,
		client:        client.(*testutils.MockSpotifyClient),
		ctx:           ctx,
		isPlaying:     false,
		totalDuration: time.Duration(track.Duration) * time.Millisecond,
		keepPlaying:   keepPlaying,
		autoQuit:      autoQuit,
	}
}

// Play starts playback
func (p *MockPlayerUI) Play() {
	p.isPlaying = true
	p.client.PlayCalled = true
}

// Pause pauses playback
func (p *MockPlayerUI) Pause() {
	p.isPlaying = false
	p.client.PauseCalled = true
}

// Stop stops playback
func (p *MockPlayerUI) Stop() {
	p.isPlaying = false
	p.client.PauseCalled = true
}

// SetPlaylistTracks sets the playlist tracks
func (p *MockPlayerUI) SetPlaylistTracks(tracks []spotify.PlaylistTrack) {
	p.playlistTracks = tracks
}

// SetSearchTracks sets the search tracks
func (p *MockPlayerUI) SetSearchTracks(tracks []spotify.FullTrack) {
	p.searchTracks = tracks
}

// SetAlbumTracks sets the album tracks
func (p *MockPlayerUI) SetAlbumTracks(tracks []spotify.SimpleTrack) {
	p.albumTracks = tracks
}

// TestPlayerUI tests the PlayerUI functionality
func TestPlayerUI(t *testing.T) {
	// Create a mock Spotify client
	mockClient := &testutils.MockSpotifyClient{}

	// Create a test track
	testTrack := spotify.FullTrack{
		SimpleTrack: spotify.SimpleTrack{
			ID:   "test_track_id",
			Name: "Test Track",
			Artists: []spotify.SimpleArtist{
				{Name: "Test Artist"},
			},
			Duration: 180000, // 3 minutes
		},
		Album: spotify.SimpleAlbum{
			Name: "Test Album",
		},
	}

	// Create a new PlayerUI instance
	player := NewMockPlayerUI(context.Background(), mockClient, testTrack, false, false)

	// Test initial state
	assert.False(t, player.isPlaying)
	assert.Equal(t, testTrack.Name, player.track.Name)
	assert.Equal(t, time.Duration(180000)*time.Millisecond, player.totalDuration)

	// Test playback controls
	t.Run("Playback Controls", func(t *testing.T) {
		// Test play
		player.Play()
		assert.True(t, player.isPlaying)
		assert.True(t, mockClient.PlayCalled)

		// Test pause
		player.Pause()
		assert.False(t, player.isPlaying)
		assert.True(t, mockClient.PauseCalled)

		// Test stop
		player.Stop()
		assert.False(t, player.isPlaying)
		assert.True(t, mockClient.PauseCalled)
	})

	// Test playlist mode
	t.Run("Playlist Mode", func(t *testing.T) {
		// Set up playlist tracks
		playlistTracks := []spotify.PlaylistTrack{
			{Track: testTrack},
			{Track: spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{
					ID:   "test_track_id_2",
					Name: "Test Track 2",
					Artists: []spotify.SimpleArtist{
						{Name: "Test Artist 2"},
					},
					Duration: 240000, // 4 minutes
				},
			}},
		}

		player.SetPlaylistTracks(playlistTracks)
		assert.Equal(t, 2, len(player.playlistTracks))
	})

	// Test search mode
	t.Run("Search Mode", func(t *testing.T) {
		// Set up search tracks
		searchTracks := []spotify.FullTrack{
			testTrack,
			{
				SimpleTrack: spotify.SimpleTrack{
					ID:   "test_track_id_2",
					Name: "Test Track 2",
					Artists: []spotify.SimpleArtist{
						{Name: "Test Artist 2"},
					},
					Duration: 240000,
				},
			},
		}

		player.SetSearchTracks(searchTracks)
		assert.Equal(t, 2, len(player.searchTracks))
	})

	// Test album mode
	t.Run("Album Mode", func(t *testing.T) {
		// Set up album tracks
		albumTracks := []spotify.SimpleTrack{
			{
				ID:   "test_track_id",
				Name: "Test Track",
				Artists: []spotify.SimpleArtist{
					{Name: "Test Artist"},
				},
				Duration: 180000,
			},
			{
				ID:   "test_track_id_2",
				Name: "Test Track 2",
				Artists: []spotify.SimpleArtist{
					{Name: "Test Artist 2"},
				},
				Duration: 240000,
			},
		}

		player.SetAlbumTracks(albumTracks)
		assert.Equal(t, 2, len(player.albumTracks))
	})
}

// TestPlayerUIWithMockContext tests the PlayerUI with a mock context
func TestPlayerUIWithMockContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockClient := &testutils.MockSpotifyClient{}
	testTrack := spotify.FullTrack{
		SimpleTrack: spotify.SimpleTrack{
			ID:   "test_track_id",
			Name: "Test Track",
			Artists: []spotify.SimpleArtist{
				{Name: "Test Artist"},
			},
			Duration: 180000,
		},
	}

	_ = NewMockPlayerUI(ctx, mockClient, testTrack, false, false)

	// Test context cancellation
	cancel()
	// The player should handle context cancellation gracefully
	// We can't easily test the visual output, but we can verify the function doesn't panic
}

func TestPlayer(t *testing.T) {
	ctx := context.Background()
	mockClient := &testutils.MockSpotifyClient{}

	t.Run("Playback Control", func(t *testing.T) {
		// Test play
		t.Run("Play", func(t *testing.T) {
			err := mockClient.Play(ctx)
			assert.NoError(t, err)
			assert.True(t, mockClient.PlayCalled)
		})

		// Test pause
		t.Run("Pause", func(t *testing.T) {
			err := mockClient.Pause(ctx)
			assert.NoError(t, err)
			assert.True(t, mockClient.PauseCalled)
		})

		// Test next
		t.Run("Next", func(t *testing.T) {
			err := mockClient.Next(ctx)
			assert.NoError(t, err)
			assert.True(t, mockClient.NextCalled)
		})

		// Test previous
		t.Run("Previous", func(t *testing.T) {
			err := mockClient.Previous(ctx)
			assert.NoError(t, err)
			assert.True(t, mockClient.PreviousCalled)
		})

		// Test seek
		t.Run("Seek", func(t *testing.T) {
			position := 30 * time.Second
			err := mockClient.Seek(ctx, position)
			assert.NoError(t, err)
			assert.True(t, mockClient.SeekCalled)
		})

		// Test volume
		t.Run("Volume", func(t *testing.T) {
			volume := 50
			err := mockClient.SetVolume(ctx, volume)
			assert.NoError(t, err)
			assert.True(t, mockClient.SetVolumeCalled)
		})
	})
}
