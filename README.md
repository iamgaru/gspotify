# gspotty

<p align="center">
  <img src="assets/images/gs-gopher.png" alt="GSpotify Gopher" width="300">
</p>

<p align="center">
  <img src="assets/images/dubstep-playlist-search.png" alt="Playlist Search" width="800">
  <br>
  <em>Searching for dubstep playlists</em>
</p>

<p align="center">
  <img src="assets/images/multiple-players.png" alt="Multiple Players" width="800">
  <br>
  <em>Support for multiple players</em>
</p>

<p align="center">
  <img src="assets/images/quick-play.png" alt="Quick Play" width="800">
  <br>
  <em>Quick play feature in action</em>
</p>

## Table of Contents

- [Project Structure](#project-structure)
- [Features](#features)
- [Building and Testing](#building-and-testing)
  - [Prerequisites](#prerequisites)
  - [Build Commands](#build-commands)
  - [Testing](#testing)
- [Installation](#installation)
- [Usage](#usage)
  - [Authentication](#authentication)
  - [Command Flags](#command-flags)
  - [Examples](#examples)
    - [Basic Search](#basic-search)
    - [Changing Search Type](#changing-search-type)
    - [Additional Options](#additional-options)
    - [Combined Options](#combined-options)
    - [User Profile Lookup](#user-profile-lookup)
- [Interactive Mode](#interactive-mode)
  - [Playing Music](#playing-music)
- [Music Player Controls](#music-player-controls)
  - [Keyboard Controls](#keyboard-controls)
  - [Playback Modes](#playback-modes)
  - [Device Management](#device-management)
- [Output Format](#output-format)
  - [Search Results Display](#search-results-display)
  - [Player Interface](#player-interface)
  - [Error Handling](#error-handling)
- [Notes](#notes)
- [Quick Play Script](#quick-play-script)
  - [Installation](#installation-1)
  - [Usage](#usage-1)
- [License](#license)
- [Author & Version](#author--version)

A simple command-line interface for searching and playing Spotify tracks, albums, and playlists.

## Project Structure

```
gspotty/
├── assets/
│   └── images/          # Application images and screenshots
├── cmd/
│   └── gspotty/          # Main application entry point
├── internal/
│   ├── cli/             # CLI implementation and Spotify client integration
│   ├── config/          # Configuration management
│   ├── menu/            # Interactive menu implementation
│   ├── player/          # Music player implementation
│   ├── profile/         # User profile functionality
│   ├── testutils/       # Test utilities and mocks
│   ├── ui/              # UI components
│   └── utils/           # Utility functions
├── scripts/             # Convenience scripts
│   └── play            # Script for quick music playback
├── configs/             # Configuration files
├── Makefile            # Build and test automation
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── LICENSE             # Project license
└── README.md           # Project documentation
```

## Features

- Search for tracks, albums, or playlists
- Display results in a tabular format
- Show detailed information about search results
- Limit the number of results displayed
- Color-coded output for better readability
- Interactive menu mode for easier searching
- Return to menu option after viewing search results
- Simple single-letter flags for easy command usage
- User profile lookup functionality
- Built-in music player with playback controls:
  - Play/Pause
  - Next/Previous track
  - Seek position
  - Volume control
- Keep music playing option even after exiting the player interface
- Support for playlist, search, and album playback modes with next track functionality
- Automatic looping in playlist, search, and album modes when "Keep Playing" is enabled
- Convenience scripts for common operations
- Secure token management with automatic refresh
- Cross-platform support (Linux, Windows, macOS)

## Building and Testing

### Prerequisites

1. Go 1.x or higher
2. Make

### Build Commands

Build the application using Make:
```bash
# Build the binary
make build

# Clean build artifacts
make clean

# Run all tests
make test

# Show available make commands
make help
```

Or build manually:
```bash
go build -o gspotty ./cmd/gspotty
```

### Testing

The project includes comprehensive tests for all components. Run the tests using:

```bash
make test
```

Tests are organized into:
- Unit tests for individual packages
- Integration tests for end-to-end functionality
- Mock implementations for external dependencies (e.g., Spotify API)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/iamgaru/gspotty.git
   cd gspotty
   ```

2. Build the application:
   ```
   make build
   ```