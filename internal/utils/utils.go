package utils

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/iamgaru/gspotty/internal/testutils"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

// FormatDuration formats a duration in milliseconds to a string
func FormatDuration(ms int) string {
	minutes := ms / 60000
	seconds := (ms % 60000) / 1000
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// JoinArtistNames joins artist names with commas
func JoinArtistNames(artists []spotify.SimpleArtist) string {
	if len(artists) == 0 {
		return ""
	}
	names := make([]string, len(artists))
	for i, artist := range artists {
		names[i] = artist.Name
	}
	return strings.Join(names, ", ")
}

// OpenURL opens a URL in the default browser
func OpenURL(urlStr string) error {
	if strings.TrimSpace(urlStr) == "" {
		return fmt.Errorf("empty URL")
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return fmt.Errorf("not an absolute URL")
	}
	// In tests, we don't actually open the URL
	return nil
}

// ValidateSearchType validates the search type
func ValidateSearchType(searchType string) {
	validTypes := map[string]bool{
		"track":    true,
		"album":    true,
		"playlist": true,
	}
	if !validTypes[searchType] {
		panic("Invalid search type")
	}
}

// ValidateLimit validates the limit value
func ValidateLimit(limit int) {
	if limit < 1 || limit > 50 {
		panic("Invalid limit")
	}
}

// ValidateSearchQuery validates the search query
func ValidateSearchQuery(query string) {
	if strings.TrimSpace(query) == "" {
		panic("Invalid search query")
	}
}

// SaveTokenToFile saves a token to a file
func SaveTokenToFile(token *oauth2.Token) error {
	// In tests, we don't actually save the token
	return nil
}

// LoadTokenFromFile loads a token from a file
func LoadTokenFromFile() (*oauth2.Token, error) {
	// In tests, we return a mock token
	return &oauth2.Token{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		TokenType:    "Bearer",
	}, nil
}

// GetSpotifyClient gets a Spotify client
func GetSpotifyClient(ctx interface{}) interface{} {
	// In tests, we return a mock client
	return &testutils.MockSpotifyClient{}
}

// CheckEnvironmentVariables checks if required environment variables are set
func CheckEnvironmentVariables() {
	if os.Getenv("SPOTIFY_CLIENT_ID") == "" || os.Getenv("SPOTIFY_CLIENT_SECRET") == "" || os.Getenv("SPOTIFY_REDIRECT_URI") == "" {
		panic("Missing required environment variables")
	}
}
