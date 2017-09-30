package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func main() {

	u, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %s", err)
	}

	confFile := fmt.Sprintf("%s/.spotify.json", u.HomeDir)

	f, err := os.Open(confFile)
	if err != nil {
		log.Fatalf("failed to open config file: %s", err)
	}

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		log.Fatalf("failed to load config file: %s", err)
	}

	// NewAuthenticator loads OAuth settings from env
	auth := spotify.NewAuthenticator("")

	client := auth.NewClient(tok)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "a", "add":
			addToLibrary(&client)
		case "n", "next":
			nextSong(&client)
		}
	} else {
		addToLibrary(&client)
	}

}

func addToLibrary(client *spotify.Client) {
	cp, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		log.Fatalf("failed to load current song: %s", err)
	}

	err = client.AddTracksToLibrary(cp.Item.ID)
	if err != nil {
		log.Fatalf("failed to add track: %s", err)
	}

	fmt.Println(cp.Item.Name, "saved to library")
}

func nextSong(client *spotify.Client) {
	err := client.Next()
	if err != nil {
		log.Fatal("failed to move to next song")
	}
}
