package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmb3/spotify/v2"
)

// PlayerUI represents a UI for playing tracks and displaying track info
type PlayerUI struct {
	app            *tview.Application
	flex           *tview.Flex
	progressBar    *tview.TextView
	infoText       *tview.TextView
	track          spotify.FullTrack
	client         *spotify.Client
	ctx            context.Context
	returnToMenu   func()
	timer          *time.Timer
	startTime      time.Time
	isPlaying      bool
	totalDuration  time.Duration
	pausedPosition time.Duration // Add a field to store the paused position
	keepPlaying    bool          // Flag to determine if music should keep playing when exiting
	autoQuit       bool          // Flag to immediately exit after starting playback
}

// NewPlayerUI creates a new player UI
func NewPlayerUI(ctx context.Context, client *spotify.Client, track spotify.FullTrack, keepPlaying bool, autoQuit bool) *PlayerUI {
	app := tview.NewApplication()
	progressBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	infoText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	playerUI := &PlayerUI{
		app:           app,
		progressBar:   progressBar,
		infoText:      infoText,
		track:         track,
		client:        client,
		ctx:           ctx,
		totalDuration: time.Duration(track.Duration) * time.Millisecond,
		keepPlaying:   keepPlaying, // Use the provided value instead of hardcoded default
		autoQuit:      autoQuit,    // Set the auto-quit flag
	}

	// Create layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(infoText, 0, 3, false).
		AddItem(progressBar, 1, 0, false)

	flex.SetBorder(true).
		SetTitle(fmt.Sprintf(" Now Playing: %s ", track.Name)).
		SetTitleAlign(tview.AlignCenter)

	playerUI.flex = flex

	// Set up key bindings
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			// Only stop playback if keepPlaying is false
			if !playerUI.keepPlaying {
				playerUI.stopPlayback()
			}
			if playerUI.returnToMenu != nil {
				app.Stop()
				playerUI.returnToMenu()
				return nil
			}
			app.Stop()
		}

		// Handle space key for play/pause
		if event.Rune() == ' ' {
			if playerUI.isPlaying {
				playerUI.pausePlayback()
			} else {
				playerUI.startPlayback()
			}
		}

		// Handle 'k' key to toggle keep playing mode
		if event.Rune() == 'k' {
			playerUI.keepPlaying = !playerUI.keepPlaying
			playerUI.updateInfoText()
		}

		return event
	})

	// Prepare player UI
	playerUI.updateInfoText()

	return playerUI
}

// updateInfoText updates the track information display
func (p *PlayerUI) updateInfoText() {
	artists := make([]string, len(p.track.Artists))
	for i, artist := range p.track.Artists {
		artists[i] = artist.Name
	}

	keepPlayingStatus := "OFF"
	if p.keepPlaying {
		keepPlayingStatus = "ON"
	}

	info := fmt.Sprintf(
		"[green]Track:[white] %s\n[green]Artists:[white] %s\n[green]Album:[white] %s\n[green]Release Date:[white] %s\n\n"+
			"[yellow]Press Space to play/pause. Press 'k' to toggle keep playing (%s). Press Esc to return.[white]",
		p.track.Name,
		strings.Join(artists, ", "),
		p.track.Album.Name,
		p.track.Album.ReleaseDate,
		keepPlayingStatus,
	)

	p.infoText.SetText(info)
}

// Play starts the playback UI
func (p *PlayerUI) Play() {
	// Start playback and get the result channel
	resultCh := p.startPlayback()

	// If autoQuit is enabled, print a message and return without starting the UI
	if p.autoQuit {
		artists := make([]string, len(p.track.Artists))
		for i, artist := range p.track.Artists {
			artists[i] = artist.Name
		}

		fmt.Printf("Now playing: %s by %s from the album %s\n",
			p.track.Name,
			strings.Join(artists, ", "),
			p.track.Album.Name)
		fmt.Println("Waiting for playback to start...")

		// Wait for playback to start or error to occur
		// Add a timeout to prevent hanging indefinitely
		select {
		case err := <-resultCh:
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println("Playback started. The application will exit but music will continue playing.")
		case <-time.After(10 * time.Second):
			fmt.Println("Timed out waiting for playback to start. The Spotify client may still begin playback shortly.")
		}
		return
	}

	// Otherwise, display the player UI
	p.app.SetRoot(p.flex, true).EnableMouse(true)
	if err := p.app.Run(); err != nil {
		fmt.Printf("Error running player UI: %v\n", err)
	}
}

