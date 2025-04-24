package config

import (
	"flag"
	"fmt"
	"os"
)

// SearchType represents the type of search to perform
type SearchType string

const (
	SearchTypeTrack    SearchType = "track"
	SearchTypeAlbum    SearchType = "album"
	SearchTypePlaylist SearchType = "playlist"
)

// Config holds all the command line configuration
type Config struct {
	SearchType   SearchType
	SearchQuery  string
	ArtistName   string
	Limit        int
	ShowDetails  bool
	Interactive  bool
	ReturnToMenu bool
	KeepPlaying  bool
	AutoPlay     bool
	StopPlayback bool
	UserID       string
}

// isSearchOperation returns true if the command requires a search operation
func (c *Config) isSearchOperation() bool {
	return !c.Interactive && !c.StopPlayback && c.UserID == ""
}

// ParseFlags parses command line flags and returns a Config struct
func ParseFlags() *Config {
	cfg := &Config{}

	// Define command line flags
	searchType := flag.String("t", "track", "Type of search: track, album, or playlist")
	flag.StringVar(&cfg.SearchQuery, "q", "", "Search query")
	flag.StringVar(&cfg.ArtistName, "a", "", "Artist name to filter results (only for track search)")
	flag.IntVar(&cfg.Limit, "l", 5, "Number of results to display")
	flag.BoolVar(&cfg.ShowDetails, "d", false, "Show detailed information about the results")
	flag.BoolVar(&cfg.Interactive, "i", false, "Run in interactive mode with a menu interface")
	flag.BoolVar(&cfg.ReturnToMenu, "r", false, "Return to interactive menu after viewing search results")
	flag.BoolVar(&cfg.KeepPlaying, "k", false, "Keep music playing when exiting the player interface")
	flag.BoolVar(&cfg.AutoPlay, "p", false, "Automatically play the first result and exit")
	flag.BoolVar(&cfg.StopPlayback, "s", false, "Stop the currently playing track")
	flag.StringVar(&cfg.UserID, "u", "", "Spotify user ID to look up profile information")

	// Add long flag alternatives
	flag.StringVar(searchType, "type", "track", "")
	flag.StringVar(&cfg.SearchQuery, "query", "", "")
	flag.StringVar(&cfg.ArtistName, "artist", "", "")
	flag.IntVar(&cfg.Limit, "limit", 5, "")
	flag.BoolVar(&cfg.ShowDetails, "details", false, "")
	flag.BoolVar(&cfg.Interactive, "interactive", false, "")
	flag.BoolVar(&cfg.ReturnToMenu, "return-to-menu", false, "")
	flag.BoolVar(&cfg.KeepPlaying, "keep-playing", false, "")
	flag.BoolVar(&cfg.AutoPlay, "auto-play", false, "")
	flag.BoolVar(&cfg.StopPlayback, "stop", false, "")
	flag.StringVar(&cfg.UserID, "user", "", "")

	// Set usage information
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A CLI tool to search and play Spotify tracks, albums, and playlists.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")

		flag.VisitAll(func(f *flag.Flag) {
			if f.Usage != "" {
				fmt.Fprintf(os.Stderr, "  -%s", f.Name)
				name, usage := flag.UnquoteUsage(f)
				if len(name) > 0 {
					fmt.Fprintf(os.Stderr, " %s", name)
				}
				if len(f.Name) <= 1 && f.DefValue == "false" {
					fmt.Fprintf(os.Stderr, "\t%s", usage)
				} else {
					fmt.Fprintf(os.Stderr, "\n\t%s", usage)
					if f.DefValue != "" && f.DefValue != "0" && f.DefValue != "false" {
						fmt.Fprintf(os.Stderr, " (default %v)", f.DefValue)
					}
				}
				fmt.Fprint(os.Stderr, "\n")
			}
		})

		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -t track -q \"Bohemian Rhapsody\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t track -q \"Bohemian Rhapsody\" -a \"Queen\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t album -q \"Dark Side of the Moon\" -l 3\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t playlist -q \"workout\" -d\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -i\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -q \"Bohemian Rhapsody\" -r\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -q \"Bohemian Rhapsody\" -k\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -q \"Bohemian Rhapsody\" -p\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -s\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -u spotify\n", os.Args[0])
	}

	flag.Parse()

	// Validate and set search type
	cfg.SearchType = SearchType(*searchType)
	if cfg.isSearchOperation() && !cfg.IsValidSearchType() {
		fmt.Fprintf(os.Stderr, "Error: invalid search type '%s'. Must be one of: track, album, playlist\n", *searchType)
		flag.Usage()
		os.Exit(1)
	}

	// Validate search query only when needed
	if cfg.isSearchOperation() && cfg.SearchQuery == "" {
		fmt.Fprintf(os.Stderr, "Error: missing search query\n")
		flag.Usage()
		os.Exit(1)
	}

	return cfg
}

// IsValidSearchType checks if the search type is valid
func (c *Config) IsValidSearchType() bool {
	switch c.SearchType {
	case SearchTypeTrack, SearchTypeAlbum, SearchTypePlaylist:
		return true
	default:
		return false
	}
}
