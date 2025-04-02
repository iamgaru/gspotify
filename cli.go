package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
		state := "gspotty-auth-" + fmt.Sprintf("%d", time.Now().UnixNano())

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
func searchTracks(ctx context.Context, client *spotify.Client, query string, artist string, limit int, showDetails bool, keepPlaying bool, autoPlay bool) {
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
		fmt.Println("No tracks found matching your query.")
		return
	}

	// Auto-play the first track if enabled
	if autoPlay && len(results.Tracks.Tracks) > 0 {
		fmt.Printf("Found %d tracks matching your query.\n", len(results.Tracks.Tracks))
		fmt.Printf("Auto-playing the first track: %s by %s\n",
			results.Tracks.Tracks[0].Name,
			joinArtistNames(results.Tracks.Tracks[0].Artists))

		playerUI := NewPlayerUI(ctx, client, results.Tracks.Tracks[0], keepPlaying, autoPlay)
		playerUI.SetSearchTracks(results.Tracks.Tracks)
		playerUI.Play()
		return
	}

	// Use the scrollable UI to display results
	ui := NewResultsUI("track", ctx, client, showDetails)
	ui.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
	ui.DisplayTrackResults(ctx, client, results.Tracks.Tracks)
}

// Helper function to join artist names for display
func joinArtistNames(artists []spotify.SimpleArtist) string {
	names := make([]string, len(artists))
	for i, artist := range artists {
		names[i] = artist.Name
	}
	return strings.Join(names, ", ")
}

