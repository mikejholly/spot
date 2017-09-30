package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func main() {

	f, err := os.Open("/etc/spotify.json")
	if err != nil {
		log.Fatalf("failed to open config file: %+v", err)
	}

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		log.Fatalf("failed to load config file: %+v", err)
	}

	// NewAuthenticator loads OAuth settings from env
	auth := spotify.NewAuthenticator("")

	client := auth.NewClient(tok)

	cp, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		log.Fatalf("failed to load current song: %v", err)
	}

	err = client.AddTracksToLibrary(cp.Item.ID)
	if err != nil {
		log.Fatalf("failed to add track: %v", err)
	}

	fmt.Println(cp.Item.Name, "saved to library")
}
