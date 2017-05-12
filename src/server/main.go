package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"util"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	address = flag.String("addr", ":8081", "http service address")
	logger  *logrus.Logger
)

func main() {
	flag.Parse()
	logger = util.NewLogger()

	hub := newHub()
	go hub.run()
	http.HandleFunc("/", webHandler())
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	addr := *address
	if port, ok := os.LookupEnv("PORT"); ok {
		addr = fmt.Sprintf(":%s", port)
	}

	logger.Infof("Listening on port %v", addr)
	logger.Fatal(http.ListenAndServe(addr, nil))
}

func webHandler() func(http.ResponseWriter, *http.Request) {
	var webRoot string

	if root, ok := os.LookupEnv("SPOTCAST_WEBROOT"); ok {
		webRoot = root
	} else {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		webRoot = path.Join(dir, "../public")
	}

	homePage := fmt.Sprintf("%s/%s", webRoot, "index.html")

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.URL)

		if r.URL.Path != "/" {
			http.Error(w, "Not found", 404)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		if _, err := os.Stat(homePage); os.IsNotExist(err) {
			logger.Warning(err)
		}

		http.ServeFile(w, r, homePage)
	}
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
