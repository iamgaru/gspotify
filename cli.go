package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

// Define global variables for authentication
const (
	redirectURI = "http://localhost:8888/callback"
	tokenFile   = ".spotify_token.json"
)

// TokenInfo stores authentication tokens
type TokenInfo struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	Expiry       time.Time `json:"expiry"`
}

// getSpotifyClient initializes and returns a Spotify client with proper authentication for playback
func getSpotifyClient(ctx context.Context) *spotify.Client {
	clientID := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")

	// Check if environment variables are set
	if clientID == "" || clientSecret == "" {
		log.Fatalf("Error: SPOTIFY_ID and SPOTIFY_SECRET environment variables must be set\n" +
			"Please set them using:\n" +
			"export SPOTIFY_ID=your_client_id\n" +
			"export SPOTIFY_SECRET=your_client_secret")
	}

	// Set up authentication with required scopes for playback control
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
		),
	)

	// Try to load token from file
	token, err := loadTokenFromFile()
	if err != nil || token.AccessToken == "" || token.RefreshToken == "" || time.Now().After(token.Expiry) {
		// If no valid token, we need to perform an initial authorization
		if token.RefreshToken != "" {
			// Try to refresh the token first if we have a refresh token
			oauthToken := &oauth2.Token{
				RefreshToken: token.RefreshToken,
				TokenType:    token.TokenType,
				Expiry:       token.Expiry,
			}

			// Create a new client with the refresh token
			client := spotify.New(auth.Client(ctx, oauthToken))

			// The client will automatically refresh the token when needed
			// We can return it directly
			return client
		}

		// We need to do a one-time interactive login
		fmt.Println("You need to authorize this application to control Spotify.")
		fmt.Println("This is a one-time process. After authorization, you won't need to do this again.")
		// Generate a random state string for security
		state := "gspotify-auth-" + fmt.Sprintf("%d", time.Now().UnixNano())

		// Generate the auth URL
		authURL := auth.AuthURL(state)

		// Try to open the URL in the default browser
		fmt.Println("Opening the authorization page in your default browser...")
		openErr := openURL(authURL)
		if openErr != nil {
			// Fall back to displaying the URL if opening fails
			fmt.Printf("Could not open browser automatically. Please visit this URL manually: %s\n", authURL)
		} else {
			fmt.Println("Browser opened. Please complete the authorization in your browser.")
			fmt.Println("Waiting for callback from Spotify...")
		}

		// Set up temporary HTTP server to handle the callback
		ch := make(chan *oauth2.Token)
		errCh := make(chan error)

		// Create a server with timeouts to prevent hanging
		server := &http.Server{
			Addr:         ":8888",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		}

		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			// Check for error parameter
			if errParam := r.URL.Query().Get("error"); errParam != "" {
				errCh <- fmt.Errorf("Spotify authorization error: %s", errParam)
				fmt.Fprintf(w, "Authorization failed: %s. Please close this window and try again.", errParam)
				return
			}

			// Get state and code from the request
			receivedState := r.URL.Query().Get("state")
			if receivedState != state {
				errCh <- fmt.Errorf("state mismatch: expected %s, got %s", state, receivedState)
				http.Error(w, "State mismatch error", http.StatusBadRequest)
				return
			}

			// Attempt to get token
			token, err := auth.Token(r.Context(), state, r)
			if err != nil {
				errCh <- fmt.Errorf("failed to get token: %v", err)
				http.Error(w, "Failed to get token", http.StatusInternalServerError)
				return
			}

			// Send the token to the channel
			ch <- token
			fmt.Fprintf(w, "Authorization successful! You can close this window and return to the application.")
		})

		// Start the server in a goroutine
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("server error: %v", err)
			}
		}()

		// Set a timeout for the auth process
		authTimeout := time.After(2 * time.Minute)

		// Wait for the token, error, or timeout
		var token *oauth2.Token
		select {
		case token = <-ch:
			fmt.Println("Authorization successful!")
		case err := <-errCh:
			server.Shutdown(ctx)
			log.Fatalf("Authorization failed: %v", err)
		case <-authTimeout:
			server.Shutdown(ctx)
			log.Fatalf("Authorization timed out. Please try again.")
		}

		// Save token to file for future use
		if err := saveTokenToFile(token); err != nil {
			fmt.Printf("Warning: Failed to save token: %v\n", err)
		}

		// Shutdown server
		server.Shutdown(ctx)

		return spotify.New(auth.Client(ctx, token))
	}

	// Create OAuth2 token from stored token
	oauthToken := &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	}

	// Return client with valid token
	return spotify.New(auth.Client(ctx, oauthToken))
}

// loadTokenFromFile loads authentication token from file
func loadTokenFromFile() (TokenInfo, error) {
	var token TokenInfo

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return token, err
	}

	tokenPath := filepath.Join(homeDir, tokenFile)
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(data, &token)
	return token, err
}

