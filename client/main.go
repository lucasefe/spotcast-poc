package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Sirupsen/logrus"

	"gitlab.com/lucasefe/spotcast/spoty"
	"gitlab.com/lucasefe/spotcast/util"
)

var (
	serverAddress = flag.String("remote", "localhost:8081", "remote server host:port")
	devmode       = flag.Bool("dev", false, "enable dev mode")
	httpEnabled   = flag.Bool("httpEnabled", false, "enable http server")

	channel *Channel
	player  spoty.Session

	logger *logrus.Logger
)

func main() {
	flag.Parse()
	logger = util.NewLogger()
	player = getPlayer(*devmode)

	closeWebsocket := make(chan bool)
	defer close(closeWebsocket)

	c, err := NewChannel(*serverAddress)
	if err != nil {
		logger.Fatalf("Could not create channel: %v", err)
	}

	channel = c
	go channel.Connect(closeWebsocket)
	defer channel.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var srv *http.Server
	if *httpEnabled {
		srv = startHTTPServer()
	}

	go playerPoller()

mainLoop:
	for {
		select {
		case message := <-channel.Receive:
			handleAction(message)
			break

		case <-interrupt:
			if srv != nil {
				srv.Shutdown(context.Background())
			}
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
			logger.Fatal(err)
		}

		if lastResult == nil || (lastResult.Track.CurrentSongURI() != result.Track.CurrentSongURI()) {
			lastResult = result
			track := result.Track
			logger.Debugf("Now Playing %+s, %+v\n", track.CurrentSongTitle(), track.CurrentSongURI())
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

	logger.Warningf("Wrong data: %+v", data)
}
