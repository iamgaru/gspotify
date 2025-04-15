// Package profile provides functionality for getting Spotify user profile information.
package profile

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/zmb3/spotify/v2"
)

var userID = flag.String("u", "", "the Spotify user ID to look up")

func init() {
	// Add long flag alternative (hidden from help)
	flag.StringVar(userID, "user", "", "")
}

// SpotifyClient interface for testing
type SpotifyClient interface {
	GetUsersPublicProfile(ctx context.Context, userID spotify.ID) (*spotify.User, error)
}

// For testing purposes
var getSpotifyClientFunc = defaultGetSpotifyClient

// GetProfile gets and displays the public profile information about a Spotify user.
func GetProfile() {
	flag.Parse()

	if *userID == "" {
		fmt.Fprintf(os.Stderr, "Error: missing user ID\n")
		flag.Usage()
		return
	}

	ctx := context.Background()
	client := getSpotifyClientFunc(ctx)
	displayProfile(ctx, client, *userID)
}

// defaultGetSpotifyClient creates a new Spotify client
func defaultGetSpotifyClient(ctx context.Context) SpotifyClient {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	return spotify.New(httpClient)
}

// displayProfile displays the user profile information
func displayProfile(ctx context.Context, client SpotifyClient, userID string) {
	user, err := client.GetUsersPublicProfile(ctx, spotify.ID(userID))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Println("User ID:", user.ID)
	fmt.Println("Display name:", user.DisplayName)
	fmt.Println("Spotify URI:", string(user.URI))
	fmt.Println("Endpoint:", user.Endpoint)
	fmt.Println("Followers:", user.Followers.Count)
}
