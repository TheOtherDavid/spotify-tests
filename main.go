package main

import (
	"context"
	"fmt"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState, spotifyauth.ScopeUserModifyPlaybackState),
	)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	var client *spotify.Client

	http.HandleFunc("/callback", completeAuth)

	go func() {
		url := auth.AuthURL(state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

		// wait for auth to complete
		client = <-ch

		// use the client to make calls that require authorization
		user, err := client.CurrentUser(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as:", user.ID)
		ctx := context.Background()
		playlistId := spotify.ID("06kZdWHb9ysIAn2GdFLUeC")
		//client.GetPlaylist(ctx, id)
		playlist := getPlaylist(ctx, client, playlistId)
		tracks := getTracksFromPlaylist(ctx, client, playlist)

	}()

	http.ListenAndServe(":8080", nil)

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

func getPlaylist(ctx context.Context, client *spotify.Client, playlistId spotify.ID) *spotify.FullPlaylist {
	fmt.Println("Beginning getPlaylist")

	playlist, err := client.GetPlaylist(ctx, playlistId)

	if err != nil {
		fmt.Println(err.Error)
	}
	fmt.Println(playlist.ID)
	return playlist
}

func getTracksFromPlaylist(ctx context.Context, client *spotify.Client, playlist *spotify.FullPlaylist) []spotify.FullTrack {

	return tracks
}
