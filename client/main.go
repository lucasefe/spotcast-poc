package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gitlab.com/lucasefe/spotcast/spoty"

	"github.com/julienschmidt/httprouter"
)

var (
	addr = flag.String("addr", "localhost:8081", "http service address")

	channel *Channel
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	closeWebsocket := make(chan bool)
	defer close(closeWebsocket)

	c, err := NewChannel(*addr)
	if err != nil {
		log.Fatalf("Could not create channel: %v", err)
	}
	channel = c
	go channel.Connect(closeWebsocket)
	defer channel.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	router := httprouter.New()
	router.POST("/", echo)
	router.POST("/play", play)

	srv := startHTTPServer(router)

away:
	for {
		select {
		case message := <-channel.Receive:

			action, err := parseAction(message)
			if err != nil {
				log.Printf("Could not parse message %+v. Got error: %s", message, err)
			}

			switch action.Type {
			case "PLAY_SONG":
				playSong(action.Data)
				break
			default:
				log.Printf("Don't know what this means: %+v", message)
			}

			break

		case <-interrupt:
			srv.Shutdown(context.Background())
			closeWebsocket <- true
			break away
		}
	}

	<-time.After(time.Second)
}

func parseAction(message string) (*Action, error) {
	action := &Action{}

	if err := json.Unmarshal([]byte(message), &action); err != nil {
		return nil, err
	}

	return action, nil
}

// Action ..
type Action struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func startHTTPServer(router *httprouter.Router) *http.Server {
	srv := &http.Server{Addr: ":8080"}
	srv.Handler = router

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	return srv
}

func play(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	message, err := playSongAction()
	if err != nil {
		http.Error(w, "Error", 500) // Or Redirect?
		log.Printf("Error playing song error: %s", err)
	}

	channel.Send(message)
	fmt.Fprint(w, "Play sent")
}

func newPlayAction(songURI string) *Action {
	return &Action{
		Type: "PLAY_SONG",
		Data: map[string]string{"song": songURI},
	}
}

func playSong(data map[string]string) {
	spoty.Connect()

	if song, ok := data["song"]; ok {
		spoty.Play(song)
		return
	}

	log.Printf("Wrong data: %+v", data)
}

func playSongAction() ([]byte, error) {
	action := newPlayAction("spotify:artist:08td7MxkoHQkXnWAYD8d6Q")

	message, err := json.Marshal(action)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func echo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	channel.Send([]byte("Echo!"))
	fmt.Fprint(w, "Echo sent")
}
