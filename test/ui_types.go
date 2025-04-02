package test

import (
	"context"
	"time"

	"github.com/zmb3/spotify/v2"
)

// PlayerUI represents the player interface
type PlayerUI struct {
	isPlaying         bool
	track             spotify.FullTrack
	totalDuration     time.Duration
	startTime         time.Time
	keepPlaying       bool
	isPlaylistMode    bool
	isSearchMode      bool
	isAlbumMode       bool
	currentTrackIndex int
	playlistTracks    []spotify.PlaylistTrack
	searchTracks      []spotify.FullTrack
	albumTracks       []spotify.SimpleTrack
}

// NewPlayerUI creates a new PlayerUI instance
func NewPlayerUI(ctx context.Context, client *MockSpotifyClient, track spotify.FullTrack, keepPlaying, showDetails bool) *PlayerUI {
	return &PlayerUI{
		track:         track,
		totalDuration: time.Duration(track.Duration) * time.Millisecond,
		keepPlaying:   keepPlaying,
	}
}

// ResultsUI represents the results interface
type ResultsUI struct {
	app          interface{}
	table        interface{}
	frame        interface{}
	showDetails  bool
	keepPlaying  bool
	returnToMenu func()
}

// NewResultsUI creates a new ResultsUI instance
func NewResultsUI(searchType string, ctx context.Context, client *MockSpotifyClient, showDetails bool) *ResultsUI {
	return &ResultsUI{
		showDetails: showDetails,
	}
}

// InteractiveMenu represents the interactive menu interface
type InteractiveMenu struct {
	app         interface{}
	pages       interface{}
	client      *MockSpotifyClient
	keepPlaying bool
}

// NewInteractiveMenu creates a new InteractiveMenu instance
func NewInteractiveMenu(ctx context.Context, client *MockSpotifyClient) *InteractiveMenu {
	return &InteractiveMenu{
		client: client,
	}
}
