package main

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmb3/spotify/v2"
)

// openURL opens the specified URL in the default browser
func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

// ResultsUI represents a scrollable UI for displaying search results
type ResultsUI struct {
	app         *tview.Application
	table       *tview.Table
	frame       *tview.Frame
	results     interface{}
	resultType  string
	client      *spotify.Client
	ctx         context.Context
	showDetails bool
}

// NewResultsUI creates a new scrollable UI for displaying search results
func NewResultsUI(resultType string, ctx context.Context, client *spotify.Client, showDetails bool) *ResultsUI {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	ui := &ResultsUI{
		app:         app,
		table:       table,
		resultType:  resultType,
		client:      client,
		ctx:         ctx,
		showDetails: showDetails,
	}

	// Set up key bindings
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			app.Stop()
		case tcell.KeyEnter:
			// Get the selected row
			row, _ := table.GetSelection()
			ui.displayDetails(row)
			return nil
		}
		return event
	})

	// Set up mouse handling for clicking on Spotify links
	table.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick {
			row, col := table.GetSelection()
			// Check if the click is on the Spotify Link column (column 5)
			if row > 0 && col == 5 {
				// Get the Spotify link from the cell
				cell := table.GetCell(row, col)
				if cell != nil {
					spotifyLink := cell.Text
					if spotifyLink != "" {
						// Try to open the link
						err := openURL(spotifyLink)
						if err != nil {
							// If opening fails, show an error message
							infoModal := tview.NewModal().
								SetText(fmt.Sprintf("Could not open browser automatically.\nSpotify link: %s", spotifyLink)).
								AddButtons([]string{"OK"}).
								SetDoneFunc(func(buttonIndex int, buttonLabel string) {
									ui.app.SetRoot(ui.frame, true)
								})
							ui.app.SetRoot(infoModal, true)
						} else {
							// Show a confirmation
							infoModal := tview.NewModal().
								SetText(fmt.Sprintf("Opening in browser:\n%s", spotifyLink)).
								AddButtons([]string{"OK"}).
								SetDoneFunc(func(buttonIndex int, buttonLabel string) {
									ui.app.SetRoot(ui.frame, true)
								})
							ui.app.SetRoot(infoModal, true)
						}
						return tview.MouseLeftClick, nil // Consume the event
					}
				}
			}
		}
		return action, event
	})

	return ui
}

// DisplayTrackResults displays track search results in a scrollable UI
func (ui *ResultsUI) DisplayTrackResults(ctx context.Context, client *spotify.Client, tracks []spotify.FullTrack) {
	ui.results = tracks

	// Set up table headers
	headers := []string{"ID", "Track Name", "Artist", "Album", "Popularity", "Spotify Link", "URI"}
	for i, header := range headers {
		ui.table.SetCell(0, i, tview.NewTableCell(header).SetSelectable(false).SetAttributes(tcell.AttrBold))
	}

	// Populate table with track data
	for i, track := range tracks {
		row := i + 1 // +1 for header row

		artists := make([]string, len(track.Artists))
		for j, artist := range track.Artists {
			artists[j] = artist.Name
		}

		// Create Spotify web link
		spotifyLink := fmt.Sprintf("https://open.spotify.com/track/%s", track.ID)

		// Set cell values
		ui.table.SetCell(row, 0, tview.NewTableCell(string(track.ID)))
		ui.table.SetCell(row, 1, tview.NewTableCell(track.Name))
		ui.table.SetCell(row, 2, tview.NewTableCell(strings.Join(artists, ", ")))
		ui.table.SetCell(row, 3, tview.NewTableCell(track.Album.Name))
		ui.table.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%d", track.Popularity)))
		ui.table.SetCell(row, 5, tview.NewTableCell(spotifyLink).SetTextColor(tcell.ColorBlue))
		ui.table.SetCell(row, 6, tview.NewTableCell(string(track.URI)))
	}

	// Set up the layout
	ui.setupLayout("Track Search Results")
}