// startPlayback starts track playback
func (p *PlayerUI) startPlayback() chan error {
	// Create a channel to signal when playback has started or encountered an error
	resultCh := make(chan error, 1)

	// Start playback using Spotify Web API instead of opening URI
	go func() {
		// Get available devices first
		devices, err := p.client.PlayerDevices(p.ctx)
		if err != nil {
			if p.app != nil {
				p.app.QueueUpdateDraw(func() {
					p.progressBar.SetText(fmt.Sprintf("[red]Error getting devices: %v[white]", err))
				})
			}
			resultCh <- fmt.Errorf("error getting devices: %v", err)
			return
		}

		// Check if there are any active devices
		if len(devices) == 0 {
			if p.app != nil {
				p.app.QueueUpdateDraw(func() {
					p.progressBar.SetText("[red]No active Spotify devices found. Please open Spotify on any device first.[white]")
				})
			}
			resultCh <- fmt.Errorf("no active Spotify devices found")
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

		// Set playback options
		playOpts := &spotify.PlayOptions{
			URIs: []spotify.URI{p.track.URI},
		}

		// If we have a paused position, set the position_ms parameter to resume from that point
		if p.pausedPosition > 0 {
			positionMs := spotify.Numeric(p.pausedPosition.Milliseconds())
			playOpts.PositionMs = positionMs
		}

		// If we have a device ID, specify it
		if deviceID != "" {
			playOpts.DeviceID = &deviceID
		}

		// Start playback on the device
		err = p.client.PlayOpt(p.ctx, playOpts)
		if err != nil {
			if p.app != nil {
				p.app.QueueUpdateDraw(func() {
					p.progressBar.SetText(fmt.Sprintf("[red]Error starting playback: %v[white]", err))
				})
			}
			resultCh <- fmt.Errorf("error starting playback: %v", err)
			return
		}

		// Signal that playback has started successfully
		resultCh <- nil
	}()

	// If we're resuming from a paused state, adjust the startTime to account for the previous playback
	if p.pausedPosition > 0 {
		p.startTime = time.Now().Add(-p.pausedPosition)
	} else {
		p.startTime = time.Now()
	}

	p.isPlaying = true

	// Start a timer to update the progress bar every second
	if p.timer != nil {
		p.timer.Stop()
	}

	p.timer = time.NewTimer(time.Second)
	go func() {
		for range p.timer.C {
			if !p.isPlaying {
				break
			}

			elapsed := time.Since(p.startTime)
			if elapsed > p.totalDuration {
				p.stopPlayback()
				break
			}

			p.updateProgressBar(elapsed)
			p.timer.Reset(time.Second)
		}
	}()

	return resultCh
}

// pausePlayback pauses the current playback
func (p *PlayerUI) pausePlayback() {
	// Use Spotify Web API to pause playback instead of OS-specific commands
	go func() {
		err := p.client.Pause(p.ctx)
		if err != nil {
			p.app.QueueUpdateDraw(func() {
				p.progressBar.SetText(fmt.Sprintf("[red]Error pausing playback: %v[white]", err))
			})
			return
		}
	}()

	// Store the current position when pausing
	p.pausedPosition = time.Since(p.startTime)

	p.isPlaying = false
	if p.timer != nil {
		p.timer.Stop()
	}
}

// stopPlayback stops the current playback
func (p *PlayerUI) stopPlayback() {
	// Actually stop the playback using Spotify API
	go func() {
		err := p.client.Pause(p.ctx)
		if err != nil {
			// Just log the error, don't need to display as we're exiting anyway
			fmt.Printf("Error stopping playback: %v\n", err)
		}
	}()

	p.isPlaying = false
	if p.timer != nil {
		p.timer.Stop()
	}
}

// updateProgressBar updates the progress bar based on the elapsed time
func (p *PlayerUI) updateProgressBar(elapsed time.Duration) {
	if elapsed > p.totalDuration {
		elapsed = p.totalDuration
	}

	// Calculate percentage
	percentage := int(float64(elapsed) / float64(p.totalDuration) * 100)

	// Create a progress bar with fill characters
	barWidth := 50
	filled := barWidth * percentage / 100

	bar := "[green]"
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "[white]"

	// Format time as mm:ss
	elapsedSeconds := int(elapsed.Seconds())
	totalSeconds := int(p.totalDuration.Seconds())

	timeText := fmt.Sprintf(" %02d:%02d / %02d:%02d (%d%%)",
		elapsedSeconds/60, elapsedSeconds%60,
		totalSeconds/60, totalSeconds%60,
		percentage)

	p.app.QueueUpdateDraw(func() {
		p.progressBar.SetText(bar + timeText)
	})
}

// SetReturnToMenuFunction sets the function to return to the main menu
func (p *PlayerUI) SetReturnToMenuFunction(returnFunc func()) {
	p.returnToMenu = returnFunc
}
