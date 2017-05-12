package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"spoty"
	"util"

	"github.com/Sirupsen/logrus"
)

// Role is either follower or leader
type Role string

const (
	// LeaderRole get to decide what to play, by selecting the song on the player.
	LeaderRole = "leader"
	// FollowerRole do not listen to the current player for changes
	FollowerRole = "follower"
)

var (
	serverAddress = flag.String("remote", "localhost:8081", "remote server host:port")
	devmode       = flag.Bool("dev", false, "enable dev mode")
	httpEnabled   = flag.Bool("httpEnabled", false, "enable http server")

	channel *Channel
	player  spoty.Session
	role    = Role(FollowerRole)

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

func nowPlaying(result *spoty.Result) {
	logger.Debugf("Now Playing %+s, %+v\n", result.Track.CurrentSongTitle(), result.Track.CurrentSongURI())
}

func playerPoller() {
	var lastResult *spoty.Result

	var process = func(result *spoty.Result) {
		defer func() { lastResult = result }()

		// First run, or same song
		sameSong := lastResult != nil && result.Track.CurrentSongURI() == lastResult.Track.CurrentSongURI()
		if !sameSong {
			nowPlaying(result)
		}

		if role != LeaderRole || lastResult == nil {
			return
		}

		playToPause := lastResult.Playing && !result.Playing
		pauseToPlay := !lastResult.Playing && result.Playing
		// logger.Debugf("Player status: playToPause=%v pauseToPlay=%v\n", playToPause, pauseToPlay)

		if !sameSong {
			logger.Debugf("Sending Play")
			channel.Send(playSongAction(result.Track.CurrentSongURI()))
			return
		}

		if playToPause {
			logger.Debugf("Sending Pause")
			channel.Send(pauseAction())
			return
		}

		if pauseToPlay {
			logger.Debugf("Sending Resume")
			channel.Send(resumeAction())
		}
	}

	for {
		time.Sleep(time.Second)
		result, err := player.Status()
		if err != nil {
			logger.Fatal(err)
		}

		process(result)
	}
}

func getPlayer(faked bool) spoty.Session {
	if faked {
		return spoty.NewFakedSession()
	}

	player, err := spoty.NewSession()
	if err != nil {
		panic(err)
	}

	return player
}