// DisplayAlbumResults displays album search results in a scrollable UI
func (ui *ResultsUI) DisplayAlbumResults(ctx context.Context, client *spotify.Client, albums []spotify.SimpleAlbum) {
	ui.results = albums

	// Set up table headers
	headers := []string{"ID", "Album Name", "Artist", "Release Date", "Total Tracks", "Spotify Link", "URI"}
	for i, header := range headers {
		ui.table.SetCell(0, i, tview.NewTableCell(header).SetSelectable(false).SetAttributes(tcell.AttrBold))
	}

	// Populate table with album data
	for i, album := range albums {
		row := i + 1 // +1 for header row

		artists := make([]string, len(album.Artists))
		for j, artist := range album.Artists {
			artists[j] = artist.Name
		}

		// Create Spotify web link
		spotifyLink := fmt.Sprintf("https://open.spotify.com/album/%s", album.ID)

		// Set cell values
		ui.table.SetCell(row, 0, tview.NewTableCell(string(album.ID)))
		ui.table.SetCell(row, 1, tview.NewTableCell(album.Name))
		ui.table.SetCell(row, 2, tview.NewTableCell(strings.Join(artists, ", ")))
		ui.table.SetCell(row, 3, tview.NewTableCell(album.ReleaseDate))
		ui.table.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%d", album.TotalTracks)))
		ui.table.SetCell(row, 5, tview.NewTableCell(spotifyLink).SetTextColor(tcell.ColorBlue))
		ui.table.SetCell(row, 6, tview.NewTableCell(string(album.URI)))
	}

	// Set up the layout
	ui.setupLayout("Album Search Results")
}

// DisplayPlaylistResults displays playlist search results in a scrollable UI
func (ui *ResultsUI) DisplayPlaylistResults(ctx context.Context, client *spotify.Client, playlists []spotify.SimplePlaylist) {
	ui.results = playlists

	// Set up table headers
	headers := []string{"ID", "Playlist Name", "Owner", "Total Tracks", "Spotify Link", "URI"}
	for i, header := range headers {
		ui.table.SetCell(0, i, tview.NewTableCell(header).SetSelectable(false).SetAttributes(tcell.AttrBold))
	}

	// Populate table with playlist data
	for i, playlist := range playlists {
		row := i + 1 // +1 for header row

		// Create Spotify web link
		spotifyLink := fmt.Sprintf("https://open.spotify.com/playlist/%s", playlist.ID)

		// Set cell values
		ui.table.SetCell(row, 0, tview.NewTableCell(string(playlist.ID)))
		ui.table.SetCell(row, 1, tview.NewTableCell(playlist.Name))
		ui.table.SetCell(row, 2, tview.NewTableCell(playlist.Owner.DisplayName))
		ui.table.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%d", playlist.Tracks.Total)))
		ui.table.SetCell(row, 4, tview.NewTableCell(spotifyLink).SetTextColor(tcell.ColorBlue))
		ui.table.SetCell(row, 5, tview.NewTableCell(string(playlist.URI)))
	}

	// Set up the layout
	ui.setupLayout("Playlist Search Results")
}

// setupLayout sets up the UI layout
func (ui *ResultsUI) setupLayout(title string) {
	// Create a frame to hold the table
	ui.frame = tview.NewFrame(ui.table).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(title, true, tview.AlignCenter, tcell.ColorWhite).
		AddText("↑/↓: Navigate • Enter: Show Details • Click on Spotify Link to Open • ESC/Ctrl-C: Exit", false, tview.AlignCenter, tcell.ColorWhite)

	// Set the root and run the application
	ui.app.SetRoot(ui.frame, true).EnableMouse(true)
	if err := ui.app.Run(); err != nil {
		panic(err)
	}
}

