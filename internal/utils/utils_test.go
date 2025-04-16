package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
)

// TestFormatDuration tests the duration formatting function
func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		ms       int
		expected string
	}{
		{
			name:     "Zero duration",
			ms:       0,
			expected: "0:00",
		},
		{
			name:     "Less than a minute",
			ms:       45000,
			expected: "0:45",
		},
		{
			name:     "More than a minute",
			ms:       125000,
			expected: "2:05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.ms)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestJoinArtistNames tests the artist name joining function
func TestJoinArtistNames(t *testing.T) {
	tests := []struct {
		name     string
		artists  []spotify.SimpleArtist
		expected string
	}{
		{
			name:     "Empty artists",
			artists:  []spotify.SimpleArtist{},
			expected: "",
		},
		{
			name: "Single artist",
			artists: []spotify.SimpleArtist{
				{Name: "Artist 1"},
			},
			expected: "Artist 1",
		},
		{
			name: "Multiple artists",
			artists: []spotify.SimpleArtist{
				{Name: "Artist 1"},
				{Name: "Artist 2"},
				{Name: "Artist 3"},
			},
			expected: "Artist 1, Artist 2, Artist 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinArtistNames(tt.artists)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestOpenURL tests the URL opening function
func TestOpenURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "Invalid URL",
			url:     "not a url",
			wantErr: true,
		},
		{
			name:    "Relative URL",
			url:     "path/to/something",
			wantErr: true,
		},
		{
			name:    "Valid URL",
			url:     "https://example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidateSearchType tests the search type validation
func TestValidateSearchType(t *testing.T) {
	// Test valid search types
	assert.NotPanics(t, func() {
		ValidateSearchType("track")
	})

	// Test invalid search type
	assert.Panics(t, func() {
		ValidateSearchType("invalid")
	})
}

// TestValidateLimit tests the limit validation
func TestValidateLimit(t *testing.T) {
	// Test valid limit
	assert.NotPanics(t, func() {
		ValidateLimit(10)
	})

	// Test invalid limit
	assert.Panics(t, func() {
		ValidateLimit(0)
	})
}

// TestValidateSearchQuery tests the search query validation
func TestValidateSearchQuery(t *testing.T) {
	// Test valid query
	assert.NotPanics(t, func() {
		ValidateSearchQuery("test query")
	})

	// Test invalid query
	assert.Panics(t, func() {
		ValidateSearchQuery("")
	})
}
