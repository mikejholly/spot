package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/zmb3/spotify"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotify.NewAuthenticator(
		redirectURI,
		spotify.ScopeUserLibraryModify,
		spotify.ScopeUserModifyPlaybackState,
		spotify.ScopeUserReadPrivate,
		spotify.ScopeUserReadCurrentlyPlaying,
	)

	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Printf("\nPlease open this login URL in a browser:\n\n%s\n\n", url)

	// wait for auth to complete
	client := <-ch

	u, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %+v", err)
	}

	settingsFile := fmt.Sprintf("%s/.spotify.json", u.HomeDir)

	f, err := os.Create(settingsFile)
	if err != nil {
		log.Fatal("could not open settings file")
	}

	tok, err := client.Token()
	if err != nil {
		log.Fatal("failed to load token")
	}

	err = json.NewEncoder(f).Encode(tok)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Token saved to", settingsFile)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	ch <- &client

	fmt.Fprintf(w, "Login Completed! You may now close this window.")
}
