package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var (
	addr = flag.String("addr", "localhost:8081", "http service address")
	conn *websocket.Conn
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	closeWebsocket := make(chan bool)
	defer close(closeWebsocket)
	go connect(closeWebsocket)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	router := httprouter.New()
	router.POST("/", echo)

	srv := startHTTPServer(router)

	select {
	case <-interrupt:
		ctx := context.Background()
		srv.Shutdown(ctx)
		closeWebsocket <- true
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
	err := conn.WriteMessage(websocket.TextMessage, []byte("Echo!"))
	if err != nil {
		log.Println("write:", err)
		return
	}
	fmt.Fprint(w, "Echo sent")
}

func connect(stop chan bool) {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	conn = c
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer conn.Close()
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

loop:
	for {
		select {
		case <-stop:
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			conn.Close()
			break loop
		}
	}
}
