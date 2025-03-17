package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/context"
)

func main() {
	// Define command line flags
	searchType := flag.String("type", "track", "Type of search: track, album, or playlist")
	searchQuery := flag.String("query", "", "Search query")
	artistName := flag.String("artist", "", "Artist name to filter results (only for track search)")
	limit := flag.Int("limit", 5, "Number of results to display")
	showDetails := flag.Bool("details", false, "Show detailed information about the results")

	// Define usage information
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A CLI tool to search Spotify for tracks, albums, and playlists.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -type=track -query=\"Bohemian Rhapsody\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type=track -query=\"Bohemian Rhapsody\" -artist=\"Queen\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type=album -query=\"Dark Side of the Moon\" -limit=3\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type=playlist -query=\"workout\" -details\n", os.Args[0])
	}

	flag.Parse()

	// Validate search type
	validTypes := map[string]bool{
		"track":    true,
		"album":    true,
		"playlist": true,
	}

	if !validTypes[*searchType] {
		fmt.Fprintf(os.Stderr, "Error: invalid search type '%s'. Must be one of: track, album, playlist\n", *searchType)
		flag.Usage()
		return
	}

	// Validate search query
	if *searchQuery == "" {
		fmt.Fprintf(os.Stderr, "Error: missing search query\n")
		flag.Usage()
		return
	}

	// Initialize Spotify client
	ctx := context.Background()
	client := getSpotifyClient(ctx)

	// Perform search based on type
	switch *searchType {
	case "track":
		searchTracks(ctx, client, *searchQuery, *artistName, *limit, *showDetails)
	case "album":
		searchAlbums(ctx, client, *searchQuery, *limit, *showDetails)
	case "playlist":
		searchPlaylists(ctx, client, *searchQuery, *limit, *showDetails)
	}
}
