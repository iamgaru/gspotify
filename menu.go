package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmb3/spotify/v2"
)

// InteractiveMenu represents the main menu UI
type InteractiveMenu struct {
	app    *tview.Application
	pages  *tview.Pages
	client *spotify.Client
	ctx    context.Context
}

// NewInteractiveMenu creates a new interactive menu
func NewInteractiveMenu(ctx context.Context, client *spotify.Client) *InteractiveMenu {
	app := tview.NewApplication()
	pages := tview.NewPages()

	menu := &InteractiveMenu{
		app:    app,
		pages:  pages,
		client: client,
		ctx:    ctx,
	}

	return menu
}

// Run starts the interactive menu
func (menu *InteractiveMenu) Run() error {
	// Create the main menu
	mainMenu := menu.createMainMenu()
	menu.pages.AddPage("main", mainMenu, true, true)

	// Set the root and run the application
	menu.app.SetRoot(menu.pages, true).EnableMouse(true)
	return menu.app.Run()
}

// createMainMenu creates the main menu form
func (menu *InteractiveMenu) createMainMenu() tview.Primitive {
	// Create a form
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Spotify Search").SetTitleAlign(tview.AlignCenter)

	// Add a dropdown for search type
	searchType := "track" // Default value
	form.AddDropDown("Search Type", []string{"track", "album", "playlist"}, 0, func(option string, optionIndex int) {
		searchType = option
	})

	// Add an input field for search query
	var searchQuery string
	form.AddInputField("Search Query", "", 40, nil, func(text string) {
		searchQuery = text
	})

	// Add an input field for artist name (only relevant for track search)
	var artistName string
	form.AddInputField("Artist Name (optional, for track search)", "", 40, nil, func(text string) {
		artistName = text
	})

	// Add an input field for limit
	var limitStr string
	form.AddInputField("Number of Results (1-50)", "5", 10, func(textToCheck string, lastChar rune) bool {
		// Only allow numbers
		if len(textToCheck) > 0 {
			_, err := strconv.Atoi(textToCheck)
			return err == nil && len(textToCheck) <= 2 // Max 2 digits
		}
		return true
	}, func(text string) {
		limitStr = text
	})

	// Add a checkbox for detailed results
	var showDetails bool
	form.AddCheckbox("Show Detailed Results", false, func(checked bool) {
		showDetails = checked
	})

	// Add buttons
	form.AddButton("Search", func() {
		// Validate search query
		if searchQuery == "" {
			menu.showError("Please enter a search query")
			return
		}

		// Parse limit
		limit := 5 // Default
		if limitStr != "" {
			var err error
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit < 1 || limit > 50 {
				menu.showError("Limit must be a number between 1 and 50")
				return
			}
		}

		// Perform search based on type
		switch searchType {
		case "track":
			menu.performTrackSearch(searchQuery, artistName, limit, showDetails)
		case "album":
			menu.performAlbumSearch(searchQuery, limit, showDetails)
		case "playlist":
			menu.performPlaylistSearch(searchQuery, limit, showDetails)
		}
	})

	form.AddButton("Quit", func() {
		menu.app.Stop()
	})

	// Add key capture for escape key to quit
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			menu.app.Stop()
			return nil
		}
		return event
	})

	// Create a frame with the form and add footer text about keyboard shortcuts
	frame := tview.NewFrame(form).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("ESC/Ctrl-C: Quit", false, tview.AlignCenter, tcell.ColorWhite)

	return frame
}

// showError displays an error message modal
func (menu *InteractiveMenu) showError(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			menu.pages.SwitchToPage("main")
		})

	menu.pages.AddPage("error", modal, true, true)
	menu.pages.SwitchToPage("error")
}

// performTrackSearch searches for tracks and displays the results
func (menu *InteractiveMenu) performTrackSearch(query, artist string, limit int, showDetails bool) {
	// Combine query and artist if artist is provided
	searchQuery := query
	if artist != "" {
		searchQuery = fmt.Sprintf("%s artist:%s", query, artist)
	}

	// Search for tracks
	results, err := menu.client.Search(menu.ctx, searchQuery, spotify.SearchTypeTrack, spotify.Limit(limit))
	if err != nil {
		menu.showError(fmt.Sprintf("Error searching for tracks: %v", err))
		return
	}

	if results.Tracks == nil || len(results.Tracks.Tracks) == 0 {
		menu.showError("No tracks found.")
		return
	}

	// Stop the current application
	menu.app.Stop()

	// Use the scrollable UI to display results
	ui := NewResultsUI("track", menu.ctx, menu.client, showDetails)

	// Set up the return to menu function
	ui.SetReturnToMenuFunction(func() {
		// Create and run a new instance of the interactive menu
		newMenu := NewInteractiveMenu(menu.ctx, menu.client)
		if err := newMenu.Run(); err != nil {
			fmt.Printf("Error running interactive menu: %v\n", err)
		}
	})

	ui.DisplayTrackResults(menu.ctx, menu.client, results.Tracks.Tracks)
}

// performAlbumSearch searches for albums and displays the results
func (menu *InteractiveMenu) performAlbumSearch(query string, limit int, showDetails bool) {
	// Search for albums
	results, err := menu.client.Search(menu.ctx, query, spotify.SearchTypeAlbum, spotify.Limit(limit))
	if err != nil {
		menu.showError(fmt.Sprintf("Error searching for albums: %v", err))
		return
	}

	if results.Albums == nil || len(results.Albums.Albums) == 0 {
		menu.showError("No albums found.")
		return
	}

	// Stop the current application
	menu.app.Stop()

	// Use the scrollable UI to display results
	ui := NewResultsUI("album", menu.ctx, menu.client, showDetails)

	// Set up the return to menu function
	ui.SetReturnToMenuFunction(func() {
		// Create and run a new instance of the interactive menu
		newMenu := NewInteractiveMenu(menu.ctx, menu.client)
		if err := newMenu.Run(); err != nil {
			fmt.Printf("Error running interactive menu: %v\n", err)
		}
	})

	ui.DisplayAlbumResults(menu.ctx, menu.client, results.Albums.Albums)
}

// performPlaylistSearch searches for playlists and displays the results
func (menu *InteractiveMenu) performPlaylistSearch(query string, limit int, showDetails bool) {
	// Search for playlists
	results, err := menu.client.Search(menu.ctx, query, spotify.SearchTypePlaylist, spotify.Limit(limit))
	if err != nil {
		menu.showError(fmt.Sprintf("Error searching for playlists: %v", err))
		return
	}

	if results.Playlists == nil || len(results.Playlists.Playlists) == 0 {
		menu.showError("No playlists found.")
		return
	}

	// Stop the current application
	menu.app.Stop()

	// Use the scrollable UI to display results
	ui := NewResultsUI("playlist", menu.ctx, menu.client, showDetails)

	// Set up the return to menu function
	ui.SetReturnToMenuFunction(func() {
		// Create and run a new instance of the interactive menu
		newMenu := NewInteractiveMenu(menu.ctx, menu.client)
		if err := newMenu.Run(); err != nil {
			fmt.Printf("Error running interactive menu: %v\n", err)
		}
	})

	ui.DisplayPlaylistResults(menu.ctx, menu.client, results.Playlists.Playlists)
}
