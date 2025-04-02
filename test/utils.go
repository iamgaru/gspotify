package test

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

// formatDuration formats a duration in milliseconds to a string
func formatDuration(ms int) string {
	minutes := ms / 60000
	seconds := (ms % 60000) / 1000
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// joinArtistNames joins artist names with commas
func joinArtistNames(artists []spotify.SimpleArtist) string {
	if len(artists) == 0 {
		return ""
	}
	names := make([]string, len(artists))
	for i, artist := range artists {
		names[i] = artist.Name
	}
	return strings.Join(names, ", ")
}

// openURL opens a URL in the default browser
func openURL(urlStr string) error {
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

// validateSearchType validates the search type
func validateSearchType(searchType string) {
	validTypes := map[string]bool{
		"track":    true,
		"album":    true,
		"playlist": true,
	}
	if !validTypes[searchType] {
		panic("Invalid search type")
	}
}

// validateLimit validates the limit value
func validateLimit(limit int) {
	if limit < 1 || limit > 50 {
		panic("Invalid limit")
	}
}

// validateSearchQuery validates the search query
func validateSearchQuery(query string) {
	if strings.TrimSpace(query) == "" {
		panic("Invalid search query")
	}
}

// saveTokenToFile saves a token to a file
func saveTokenToFile(token *oauth2.Token) error {
	// In tests, we don't actually save the token
	return nil
}

// loadTokenFromFile loads a token from a file
func loadTokenFromFile() (*oauth2.Token, error) {
	// In tests, we return a mock token
	return &oauth2.Token{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		TokenType:    "Bearer",
	}, nil
}

// getSpotifyClient gets a Spotify client
func getSpotifyClient(ctx interface{}) interface{} {
	// In tests, we return a mock client
	return &MockSpotifyClient{}
}

// checkEnvironmentVariables checks if required environment variables are set
func checkEnvironmentVariables() {
	if os.Getenv("SPOTIFY_ID") == "" || os.Getenv("SPOTIFY_SECRET") == "" {
		panic("Missing required environment variables")
	}
}
