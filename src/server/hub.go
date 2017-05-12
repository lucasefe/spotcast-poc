package main

import (
	"encoding/json"
	"util"

	"github.com/Sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// logger
	log *logrus.Entry

	// a name to identify it.
	name string
}

func newHub(name string) *Hub {
	return &Hub{
		name:       name,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		log:        util.NewLogger().WithField("channel", name),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			h.removeClient(client)
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) addClient(client *Client) {
	h.clients[client] = true
	h.log.Infof("-> Clients: %d", len(h.clients))

	if len(h.clients) == 1 {
		client.send <- leaderAction()
	}
}

func (h *Hub) removeClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		h.log.Infof("<- Clients: %d", len(h.clients))
	}
}

// Action is..
type Action struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func leaderAction() []byte {
	data := map[string]string{"role": "leader"}
	action := &Action{Type: "SET_ROLE", Data: data}
	message, _ := json.Marshal(action)
	return message
}
