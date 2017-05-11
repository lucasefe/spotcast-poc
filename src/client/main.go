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

// Role is either follower or leader
type Role string

const (
	// LeaderRole get to decide what to play, by selecting the song on the player.
	LeaderRole = "leader"
	// FollowerRole don't listen to the current player for changes
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

			if role == LeaderRole {
				message, err := playSongAction(track.CurrentSongURI())
				if err != nil {
					logger.Errorf("Error attempting to send play. Error: %s", err)
					break
				}

				channel.Send(message)
			}
		}
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

func playSong(data map[string]string) {
	song, ok := data["song"]

	if !ok {
		logger.Warningf("Wrong data: %+v", data)
		return
	}

	if role == LeaderRole {
		logger.Warningf("Not gonna play song %s, I'm a leader!", song)
		return
	}

	player.Play(song)

}

func setRole(data map[string]string) {
	if newRole, ok := data["role"]; ok {
		logger.Infof("Setting new role: %s", newRole)
		role = Role(newRole)
	}
}
