package main

import (
	"context"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmb3/spotify/v2"
)

// PlayerUI represents a UI for playing tracks and displaying track info
type PlayerUI struct {
	app           *tview.Application
	flex          *tview.Flex
	progressBar   *tview.TextView
	infoText      *tview.TextView
	imageView     *tview.TextView
	track         spotify.FullTrack
	client        *spotify.Client
	ctx           context.Context
	returnToMenu  func()
	timer         *time.Timer
	startTime     time.Time
	isPlaying     bool
	totalDuration time.Duration
}

// NewPlayerUI creates a new player UI
func NewPlayerUI(ctx context.Context, client *spotify.Client, track spotify.FullTrack) *PlayerUI {
	app := tview.NewApplication()
	progressBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	infoText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	imageView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	playerUI := &PlayerUI{
		app:           app,
		progressBar:   progressBar,
		infoText:      infoText,
		imageView:     imageView,
		track:         track,
		client:        client,
		ctx:           ctx,
		totalDuration: time.Duration(track.Duration) * time.Millisecond,
	}

	// Create layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(imageView, 0, 10, false).
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
			playerUI.stopPlayback()
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

		return event
	})

	// Prepare player UI
	playerUI.updateInfoText()
	playerUI.updateArtwork()

	return playerUI
}

// updateInfoText updates the track information display
func (p *PlayerUI) updateInfoText() {
	artists := make([]string, len(p.track.Artists))
	for i, artist := range p.track.Artists {
		artists[i] = artist.Name
	}

	info := fmt.Sprintf(
		"[green]Track:[white] %s\n[green]Artists:[white] %s\n[green]Album:[white] %s\n[green]Release Date:[white] %s\n\n"+
			"[yellow]Press Space to play/pause. Press Esc to return.[white]",
		p.track.Name,
		strings.Join(artists, ", "),
		p.track.Album.Name,
		p.track.Album.ReleaseDate,
	)

	p.infoText.SetText(info)
}

// updateArtwork attempts to display album artwork in ASCII form
func (p *PlayerUI) updateArtwork() {
	if len(p.track.Album.Images) > 0 {
		// Get the medium size image URL
		imageURL := p.track.Album.Images[0].URL
		if len(p.track.Album.Images) > 1 {
			imageURL = p.track.Album.Images[1].URL
		}

		go func() {
			asciiArt := fetchAndConvertToAscii(imageURL)
			p.app.QueueUpdateDraw(func() {
				p.imageView.SetText(asciiArt)
			})
		}()
	}
}

// fetchAndConvertToAscii fetches an image from a URL and converts it to ASCII art
func fetchAndConvertToAscii(url string) string {
	// Fetch the image
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("[red]Error fetching image: %v[white]", err)
	}
	defer resp.Body.Close()

	// For now, return a placeholder
	return `[yellow]
    _____________________
   /                     \
  |       ALBUM ART       |
  |                       |
  |      ___________      |
  |     /           \     |
  |    |     ♫      |     |
  |    |    ♫ ♫     |     |
  |    |   ♫   ♫    |     |
  |    |     ♫      |     |
  |     \___________/     |
  |                       |
  |                       |
   \_____________________/
   [white]`
}

// Play starts the playback UI
func (p *PlayerUI) Play() {
	p.startPlayback()
	p.app.SetRoot(p.flex, true).EnableMouse(true)
	if err := p.app.Run(); err != nil {
		fmt.Printf("Error running player UI: %v\n", err)
	}
}

// startPlayback starts track playback
func (p *PlayerUI) startPlayback() {
	// Start playback on your default device
	go func() {
		err := openSpotifyURI(string(p.track.URI))
		if err != nil {
			p.app.QueueUpdateDraw(func() {
				p.progressBar.SetText(fmt.Sprintf("[red]Error starting playback: %v[white]", err))
			})
			return
		}
	}()

	p.isPlaying = true
	p.startTime = time.Now()

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
}

// pausePlayback pauses the current playback
func (p *PlayerUI) pausePlayback() {
	// Execute the Spotify pause command based on OS
	go func() {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("osascript", "-e", `tell application "Spotify" to pause`)
		case "windows":
			cmd = exec.Command("powershell", "-c", "(New-Object -ComObject WScript.Shell).SendKeys(' ')")
		case "linux":
			cmd = exec.Command("dbus-send", "--print-reply", "--dest=org.mpris.MediaPlayer2.spotify",
				"/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player.Pause")
		default:
			p.app.QueueUpdateDraw(func() {
				p.progressBar.SetText("[red]Pause not supported on this OS[white]")
			})
			return
		}

		err := cmd.Run()
		if err != nil {
			p.app.QueueUpdateDraw(func() {
				p.progressBar.SetText(fmt.Sprintf("[red]Error pausing playback: %v[white]", err))
			})
			return
		}
	}()

	p.isPlaying = false
	if p.timer != nil {
		p.timer.Stop()
	}
}

// stopPlayback stops the current playback
func (p *PlayerUI) stopPlayback() {
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

// openSpotifyURI opens the specified Spotify URI
func openSpotifyURI(uri string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", uri)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", uri)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", uri)
	}

	return cmd.Run()
}
