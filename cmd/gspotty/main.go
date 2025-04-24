package main

import (
	"fmt"
	"os"

	"github.com/iamgaru/gspotty/internal/cli"
	"github.com/iamgaru/gspotty/internal/config"
	"github.com/iamgaru/gspotty/internal/env"
	"github.com/iamgaru/gspotty/internal/menu"
	"github.com/iamgaru/gspotty/internal/profile"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/net/context"
)

func main() {
	// Check environment variables first
	if err := env.CheckSpotifyCredentials(); err != nil {
		os.Exit(1)
	}

	// Parse command line flags
	cfg := config.ParseFlags()

	// Initialize Spotify client
	ctx := context.Background()
	client := cli.GetSpotifyClient(ctx)

	// Handle different command modes
	if err := handleCommand(ctx, client, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleCommand(ctx context.Context, client *spotify.Client, cfg *config.Config) error {
	// Handle user profile lookup
	if cfg.UserID != "" {
		profile.GetProfile(cfg.UserID)
		return nil
	}

	// Handle stop playback
	if cfg.StopPlayback {
		cli.StopCurrentlyPlaying(ctx, client)
		return nil
	}

	// Handle interactive mode
	if cfg.Interactive {
		interactiveMenu := menu.NewInteractiveMenu(ctx, client)
		interactiveMenu.SetKeepPlayingFlag(cfg.KeepPlaying)
		return interactiveMenu.Run()
	}

	// Handle search operations
	return handleSearch(ctx, client, cfg)
}

func handleSearch(ctx context.Context, client *spotify.Client, cfg *config.Config) error {
	// Perform search based on type and returnToMenu flag
	if cfg.ReturnToMenu {
		switch cfg.SearchType {
		case config.SearchTypeTrack:
			cli.SearchTracksWithMenu(ctx, client, cfg.SearchQuery, cfg.ArtistName, cfg.Limit, cfg.ShowDetails, cfg.KeepPlaying, cfg.AutoPlay)
		case config.SearchTypeAlbum:
			cli.SearchAlbumsWithMenu(ctx, client, cfg.SearchQuery, cfg.Limit, cfg.ShowDetails, cfg.KeepPlaying, cfg.AutoPlay)
		case config.SearchTypePlaylist:
			cli.SearchPlaylistsWithMenu(ctx, client, cfg.SearchQuery, cfg.Limit, cfg.ShowDetails, cfg.KeepPlaying, cfg.AutoPlay)
		}
	} else {
		switch cfg.SearchType {
		case config.SearchTypeTrack:
			cli.SearchTracks(ctx, client, cfg.SearchQuery, cfg.ArtistName, cfg.Limit, cfg.ShowDetails, cfg.KeepPlaying, cfg.AutoPlay)
		case config.SearchTypeAlbum:
			cli.SearchAlbums(ctx, client, cfg.SearchQuery, cfg.Limit, cfg.ShowDetails, cfg.KeepPlaying, cfg.AutoPlay)
		case config.SearchTypePlaylist:
			cli.SearchPlaylists(ctx, client, cfg.SearchQuery, cfg.Limit, cfg.ShowDetails, cfg.KeepPlaying, cfg.AutoPlay)
		}
	}
	return nil
}