// saveTokenToFile saves authentication token to file
func saveTokenToFile(token *oauth2.Token) error {
	if token == nil {
		return fmt.Errorf("no token to save")
	}

	tokenInfo := TokenInfo{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	}

	data, err := json.Marshal(tokenInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	tokenPath := filepath.Join(homeDir, tokenFile)

	// Ensure the file permissions are restricted to current user only
	err = os.WriteFile(tokenPath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write token file: %v", err)
	}

	fmt.Printf("Token successfully saved to %s\n", tokenPath)
	return nil
}

// searchTracks searches for tracks and displays the results
func searchTracks(ctx context.Context, client *spotify.Client, query string, artist string, limit int, showDetails bool) {
	// Combine query and artist if artist is provided
	searchQuery := query
	if artist != "" {
		searchQuery = fmt.Sprintf("%s artist:%s", query, artist)
	}

	// Search for tracks
	results, err := client.Search(ctx, searchQuery, spotify.SearchTypeTrack, spotify.Limit(limit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for tracks: %v\n", err)
		return
	}

	if results.Tracks == nil || len(results.Tracks.Tracks) == 0 {
		fmt.Println("No tracks found.")
		return
	}

	// Use the scrollable UI to display results
	ui := NewResultsUI("track", ctx, client, showDetails)
	ui.DisplayTrackResults(ctx, client, results.Tracks.Tracks)
}

// searchAlbums searches for albums and displays the results
func searchAlbums(ctx context.Context, client *spotify.Client, query string, limit int, showDetails bool) {
	// Search for albums
	results, err := client.Search(ctx, query, spotify.SearchTypeAlbum, spotify.Limit(limit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for albums: %v\n", err)
		return
	}

	if results.Albums == nil || len(results.Albums.Albums) == 0 {
		fmt.Println("No albums found.")
		return
	}

	// Use the scrollable UI to display results
	ui := NewResultsUI("album", ctx, client, showDetails)
	ui.DisplayAlbumResults(ctx, client, results.Albums.Albums)
}

// searchPlaylists searches for playlists and displays the results
func searchPlaylists(ctx context.Context, client *spotify.Client, query string, limit int, showDetails bool) {
	// Search for playlists
	results, err := client.Search(ctx, query, spotify.SearchTypePlaylist, spotify.Limit(limit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for playlists: %v\n", err)
		return
	}

	if results.Playlists == nil || len(results.Playlists.Playlists) == 0 {
		fmt.Println("No playlists found.")
		return
	}

	// Use the scrollable UI to display results
	ui := NewResultsUI("playlist", ctx, client, showDetails)
	ui.DisplayPlaylistResults(ctx, client, results.Playlists.Playlists)
}

// searchTracksWithMenu searches for tracks and displays the results with return to menu option
func searchTracksWithMenu(ctx context.Context, client *spotify.Client, query string, artist string, limit int, showDetails bool) {
	// Combine query and artist if artist is provided
	searchQuery := query
	if artist != "" {
		searchQuery = fmt.Sprintf("%s artist:%s", query, artist)
	}

	// Search for tracks
	results, err := client.Search(ctx, searchQuery, spotify.SearchTypeTrack, spotify.Limit(limit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for tracks: %v\n", err)
		return
	}

	if results.Tracks == nil || len(results.Tracks.Tracks) == 0 {
		fmt.Println("No tracks found.")
		return
	}

	// Use the scrollable UI to display results
	ui := NewResultsUI("track", ctx, client, showDetails)

	// Set up the return to menu function
	ui.SetReturnToMenuFunction(func() {
		// Create and run a new instance of the interactive menu
		menu := NewInteractiveMenu(ctx, client)
		if err := menu.Run(); err != nil {
			fmt.Printf("Error running interactive menu: %v\n", err)
		}
	})

	ui.DisplayTrackResults(ctx, client, results.Tracks.Tracks)
}

// searchAlbumsWithMenu searches for albums and displays the results with return to menu option
func searchAlbumsWithMenu(ctx context.Context, client *spotify.Client, query string, limit int, showDetails bool) {
	// Search for albums
	results, err := client.Search(ctx, query, spotify.SearchTypeAlbum, spotify.Limit(limit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for albums: %v\n", err)
		return
	}

	if results.Albums == nil || len(results.Albums.Albums) == 0 {
		fmt.Println("No albums found.")
		return
	}

	// Use the scrollable UI to display results
	ui := NewResultsUI("album", ctx, client, showDetails)

	// Set up the return to menu function
	ui.SetReturnToMenuFunction(func() {
		// Create and run a new instance of the interactive menu
		menu := NewInteractiveMenu(ctx, client)
		if err := menu.Run(); err != nil {
			fmt.Printf("Error running interactive menu: %v\n", err)
		}
	})

	ui.DisplayAlbumResults(ctx, client, results.Albums.Albums)
}

// searchPlaylistsWithMenu searches for playlists and displays the results with return to menu option
func searchPlaylistsWithMenu(ctx context.Context, client *spotify.Client, query string, limit int, showDetails bool) {
	// Search for playlists
	results, err := client.Search(ctx, query, spotify.SearchTypePlaylist, spotify.Limit(limit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for playlists: %v\n", err)
		return
	}

	if results.Playlists == nil || len(results.Playlists.Playlists) == 0 {
		fmt.Println("No playlists found.")
		return
	}

	// Use the scrollable UI to display results
	ui := NewResultsUI("playlist", ctx, client, showDetails)

	// Set up the return to menu function
	ui.SetReturnToMenuFunction(func() {
		// Create and run a new instance of the interactive menu
		menu := NewInteractiveMenu(ctx, client)
		if err := menu.Run(); err != nil {
			fmt.Printf("Error running interactive menu: %v\n", err)
		}
	})

	ui.DisplayPlaylistResults(ctx, client, results.Playlists.Playlists)
}
