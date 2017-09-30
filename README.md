# spot

Spot is a simple command line tool that offers basic control over currently playing Spotify tracks. I made it because I find switching to the player distracting while working. It currently supports adding a song to your library (Spotify's equivalent to liking) and

## Setup

### 1. Install

Run `make install` in this repo's directory to install. Alternative, run `go install github.com/mikejholly/spot/cmd/spot`.

### 2. Authorize

Spotify only supports 2-legged OAuth and forces you to run a web server to complete the handshake. I've included a script (mostly cribbed from the Spotify Go library) which makes this easy.

1. Go to https://developer.spotify.com/my-applications
1. Create a new application, set the Redirect URI to `http://localhost:8080/callback`.
1. Add `SPOTIFY_ID` and `SPOTIFY_SECRET` to your enviroment. I've added these to my `.zshrc` file.
1. Run `make auth`. Open the URL and authorize the application.
1. You should have a valid settings file at `~/.spotify.json`.

## Usage

### Liking

Use `spot a` or `spot add` to add the currently playing song to your library.

### Nexting

Use `spot n` or `spot next` to skip to the next song in your playlist.
