package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/context"
)

func main() {
	// Define command line flags
	var (
		searchType   = flag.String("type", "track", "Type of search: track, album, or playlist")
		searchQuery  = flag.String("query", "", "Search query")
		artistName   = flag.String("artist", "", "Artist name to filter results (only for track search)")
		limit        = flag.Int("limit", 5, "Number of results to display")
		showDetails  = flag.Bool("details", false, "Show detailed information about the results")
		interactive  = flag.Bool("interactive", false, "Run in interactive mode with a menu interface")
		returnToMenu = flag.Bool("return-to-menu", false, "Return to interactive menu after viewing search results")
	)

	// Add short flag alternatives
	flag.StringVar(searchType, "t", "track", "Short for -type")
	flag.StringVar(searchQuery, "q", "", "Short for -query")
	flag.StringVar(artistName, "a", "", "Short for -artist")
	flag.IntVar(limit, "l", 5, "Short for -limit")
	flag.BoolVar(showDetails, "d", false, "Short for -details")
	flag.BoolVar(interactive, "i", false, "Short for -interactive")
	flag.BoolVar(returnToMenu, "r", false, "Short for -return-to-menu")

	// Define usage information
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A CLI tool to search Spotify for tracks, albums, and playlists.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -type=track -query=\"Bohemian Rhapsody\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t track -q \"Bohemian Rhapsody\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type=track -query=\"Bohemian Rhapsody\" -artist=\"Queen\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t track -q \"Bohemian Rhapsody\" -a \"Queen\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type=album -query=\"Dark Side of the Moon\" -limit=3\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t album -q \"Dark Side of the Moon\" -l 3\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type=playlist -query=\"workout\" -details\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t playlist -q \"workout\" -d\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -interactive\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -i\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -query=\"Bohemian Rhapsody\" -return-to-menu\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -q \"Bohemian Rhapsody\" -r\n", os.Args[0])
	}

	flag.Parse()

	// Initialize Spotify client
	ctx := context.Background()
	client := getSpotifyClient(ctx)

	// Check if interactive mode is enabled
	if *interactive {
		// Start interactive menu
		menu := NewInteractiveMenu(ctx, client)
		if err := menu.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running interactive menu: %v\n", err)
		}
		return
	}

	// Validate search type for non-interactive mode
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

	// Validate search query for non-interactive mode
	if *searchQuery == "" {
		fmt.Fprintf(os.Stderr, "Error: missing search query\n")
		flag.Usage()
		return
	}

	// Perform search based on type and returnToMenu flag
	if *returnToMenu {
		switch *searchType {
		case "track":
			searchTracksWithMenu(ctx, client, *searchQuery, *artistName, *limit, *showDetails)
		case "album":
			searchAlbumsWithMenu(ctx, client, *searchQuery, *limit, *showDetails)
		case "playlist":
			searchPlaylistsWithMenu(ctx, client, *searchQuery, *limit, *showDetails)
		}
	} else {
		switch *searchType {
		case "track":
			searchTracks(ctx, client, *searchQuery, *artistName, *limit, *showDetails)
		case "album":
			searchAlbums(ctx, client, *searchQuery, *limit, *showDetails)
		case "playlist":
			searchPlaylists(ctx, client, *searchQuery, *limit, *showDetails)
		}
	}
}
