# GSpotify CLI

A simple command-line interface for searching Spotify's catalog for tracks, albums, and playlists.

## Features

- Search for tracks, albums, or playlists
- Display results in a tabular format
- Show detailed information about search results
- Limit the number of results displayed
- Color-coded output for better readability

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

### Options

- `-type`: Type of search (track, album, or playlist). Default is "track".
- `-query`: Search query. Required.
- `-artist`: Artist name to filter results (only for track search). Optional.
- `-limit`: Number of results to display. Default is 5.
- `-details`: Show detailed information about the results. Default is false.

### Examples

Search for tracks:
```
./gspotify -type=track -query="Bohemian Rhapsody"
```

Search for tracks by a specific artist:
```
./gspotify -type=track -query="Bohemian Rhapsody" -artist="Queen"
```

Search for albums with a limit of 3 results:
```
./gspotify -type=album -query="Dark Side of the Moon" -limit=3
```

Search for playlists with detailed information:
```
./gspotify -type=playlist -query="workout" -details
```

## Output

The application displays results in a tabular format with relevant information for each type of search:

- **Tracks**: ID, Track Name, Artist, Album, Popularity, URI
- **Albums**: ID, Album Name, Artist, Release Date, Total Tracks, URI
- **Playlists**: ID, Playlist Name, Owner, Total Tracks, URI

When the `-details` flag is used, additional information is displayed:

- **Tracks**: Audio features (Energy, Danceability, Valence, Tempo)
- **Albums**: First few tracks
- **Playlists**: Description and sample tracks

## Notes

- The application uses client credentials flow for authentication, so it can only access public data.
- The Spotify API has rate limits, so excessive usage may result in temporary blocks.
- You need to obtain your own Spotify API credentials from the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/). 