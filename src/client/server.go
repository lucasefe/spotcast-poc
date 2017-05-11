package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func startHTTPServer() *http.Server {
	router := httprouter.New()
	router.POST("/play/:songURI", play)
	router.POST("/pause", pause)
	router.POST("/stop", pause)

	addr := ":8080"
	srv := &http.Server{Addr: addr}
	srv.Handler = router

	go func() {
		logger.Infof("HTTP Server listening on: %s", addr)
		if err := srv.ListenAndServe(); err != nil {
			logger.Warningf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	return srv
}

func play(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	songURI := ps.ByName("songURI")
	message, err := playSongAction(songURI)
	if err != nil {
		http.Error(w, "Error", 500) // Or Redirect?
		logger.Errorf("Error attempting to send play. Error: %s", err)
	}

	channel.Send(message)
	fmt.Fprint(w, fmt.Sprintf("Requesting play of song: %+v\n", songURI))
}

func pause(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	message, err := pauseAction()
	if err != nil {
		http.Error(w, "Error", 500) // Or Redirect?
		logger.Errorf("Error attempting to send pause. Error: %s", err)
	}

	channel.Send(message)
	fmt.Fprint(w, "Requesting pause")
}
