package env

import (
	"fmt"
	"os"
)

// CheckSpotifyCredentials verifies that required Spotify API credentials are set
func CheckSpotifyCredentials() error {
	clientID := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")

	if clientID == "" || clientSecret == "" {
		fmt.Println("=================================================================")
		fmt.Println("ERROR: Spotify API credentials not properly configured")
		fmt.Println("=================================================================")

		if clientID == "" {
			fmt.Println("Missing SPOTIFY_ID environment variable")
		}

		if clientSecret == "" {
			fmt.Println("Missing SPOTIFY_SECRET environment variable")
		}

		fmt.Println("\nTo set up your credentials:")
		fmt.Println("1. Go to https://developer.spotify.com/dashboard/")
		fmt.Println("2. Log in and create a new app")
		fmt.Println("3. Set the redirect URI to http://localhost:8888/callback in your app settings")
		fmt.Println("4. Set these environment variables with your credentials:")
		fmt.Println("   export SPOTIFY_ID=your_client_id")
		fmt.Println("   export SPOTIFY_SECRET=your_client_secret")
		fmt.Println("=================================================================")
		return fmt.Errorf("spotify credentials not configured")
	}

	return nil
}
