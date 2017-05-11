package main

import (
	"flag"
	"net/http"
	"util"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	addr   = flag.String("addr", ":8081", "http service address")
	logger *logrus.Logger
)

func main() {
	flag.Parse()
	logger = util.NewLogger()

	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	logger.Infof("Listening on port %v", *addr)
	logger.Fatal(http.ListenAndServe(*addr, nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	logger.Debug(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}
