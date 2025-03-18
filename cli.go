package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

// getSpotifyClient initializes and returns a Spotify client using client credentials flow
func getSpotifyClient(ctx context.Context) *spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     "4592901f2a854ff0bde6d5a348f29539",
		ClientSecret: "525f64edf0654150bc74434e7c6ee68c",
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("Error getting token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	return client
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
