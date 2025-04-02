package test

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
		{"Zero duration", 0, "0:00"},
		{"One minute", 60000, "1:00"},
		{"Two minutes", 120000, "2:00"},
		{"One minute thirty seconds", 90000, "1:30"},
		{"Two minutes fifteen seconds", 135000, "2:15"},
		{"Five minutes", 300000, "5:00"},
		{"Ten minutes", 600000, "10:00"},
		{"One hour", 3600000, "60:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.ms)
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
			name: "Single artist",
			artists: []spotify.SimpleArtist{
				{Name: "Queen"},
			},
			expected: "Queen",
		},
		{
			name: "Multiple artists",
			artists: []spotify.SimpleArtist{
				{Name: "Queen"},
				{Name: "David Bowie"},
			},
			expected: "Queen, David Bowie",
		},
		{
			name: "Three artists",
			artists: []spotify.SimpleArtist{
				{Name: "Queen"},
				{Name: "David Bowie"},
				{Name: "Freddie Mercury"},
			},
			expected: "Queen, David Bowie, Freddie Mercury",
		},
		{
			name:     "Empty artists",
			artists:  []spotify.SimpleArtist{},
			expected: "",
		},
		{
			name: "Artists with special characters",
			artists: []spotify.SimpleArtist{
				{Name: "Queen & Company"},
				{Name: "David Bowie (feat. Mick Jagger)"},
			},
			expected: "Queen & Company, David Bowie (feat. Mick Jagger)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinArtistNames(tt.artists)
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
		{"Valid URL", "https://www.spotify.com", false},
		{"Invalid URL", "not_a_url", true},
		{"Empty URL", "", true},
		{"Malformed URL", "://invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := openURL(tt.url)
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
	tests := []struct {
		name        string
		searchType  string
		shouldPanic bool
	}{
		{"Valid track type", "track", false},
		{"Valid album type", "album", false},
		{"Valid playlist type", "playlist", false},
		{"Invalid type", "invalid", true},
		{"Empty type", "", true},
		{"Case sensitive", "TRACK", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					validateSearchType(tt.searchType)
				})
			} else {
				assert.NotPanics(t, func() {
					validateSearchType(tt.searchType)
				})
			}
		})
	}
}

// TestValidateLimit tests the limit validation
func TestValidateLimit(t *testing.T) {
	tests := []struct {
		name        string
		limit       int
		shouldPanic bool
	}{
		{"Valid limit 1", 1, false},
		{"Valid limit 5", 5, false},
		{"Valid limit 50", 50, false},
		{"Invalid limit 0", 0, true},
		{"Invalid limit -1", -1, true},
		{"Invalid limit 51", 51, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					validateLimit(tt.limit)
				})
			} else {
				assert.NotPanics(t, func() {
					validateLimit(tt.limit)
				})
			}
		})
	}
}

// TestValidateSearchQuery tests the search query validation
func TestValidateSearchQuery(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		shouldPanic bool
	}{
		{"Valid query", "test query", false},
		{"Empty query", "", true},
		{"Whitespace query", "   ", true},
		{"Special characters", "test@#$%^&*()", false},
		{"Long query", "this is a very long search query that should still be valid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					validateSearchQuery(tt.query)
				})
			} else {
				assert.NotPanics(t, func() {
					validateSearchQuery(tt.query)
				})
			}
		})
	}
}
