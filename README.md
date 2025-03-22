# GSpotify CLI

A simple command-line interface for searching and playing Spotify tracks, albums, and playlists.

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
- Built-in music player with playback controls
- Keep music playing option even after exiting the player interface

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/iamgaru/gspotify.git
   cd gspotify
   ```

2. Build the application:
   ```
   go build
   ```

3. Set up Spotify API credentials:
   ```
   export SPOTIFY_ID=your_client_id
   export SPOTIFY_SECRET=your_client_secret
   ```

## Usage

```
./gspotify [options]
```

### Command Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-t` | Type of search: track, album, or playlist | "track" |
| `-q` | Search query | Required |
| `-a` | Artist name to filter results (only for track search) | Optional |
| `-l` | Number of results to display | 5 |
| `-d` | Show detailed information about the results | false |
| `-i` | Run in interactive mode with a menu interface | false |
| `-r` | Return to interactive menu after viewing search results | false |
| `-k` | Keep music playing when exiting the player interface | false |
| `-p` | Automatically play the first result and exit | false |
| `-u` | Spotify user ID to look up profile information | Optional |

### Examples

#### Basic Search

Search for tracks:
```
./gspotify -q "Bohemian Rhapsody"
```

Search for tracks by a specific artist:
```
./gspotify -q "Bohemian Rhapsody" -a "Queen"
```

#### Changing Search Type

Search for albums:
```
./gspotify -t album -q "Dark Side of the Moon"
```

Search for playlists:
```
./gspotify -t playlist -q "workout"
```

#### Additional Options

Limit results to 3:
```
./gspotify -q "Dark Side of the Moon" -l 3
```

Show detailed information:
```
./gspotify -q "workout" -d
```

Run in interactive mode:
```
./gspotify -i
```

Search and return to menu:
```
./gspotify -q "Bohemian Rhapsody" -r
```

Play music and keep it playing when exiting the player:
```
./gspotify -q "Bohemian Rhapsody" -k
```

Automatically play the first result:
```
./gspotify -q "Bohemian Rhapsody" -p
```

Automatically play the first result and continue playing after exit:
```
./gspotify -q "Bohemian Rhapsody" -p -k
```

#### Combined Options

Search for Queen albums with detailed information:
```
./gspotify -t album -q "Queen" -d
```

Search for workout playlists, limit to 10, and show details:
```
./gspotify -t playlist -q "workout" -l 10 -d
```

#### User Profile Lookup

Look up a Spotify user's public profile:
```
./gspotify -u spotify
```

## Interactive Mode

When running in interactive mode, the application presents a user-friendly form where you can:

1. Select the search type (track, album, or playlist)
2. Enter your search query
3. Specify an artist name (for track searches)
4. Set the number of results to display
5. Choose whether to show detailed information

After submitting the form, the search results will be displayed in the same tabular format as the non-interactive mode.

### Playing Music

In interactive mode, you can play music by:

1. Searching for tracks (choose "track" as the search type)
2. Selecting a track from the search results
3. Clicking "Play" on the track details screen

You'll then be taken to the player interface where you can control playback.

## Music Player Controls

When playing a track, the player interface provides the following controls:

| Key | Function |
|-----|----------|
| Space | Play/Pause the current track |
| k | Toggle "Keep Playing" mode (ON/OFF) |
| Esc | Return to the previous menu |

### Keep Playing Mode

The "Keep Playing" feature allows you to continue listening to the current track even after exiting the player interface. This is useful when you want to continue browsing or searching while the music plays.

To use this feature:
1. While in the player interface, press `k` to toggle "Keep Playing" mode ON
2. The status will be displayed in the player interface
3. Press Esc to return to the menu while music continues playing
4. To stop playback later, return to the player interface and press Space to pause

Example workflow:
```
1. Search for a track in interactive mode
2. Select a track to play
3. In the player interface, press 'k' to enable Keep Playing
4. Press Esc to return to the menu while music continues
5. Continue browsing or searching while listening
```

## Output

The application displays results in a tabular format with relevant information for each type of search:

- **Tracks**: ID, Track Name, Artist, Album, Popularity, Spotify Link, URI
- **Albums**: ID, Album Name, Artist, Release Date, Total Tracks, Spotify Link, URI
- **Playlists**: ID, Playlist Name, Owner, Total Tracks, Spotify Link, URI

When the `-d` flag is used or "Show Detailed Results" is selected in interactive mode, additional information is displayed:

- **Tracks**: Audio features (Energy, Danceability, Valence, Tempo)
- **Albums**: First few tracks
- **Playlists**: Description and sample tracks

## Notes

- The application uses client credentials flow for authentication, so it can only access public data.
- The Spotify API has rate limits, so excessive usage may result in temporary blocks.
- You need to obtain your own Spotify API credentials from the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/).
- Music playback requires an active Spotify device (such as the Spotify desktop app or web player). 