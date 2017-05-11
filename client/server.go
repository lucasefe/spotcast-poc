package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func startHTTPServer() *http.Server {
	router := httprouter.New()
	router.POST("/play/:songURI", play)
	router.POST("/pause", pause)
	router.POST("/stop", pause)

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
	songURI := ps.ByName("songURI")
	message, err := playSongAction(songURI)
	if err != nil {
		http.Error(w, "Error", 500) // Or Redirect?
		log.Printf("Error attempting to send play. Error: %s", err)
	}

	channel.Send(message)
	fmt.Fprint(w, fmt.Sprintf("Requesting play of song: %+v\n", songURI))
}

func pause(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	message, err := pauseAction()
	if err != nil {
		http.Error(w, "Error", 500) // Or Redirect?
		log.Printf("Error attempting to send pause. Error: %s", err)
	}

	channel.Send(message)
	fmt.Fprint(w, "Requesting pause")
}
