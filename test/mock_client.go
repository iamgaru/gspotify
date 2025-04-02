package test

import (
	"context"
	"time"

	"github.com/zmb3/spotify/v2"
)

// MockSpotifyClient implements a mock Spotify client for testing
type MockSpotifyClient struct {
	playCalled        bool
	pauseCalled       bool
	nextCalled        bool
	previousCalled    bool
	seekCalled        bool
	volumeCalled      bool
	searchCalled      bool
	getTrackCalled    bool
	getAlbumCalled    bool
	getPlaylistCalled bool
}

// Play implements the Play method
func (m *MockSpotifyClient) Play(ctx context.Context) error {
	m.playCalled = true
	return nil
}

// Pause implements the Pause method
func (m *MockSpotifyClient) Pause(ctx context.Context) error {
	m.pauseCalled = true
	return nil
}

// Next implements the Next method
func (m *MockSpotifyClient) Next(ctx context.Context) error {
	m.nextCalled = true
	return nil
}

// Previous implements the Previous method
func (m *MockSpotifyClient) Previous(ctx context.Context) error {
	m.previousCalled = true
	return nil
}

// Seek implements the Seek method
func (m *MockSpotifyClient) Seek(ctx context.Context, position time.Duration) error {
	m.seekCalled = true
	return nil
}

// SetVolume implements the SetVolume method
func (m *MockSpotifyClient) SetVolume(ctx context.Context, volume int) error {
	m.volumeCalled = true
	return nil
}

// Search implements the Search method
func (m *MockSpotifyClient) Search(ctx context.Context, query string, searchType spotify.SearchType, opts ...spotify.RequestOption) (*spotify.SearchResult, error) {
	m.searchCalled = true
	return &spotify.SearchResult{
		Tracks: &spotify.FullTrackPage{
			Tracks: []spotify.FullTrack{
				{
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
				},
			},
		},
		Albums: &spotify.SimpleAlbumPage{
			Albums: []spotify.SimpleAlbum{
				{
					ID:   "test_album_id",
					Name: "Test Album",
					Artists: []spotify.SimpleArtist{
						{Name: "Test Artist"},
					},
				},
			},
		},
		Playlists: &spotify.SimplePlaylistPage{
			Playlists: []spotify.SimplePlaylist{
				{
					ID:   "test_playlist_id",
					Name: "Test Playlist",
					Owner: spotify.User{
						DisplayName: "Test User",
					},
				},
			},
		},
	}, nil
}

// GetTrack implements the GetTrack method
func (m *MockSpotifyClient) GetTrack(ctx context.Context, id spotify.ID) (*spotify.FullTrack, error) {
	m.getTrackCalled = true
	return &spotify.FullTrack{
		SimpleTrack: spotify.SimpleTrack{
			ID:   id,
			Name: "Test Track",
			Artists: []spotify.SimpleArtist{
				{Name: "Test Artist"},
			},
			Duration: 180000,
		},
		Album: spotify.SimpleAlbum{
			Name: "Test Album",
		},
	}, nil
}

// GetAlbum implements the GetAlbum method
func (m *MockSpotifyClient) GetAlbum(ctx context.Context, id spotify.ID) (*spotify.FullAlbum, error) {
	m.getAlbumCalled = true
	return &spotify.FullAlbum{
		SimpleAlbum: spotify.SimpleAlbum{
			ID:   id,
			Name: "Test Album",
			Artists: []spotify.SimpleArtist{
				{Name: "Test Artist"},
			},
		},
	}, nil
}

// GetPlaylist implements the GetPlaylist method
func (m *MockSpotifyClient) GetPlaylist(ctx context.Context, id spotify.ID) (*spotify.FullPlaylist, error) {
	m.getPlaylistCalled = true
	return &spotify.FullPlaylist{
		SimplePlaylist: spotify.SimplePlaylist{
			ID:   id,
			Name: "Test Playlist",
			Owner: spotify.User{
				DisplayName: "Test User",
			},
		},
		Tracks: spotify.PlaylistTrackPage{
			Tracks: []spotify.PlaylistTrack{
				{
					Track: spotify.FullTrack{
						SimpleTrack: spotify.SimpleTrack{
							ID:   "test_track_id",
							Name: "Test Track",
							Artists: []spotify.SimpleArtist{
								{Name: "Test Artist"},
							},
							Duration: 180000,
						},
					},
				},
			},
		},
	}, nil
}

// GetAlbumTracks implements the GetAlbumTracks method
func (m *MockSpotifyClient) GetAlbumTracks(ctx context.Context, id spotify.ID) (*spotify.SimpleTrackPage, error) {
	return &spotify.SimpleTrackPage{
		Tracks: []spotify.SimpleTrack{
			{
				ID:   "test_track_id",
				Name: "Test Track",
				Artists: []spotify.SimpleArtist{
					{Name: "Test Artist"},
				},
				Duration: 180000,
			},
		},
	}, nil
}
