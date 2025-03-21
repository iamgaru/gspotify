# GSpotify CLI

<p align="center">
  <img src="gs-gopher.png" alt="GSpotify Gopher" width="300">
</p>

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

4. **Authorization**: On first use, you will be prompted to authorize the application to access your Spotify account. This is a one-time process where a browser window will open for you to sign in to Spotify and grant the necessary permissions.

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
- **First-time Authorization**: When you first use the application, you will be prompted to authorize it to play Spotify tracks. A browser window will open automatically, and you'll need to log in to your Spotify account and approve the requested permissions. This authorization only happens once, and the app will save your credentials for future use.
- Music playback requires an active Spotify device (such as the Spotify desktop app or web player). 

## Quick Play Script

For a more convenient way to play music, you can use the included "play" bash script:

### Installation

1. Make the script executable:
   ```
   chmod +x play
   ```

2. Move it to a directory in your PATH (optional):
   ```
   sudo mv play /usr/local/bin/
   ```

### Usage

Simply run the script followed by your search query:
```
play Bohemian Rhapsody
```

This is equivalent to:
```
./gspotify -t track -q "Bohemian Rhapsody" -p -k
```

The script will:
1. Search for the specified track
2. Automatically play the first match (-p flag)
3. Continue playing even after exiting (-k flag)
4. Display a confirmation message

The script is a convenient shorthand when you just want to quickly play a song without interacting with the search results. 

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author & Version

**Author:** Nick Conolly  
**Version:** 0.0.2  
**GitHub:** [https://github.com/iamgaru](https://github.com/iamgaru) 