// displayDetails shows detailed information about the selected item
func (ui *ResultsUI) displayDetails(row int) {
	if row == 0 { // Header row
		return
	}

	// Create a modal for displaying details
	modal := tview.NewModal()

	var text string
	var spotifyLink string
	var spotifyURI string

	switch ui.resultType {
	case "track":
		tracks := ui.results.([]spotify.FullTrack)
		if row-1 < len(tracks) {
			track := tracks[row-1]
			artists := make([]string, len(track.Artists))
			for i, artist := range track.Artists {
				artists[i] = artist.Name
			}

			spotifyLink = fmt.Sprintf("https://open.spotify.com/track/%s", track.ID)
			spotifyURI = string(track.URI)

			text = fmt.Sprintf("Track: %s\nArtist(s): %s\nAlbum: %s\nPopularity: %d\nSpotify Link: %s\nURI: %s",
				track.Name,
				strings.Join(artists, ", "),
				track.Album.Name,
				track.Popularity,
				spotifyLink,
				track.URI)

			// Add audio features if showDetails is true
			if ui.showDetails {
				features, err := ui.client.GetAudioFeatures(ui.ctx, track.ID)
				if err == nil && len(features) > 0 {
					text += fmt.Sprintf("\n\nAudio Features:\nEnergy: %.2f\nDanceability: %.2f\nValence: %.2f\nTempo: %.2f",
						features[0].Energy,
						features[0].Danceability,
						features[0].Valence,
						features[0].Tempo)
				}
			}
		}
	case "album":
		albums := ui.results.([]spotify.SimpleAlbum)
		if row-1 < len(albums) {
			album := albums[row-1]
			artists := make([]string, len(album.Artists))
			for i, artist := range album.Artists {
				artists[i] = artist.Name
			}

			spotifyLink = fmt.Sprintf("https://open.spotify.com/album/%s", album.ID)
			spotifyURI = string(album.URI)

			text = fmt.Sprintf("Album: %s\nArtist(s): %s\nRelease Date: %s\nTotal Tracks: %d\nSpotify Link: %s\nURI: %s",
				album.Name,
				strings.Join(artists, ", "),
				album.ReleaseDate,
				album.TotalTracks,
				spotifyLink,
				album.URI)

			// Get album tracks
			albumTracks, err := ui.client.GetAlbumTracks(ui.ctx, album.ID)
			if err == nil && albumTracks != nil && len(albumTracks.Tracks) > 0 {
				// Create a list view for tracks
				trackList := tview.NewList().
					SetMainTextColor(tcell.ColorWhite).
					SetSelectedTextColor(tcell.ColorBlack).
					SetSelectedBackgroundColor(tcell.ColorGreen)

				// Add tracks to the list
				for i, track := range albumTracks.Tracks {
					trackSpotifyLink := fmt.Sprintf("https://open.spotify.com/track/%s", track.ID)
					trackURI := string(track.URI)

					// Create a closure to capture the current track's link and URI
					trackList.AddItem(fmt.Sprintf("%d. %s", i+1, track.Name),
						fmt.Sprintf("Duration: %s", formatDuration(track.Duration)),
						rune('1'+i),
						func(trackLink, trackUri string) func() {
							return func() {
								// Try to open the link in the default browser
								err := openURL(trackLink)
								if err != nil {
									// If opening the browser fails, just show the link
									infoModal := tview.NewModal().
										SetText(fmt.Sprintf("Could not open browser automatically.\nSpotify link: %s\nURI: %s", trackLink, trackUri)).
										AddButtons([]string{"OK"}).
										SetDoneFunc(func(buttonIndex int, buttonLabel string) {
											ui.app.SetRoot(trackList, true)
										})
									ui.app.SetRoot(infoModal, true)
								} else {
									// Show a confirmation that the link was opened
									infoModal := tview.NewModal().
										SetText(fmt.Sprintf("Opening in browser:\n%s", trackLink)).
										AddButtons([]string{"OK"}).
										SetDoneFunc(func(buttonIndex int, buttonLabel string) {
											ui.app.SetRoot(trackList, true)
										})
									ui.app.SetRoot(infoModal, true)
								}
							}
						}(trackSpotifyLink, trackURI))
				}

				// Add a "Back" option at the end of the list
				trackList.AddItem("Back to Album Details", "", 'b', func() {
					// Create a modal with album details
					detailsModal := tview.NewModal().
						SetText(text).
						AddButtons([]string{"Open Album in Spotify", "View Tracks", "Close"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							if buttonIndex == 0 && spotifyLink != "" {
								// Open album link
								err := openURL(spotifyLink)
								if err != nil {
									infoModal := tview.NewModal().
										SetText(fmt.Sprintf("Could not open browser automatically.\nSpotify link: %s\nURI: %s", spotifyLink, spotifyURI)).
										AddButtons([]string{"OK"}).
										SetDoneFunc(func(buttonIndex int, buttonLabel string) {
											ui.app.SetRoot(ui.frame, true)
										})
									ui.app.SetRoot(infoModal, true)
								} else {
									infoModal := tview.NewModal().
										SetText(fmt.Sprintf("Opening in browser:\n%s", spotifyLink)).
										AddButtons([]string{"OK"}).
										SetDoneFunc(func(buttonIndex int, buttonLabel string) {
											ui.app.SetRoot(ui.frame, true)
										})
									ui.app.SetRoot(infoModal, true)
								}
							} else if buttonIndex == 1 {
								// Go back to track list
								ui.app.SetRoot(trackList, true)
							} else {
								// Close and return to main view
								ui.app.SetRoot(ui.frame, true)
							}
						})

					ui.app.SetRoot(detailsModal, true)
				})

				// Set up key bindings for the track list
				trackList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					if event.Key() == tcell.KeyEscape {
						// Create a modal with album details
						detailsModal := tview.NewModal().
							SetText(text).
							AddButtons([]string{"Open Album in Spotify", "View Tracks", "Close"}).
							SetDoneFunc(func(buttonIndex int, buttonLabel string) {
								if buttonIndex == 0 && spotifyLink != "" {
									// Open album link
									err := openURL(spotifyLink)
									if err != nil {
										infoModal := tview.NewModal().
											SetText(fmt.Sprintf("Could not open browser automatically.\nSpotify link: %s\nURI: %s", spotifyLink, spotifyURI)).
											AddButtons([]string{"OK"}).
											SetDoneFunc(func(buttonIndex int, buttonLabel string) {
												ui.app.SetRoot(ui.frame, true)
											})
										ui.app.SetRoot(infoModal, true)
									} else {
										infoModal := tview.NewModal().
											SetText(fmt.Sprintf("Opening in browser:\n%s", spotifyLink)).
											AddButtons([]string{"OK"}).
											SetDoneFunc(func(buttonIndex int, buttonLabel string) {
												ui.app.SetRoot(ui.frame, true)
											})
										ui.app.SetRoot(infoModal, true)
									}
								} else if buttonIndex == 1 {
									// Go back to track list
									ui.app.SetRoot(trackList, true)
								} else {
									// Close and return to main view
									ui.app.SetRoot(ui.frame, true)
								}
							})

						ui.app.SetRoot(detailsModal, true)
						return nil
					}
					return event
				})

				// Set the title for the track list
				trackListFrame := tview.NewFrame(trackList).
					SetBorders(0, 0, 0, 0, 0, 0).
					AddText(fmt.Sprintf("Tracks from '%s'", album.Name), true, tview.AlignCenter, tcell.ColorWhite).
					AddText("↑/↓: Navigate • Enter: Open Track • ESC: Back to Album Details", false, tview.AlignCenter, tcell.ColorWhite)

				// Show the track list instead of the modal
				ui.app.SetRoot(trackListFrame, true)
				return
			}

			// If we couldn't get tracks or there was an error, just show the album details
			// Add buttons to the modal
			buttons := []string{"Open in Spotify", "Close"}

			modal.SetText(text).
				AddButtons(buttons).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonIndex == 0 && spotifyLink != "" {
						// Try to open the link in the default browser
						err := openURL(spotifyLink)
						if err != nil {
							// If opening the browser fails, just show the link
							infoModal := tview.NewModal().
								SetText(fmt.Sprintf("Could not open browser automatically.\nSpotify link: %s\nURI: %s", spotifyLink, spotifyURI)).
								AddButtons([]string{"OK"}).
								SetDoneFunc(func(buttonIndex int, buttonLabel string) {
									ui.app.SetRoot(ui.frame, true)
								})
							ui.app.SetRoot(infoModal, true)
						} else {
							// Show a confirmation that the link was opened
							infoModal := tview.NewModal().
								SetText(fmt.Sprintf("Opening in browser:\n%s", spotifyLink)).
								AddButtons([]string{"OK"}).
								SetDoneFunc(func(buttonIndex int, buttonLabel string) {
									ui.app.SetRoot(ui.frame, true)
								})
							ui.app.SetRoot(infoModal, true)
						}
					} else {
						ui.app.SetRoot(ui.frame, true)
					}
				})

			// Display the modal
			ui.app.SetRoot(modal, true)
			return
		}
	case "playlist":
		playlists := ui.results.([]spotify.SimplePlaylist)
		if row-1 < len(playlists) {
			playlist := playlists[row-1]

			spotifyLink = fmt.Sprintf("https://open.spotify.com/playlist/%s", playlist.ID)
			spotifyURI = string(playlist.URI)

			text = fmt.Sprintf("Playlist: %s\nOwner: %s\nTotal Tracks: %d\nSpotify Link: %s\nURI: %s",
				playlist.Name,
				playlist.Owner.DisplayName,
				playlist.Tracks.Total,
				spotifyLink,
				playlist.URI)

			if playlist.Description != "" {
				text += fmt.Sprintf("\nDescription: %s", playlist.Description)
			}

			// Add playlist tracks if showDetails is true
			if ui.showDetails {
				fullPlaylist, err := ui.client.GetPlaylist(ui.ctx, playlist.ID)
				if err == nil && fullPlaylist != nil && len(fullPlaylist.Tracks.Tracks) > 0 {
					text += "\n\nSample Tracks:"

					// Show up to 10 tracks
					trackLimit := 10
					if trackLimit > len(fullPlaylist.Tracks.Tracks) {
						trackLimit = len(fullPlaylist.Tracks.Tracks)
					}

					for i := 0; i < trackLimit; i++ {
						trackName := fmt.Sprintf("Track %d", i+1)

						// Try to get the track name
						if fullPlaylist.Tracks.Tracks[i].Track.Name != "" {
							trackName = fullPlaylist.Tracks.Tracks[i].Track.Name
						}

						text += fmt.Sprintf("\n%d. %s", i+1, trackName)
					}

					totalTracks := int(playlist.Tracks.Total)
					if totalTracks > trackLimit {
						text += fmt.Sprintf("\n...and %d more tracks", totalTracks-trackLimit)
					}
				}
			}
		}
	}

	// Add buttons to the modal
	buttons := []string{"Open in Spotify", "Close"}

	modal.SetText(text).
		AddButtons(buttons).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 && spotifyLink != "" {
				// Try to open the link in the default browser
				err := openURL(spotifyLink)
				if err != nil {
					// If opening the browser fails, just show the link
					infoModal := tview.NewModal().
						SetText(fmt.Sprintf("Could not open browser automatically.\nSpotify link: %s\nURI: %s", spotifyLink, spotifyURI)).
						AddButtons([]string{"OK"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							ui.app.SetRoot(ui.frame, true)
						})
					ui.app.SetRoot(infoModal, true)
				} else {
					// Show a confirmation that the link was opened
					infoModal := tview.NewModal().
						SetText(fmt.Sprintf("Opening in browser:\n%s", spotifyLink)).
						AddButtons([]string{"OK"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							ui.app.SetRoot(ui.frame, true)
						})
					ui.app.SetRoot(infoModal, true)
				}
			} else {
				ui.app.SetRoot(ui.frame, true)
			}
		})

	// Display the modal
	ui.app.SetRoot(modal, true)
}

// formatDuration formats milliseconds into a human-readable duration string (MM:SS)
func formatDuration(ms spotify.Numeric) string {
	totalSeconds := int(ms) / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
