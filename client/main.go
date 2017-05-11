package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

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

	srv := startHTTPServer()

	go playerPoller()

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

func playerPoller() {
	var lastResult *spoty.Result

	for {
		time.Sleep(time.Second)
		result, err := player.Status()
		if err != nil {
			log.Fatal(err)
		}

		if lastResult != nil || (lastResult.Track.CurrentSongURI() != result.Track.CurrentSongURI()) {
			lastResult = result
			track := result.Track
			log.Printf("Now Playing %+s, %+v\n", track.CurrentSongTitle(), track.CurrentSongURI())
		}
	}
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

func playSong(data map[string]string) {
	if song, ok := data["song"]; ok {
		player.Play(song)
		return
	}

	log.Printf("Wrong data: %+v", data)
}