// searchAlbums searches for albums and displays the results
func searchAlbums(ctx context.Context, client *spotify.Client, query string, limit int, showDetails bool, keepPlaying bool, autoPlay bool) {
	// Search for albums
	results, err := client.Search(ctx, query, spotify.SearchTypeAlbum, spotify.Limit(limit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for albums: %v\n", err)
		return
	}

	if results.Albums == nil || len(results.Albums.Albums) == 0 {
		fmt.Println("No albums found matching your query.")
		return
	}

	// Auto-play the first album track if enabled
	if autoPlay && len(results.Albums.Albums) > 0 {
		album := results.Albums.Albums[0]
		fmt.Printf("Found %d albums matching your query.\n", len(results.Albums.Albums))
		fmt.Printf("Selected the first album: %s by %s\n",
			album.Name,
			joinArtistNames(album.Artists))

		// Get the album's tracks
		albumTracks, err := client.GetAlbumTracks(ctx, album.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting album tracks: %v\n", err)
		} else if len(albumTracks.Tracks) > 0 {
			// Get the full track info for the first track
			fullTrack, err := client.GetTrack(ctx, albumTracks.Tracks[0].ID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting full track info: %v\n", err)
			} else {
				fmt.Printf("Auto-playing the first track: %s\n", fullTrack.Name)
				playerUI := NewPlayerUI(ctx, client, *fullTrack, keepPlaying, autoPlay)
				playerUI.Play()
				return
			}
		} else {
			fmt.Println("No tracks found in the selected album.")
		}
	}

	// Use the scrollable UI to display results
	ui := NewResultsUI("album", ctx, client, showDetails)
	ui.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
	ui.DisplayAlbumResults(ctx, client, results.Albums.Albums)
}

// searchPlaylists searches for playlists and displays the results
func searchPlaylists(ctx context.Context, client *spotify.Client, query string, limit int, showDetails bool, keepPlaying bool, autoPlay bool) {
	// Search for playlists
	results, err := client.Search(ctx, query, spotify.SearchTypePlaylist, spotify.Limit(limit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for playlists: %v\n", err)
		return
	}

	if results.Playlists == nil || len(results.Playlists.Playlists) == 0 {
		fmt.Println("No playlists found matching your query.")
		return
	}

	// Auto-play the first playlist track if enabled
	if autoPlay && len(results.Playlists.Playlists) > 0 {
		playlist := results.Playlists.Playlists[0]
		fmt.Printf("Found %d playlists matching your query.\n", len(results.Playlists.Playlists))
		fmt.Printf("Selected the first playlist: %s by %s\n",
			playlist.Name,
			playlist.Owner.DisplayName)

		// Get the playlist's tracks
		playlistTracks, err := client.GetPlaylistItems(ctx, playlist.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting playlist tracks: %v\n", err)
		} else if len(playlistTracks.Items) > 0 {
			// Check if the item contains a track (not an episode)
			if playlistTracks.Items[0].Track.Track != nil {
				track := playlistTracks.Items[0].Track.Track
				fmt.Printf("Auto-playing the first track: %s by %s\n",
					track.Name,
					joinArtistNames(track.Artists))
				playerUI := NewPlayerUI(ctx, client, *track, keepPlaying, autoPlay)
				playerUI.Play()
				return
			} else if playlistTracks.Items[0].Track.Episode != nil {
				// Episodes are not supported for playback in this application
				fmt.Println("The first item is an episode, which is not supported for playback. Showing playlist instead.")
			} else {
				fmt.Println("No playable tracks found in the selected playlist.")
			}
		} else {
			fmt.Println("No tracks found in the selected playlist.")
		}
	}

	// Use the scrollable UI to display results
	ui := NewResultsUI("playlist", ctx, client, showDetails)
	ui.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
	ui.DisplayPlaylistResults(ctx, client, results.Playlists.Playlists)
}

// searchTracksWithMenu searches for tracks and displays the results with return to menu option
func searchTracksWithMenu(ctx context.Context, client *spotify.Client, query string, artist string, limit int, showDetails bool, keepPlaying bool, autoPlay bool) {
	// Combine query and artist if artist is provided
	searchQuery := query
	if artist != "" {
		searchQuery = fmt.Sprintf("%s artist:%s", query, artist)
	}

	// Search for tracks
	results, err := client.Search(ctx, searchQuery, spotify.SearchTypeTrack, spotify.Limit(limit))
	if err != nil {
		fmt.Printf("Error searching for tracks: %v\n", err)
		return
	}

	if results.Tracks == nil || len(results.Tracks.Tracks) == 0 {
		fmt.Println("No tracks found matching your query.")
		return
	}

	// Auto-play the first track if enabled
	if autoPlay && len(results.Tracks.Tracks) > 0 {
		fmt.Printf("Found %d tracks matching your query.\n", len(results.Tracks.Tracks))
		fmt.Printf("Auto-playing the first track: %s by %s\n",
			results.Tracks.Tracks[0].Name,
			joinArtistNames(results.Tracks.Tracks[0].Artists))

		playerUI := NewPlayerUI(ctx, client, results.Tracks.Tracks[0], keepPlaying, autoPlay)
		playerUI.SetReturnToMenuFunction(func() {
			// Create and run a new instance of the interactive menu
			menu := NewInteractiveMenu(ctx, client)
			menu.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
			if err := menu.Run(); err != nil {
				fmt.Printf("Error running interactive menu: %v\n", err)
			}
		})
		playerUI.Play()
		return
	}

	// Create and run a new interactive menu
	ui := NewResultsUI("track", ctx, client, showDetails)
	ui.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
	ui.SetReturnToMenuFunction(func() {
		// Create and run a new instance of the interactive menu
		menu := NewInteractiveMenu(ctx, client)
		menu.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
		if err := menu.Run(); err != nil {
			fmt.Printf("Error running interactive menu: %v\n", err)
		}
	})
	ui.DisplayTrackResults(ctx, client, results.Tracks.Tracks)
}

// searchAlbumsWithMenu searches for albums and displays the results with return to menu option
func searchAlbumsWithMenu(ctx context.Context, client *spotify.Client, query string, limit int, showDetails bool, keepPlaying bool, autoPlay bool) {
	// Search for albums
	results, err := client.Search(ctx, query, spotify.SearchTypeAlbum, spotify.Limit(limit))
	if err != nil {
		fmt.Printf("Error searching for albums: %v\n", err)
		return
	}

	if results.Albums == nil || len(results.Albums.Albums) == 0 {
		fmt.Println("No albums found matching your query.")
		return
	}

	// Auto-play the first album track if enabled
	if autoPlay && len(results.Albums.Albums) > 0 {
		album := results.Albums.Albums[0]
		fmt.Printf("Found %d albums matching your query.\n", len(results.Albums.Albums))
		fmt.Printf("Selected the first album: %s by %s\n",
			album.Name,
			joinArtistNames(album.Artists))

		// Get the album's tracks
		albumTracks, err := client.GetAlbumTracks(ctx, album.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting album tracks: %v\n", err)
		} else if len(albumTracks.Tracks) > 0 {
			// Get the full track info for the first track
			fullTrack, err := client.GetTrack(ctx, albumTracks.Tracks[0].ID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting full track info: %v\n", err)
			} else {
				fmt.Printf("Auto-playing the first track: %s\n", fullTrack.Name)
				playerUI := NewPlayerUI(ctx, client, *fullTrack, keepPlaying, autoPlay)
				playerUI.SetReturnToMenuFunction(func() {
					// Create and run a new instance of the interactive menu
					menu := NewInteractiveMenu(ctx, client)
					menu.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
					if err := menu.Run(); err != nil {
						fmt.Printf("Error running interactive menu: %v\n", err)
					}
				})
				playerUI.Play()
				return
			}
		} else {
			fmt.Println("No tracks found in the selected album.")
		}
	}

	// Create and run a new interactive menu
	ui := NewResultsUI("album", ctx, client, showDetails)
	ui.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
	ui.SetReturnToMenuFunction(func() {
		// Create and run a new instance of the interactive menu
		menu := NewInteractiveMenu(ctx, client)
		menu.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
		if err := menu.Run(); err != nil {
			fmt.Printf("Error running interactive menu: %v\n", err)
		}
	})
	ui.DisplayAlbumResults(ctx, client, results.Albums.Albums)
}

// searchPlaylistsWithMenu searches for playlists and displays the results with return to menu option
func searchPlaylistsWithMenu(ctx context.Context, client *spotify.Client, query string, limit int, showDetails bool, keepPlaying bool, autoPlay bool) {
	// Search for playlists
	results, err := client.Search(ctx, query, spotify.SearchTypePlaylist, spotify.Limit(limit))
	if err != nil {
		fmt.Printf("Error searching for playlists: %v\n", err)
		return
	}

	if results.Playlists == nil || len(results.Playlists.Playlists) == 0 {
		fmt.Println("No playlists found matching your query.")
		return
	}

	// Auto-play the first playlist track if enabled
	if autoPlay && len(results.Playlists.Playlists) > 0 {
		playlist := results.Playlists.Playlists[0]
		fmt.Printf("Found %d playlists matching your query.\n", len(results.Playlists.Playlists))
		fmt.Printf("Selected the first playlist: %s by %s\n",
			playlist.Name,
			playlist.Owner.DisplayName)

		// Get the playlist's tracks
		playlistTracks, err := client.GetPlaylistItems(ctx, playlist.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting playlist tracks: %v\n", err)
		} else if len(playlistTracks.Items) > 0 {
			// Check if the item contains a track (not an episode)
			if playlistTracks.Items[0].Track.Track != nil {
				track := playlistTracks.Items[0].Track.Track
				fmt.Printf("Auto-playing the first track: %s by %s\n",
					track.Name,
					joinArtistNames(track.Artists))
				playerUI := NewPlayerUI(ctx, client, *track, keepPlaying, autoPlay)
				playerUI.SetReturnToMenuFunction(func() {
					// Create and run a new instance of the interactive menu
					menu := NewInteractiveMenu(ctx, client)
					menu.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
					if err := menu.Run(); err != nil {
						fmt.Printf("Error running interactive menu: %v\n", err)
					}
				})
				playerUI.Play()
				return
			} else if playlistTracks.Items[0].Track.Episode != nil {
				// Episodes are not supported for playback in this application
				fmt.Println("The first item is an episode, which is not supported for playback. Showing playlist instead.")
			} else {
				fmt.Println("No playable tracks found in the selected playlist.")
			}
		} else {
			fmt.Println("No tracks found in the selected playlist.")
		}
	}

	// Create and run a new interactive menu
	ui := NewResultsUI("playlist", ctx, client, showDetails)
	ui.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
	ui.SetReturnToMenuFunction(func() {
		// Create and run a new instance of the interactive menu
		menu := NewInteractiveMenu(ctx, client)
		menu.SetKeepPlayingFlag(keepPlaying) // Set the keep playing flag
		if err := menu.Run(); err != nil {
			fmt.Printf("Error running interactive menu: %v\n", err)
		}
	})
	ui.DisplayPlaylistResults(ctx, client, results.Playlists.Playlists)
}

// stopCurrentlyPlaying stops the currently playing track
func stopCurrentlyPlaying(ctx context.Context, client *spotify.Client) {
	// Get available devices first
	devices, err := client.PlayerDevices(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting devices: %v\n", err)
		return
	}

	// Check if there are any active devices
	if len(devices) == 0 {
		fmt.Println("No active Spotify devices found. Please open Spotify on any device first.")
		return
	}

	// Find an active device to use
	var deviceID spotify.ID
	for _, device := range devices {
		if device.Active {
			deviceID = device.ID
			break
		}
	}

	// If no active device found, use the first available one
	if deviceID == "" && len(devices) > 0 {
		deviceID = devices[0].ID
	}

	// Pause playback on the device
	err = client.PauseOpt(ctx, &spotify.PlayOptions{
		DeviceID: &deviceID,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error stopping playback: %v\n", err)
		return
	}

	fmt.Println("Playback stopped successfully.")
}
