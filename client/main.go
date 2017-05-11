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

	"github.com/julienschmidt/httprouter"
	"gitlab.com/lucasefe/spotcast/spoty"
)

var (
	serverAddress = flag.String("remote", "localhost:8081", "remote server host:port")
	verbose       = flag.Bool("verbose", false, "enable verbose mode")
	devmode       = flag.Bool("dev", false, "enable dev mode")

	channel *Channel
	player  spoty.Session
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	player = getPlayer(*devmode)

	if *verbose || *devmode {
		player.SetVerbose()
	}

	closeWebsocket := make(chan bool)
	defer close(closeWebsocket)

	c, err := NewChannel(*serverAddress)
	if err != nil {
		log.Fatalf("Could not create channel: %v", err)
	}

	channel = c
	go channel.Connect(closeWebsocket)
	defer channel.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	router := httprouter.New()
	router.POST("/play", play)

	srv := startHTTPServer(router)

mainLoop:
	for {
		select {
		case message := <-channel.Receive:
			handleAction(message)
			break

		case <-interrupt:
			srv.Shutdown(context.Background())
			closeWebsocket <- true
			break mainLoop
		}
	}

	<-time.After(time.Second)
}

func handleAction(message string) {
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
	if song, ok := data["song"]; ok {
		player.Play(song)
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

func getPlayer(faked bool) spoty.Session {
	if faked {
		return &spoty.FakedSession{}
	}

	player, err := spoty.NewSession()
	if err != nil {
		panic(err)
	}

	return player
}
