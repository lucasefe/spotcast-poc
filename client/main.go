package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	srv := startHTTPServer(router)

away:
	for {
		select {
		case message := <-channel.Receive:
			log.Printf("Got %s\n", message)
			break

		case <-interrupt:
			srv.Shutdown(context.Background())
			closeWebsocket <- true
			break away
		}
	}

	<-time.After(time.Second)
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

func echo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	channel.Send("Echo!")
	fmt.Fprint(w, "Echo sent")
}
