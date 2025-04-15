package test

import (
	"context"
	"time"

	"github.com/zmb3/spotify/v2"
)

// MockSpotifyClient is a mock implementation of the Spotify client for testing
type MockSpotifyClient struct {
	searchCalled   bool
	playCalled     bool
	pauseCalled    bool
	nextCalled     bool
	previousCalled bool
	seekCalled     bool
	volumeCalled   bool
}

// Search mocks the Search method
func (m *MockSpotifyClient) Search(ctx context.Context, query string, t spotify.SearchType, opts ...spotify.RequestOption) (*spotify.SearchResult, error) {
	m.searchCalled = true
	// Return a mock search result
	return &spotify.SearchResult{}, nil
}

// Play mocks the Play method
func (m *MockSpotifyClient) Play(ctx context.Context) error {
	m.playCalled = true
	return nil
}

// Pause mocks the Pause method
func (m *MockSpotifyClient) Pause(ctx context.Context) error {
	m.pauseCalled = true
	return nil
}

// Next mocks the Next method
func (m *MockSpotifyClient) Next(ctx context.Context) error {
	m.nextCalled = true
	return nil
}

// Previous mocks the Previous method
func (m *MockSpotifyClient) Previous(ctx context.Context) error {
	m.previousCalled = true
	return nil
}

// Seek mocks the Seek method
func (m *MockSpotifyClient) Seek(ctx context.Context, position time.Duration) error {
	m.seekCalled = true
	return nil
}

// SetVolume mocks the SetVolume method
func (m *MockSpotifyClient) SetVolume(ctx context.Context, volume int) error {
	m.volumeCalled = true
	return nil
}
