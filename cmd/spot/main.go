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
			if err := addToLibrary(&client); err != nil {
				log.Fatalf("failed to add song to library: %v\n", err)
			}
			fmt.Println("saved song to library")
		case "r", "remove":
			if err := removeFromLibrary(&client); err != nil {
				log.Fatalf("failed to remove song from library: %v\n", err)
			}
			fmt.Println("removed song from library")
		case "n", "next":
			if err := client.Next(); err != nil {
				log.Fatalf("failed to move to next song: %v\n", err)
			}
			fmt.Println("moved to next song")
		case "p", "prev":
			if err := client.Previous(); err != nil {
				log.Fatalf("failed to move to previous song: %v\n", err)
			}
			fmt.Println("moved to previous song")
		case "t", "toggle":
			fmt.Println("toggling play/pause...")
			msg, err := togglePlay(&client)
			if err != nil {
				log.Fatalf("%v: %v\n", msg, err)
			}
			fmt.Println(msg)
		case "i", "info":
			if cp, err := client.PlayerCurrentlyPlaying(); err != nil || cp.Item == nil {
				log.Fatalf("failed to get song player is currently playing: %v\n", err)
			} else {
				fmt.Printf("current song: %s - %s\n", cp.Item.Name, cp.Item.Artists[0].Name)
			}
		}
	} else {
		if err := addToLibrary(&client); err != nil {
			log.Fatalf("failed to add song to library: %v\n", err)
		}
		fmt.Println("saved song to library")
	}

}

func addToLibrary(client *spotify.Client) error {
	cp, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return err
	}

	err = client.AddTracksToLibrary(cp.Item.ID)
	if err != nil {
		return err
	}

	return nil
}

func removeFromLibrary(client *spotify.Client) error {
	cp, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return err
	}

	err = client.RemoveTracksFromLibrary(cp.Item.ID)
	if err != nil {
		return err
	}

	return nil
}

func togglePlay(client *spotify.Client) (string, error) {
	state, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		// this is also thrown if no player is active
		return "failed to retrieve player state", err
	}
	// toggle play state
	if state.Playing {
		err := client.Pause()
		if err != nil {
			return "failed to pause", err
		}
		return "paused", nil
	} else {
		err := client.Play()
		if err != nil {
			return "failed to play", err
		}
		return "started playing", nil
	}
}
