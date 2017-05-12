package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Channel is a connection/websocket abstraction
type Channel struct {
	conn    *websocket.Conn
	Receive chan string
}

// NewChannel creates a new Channel
func NewChannel(addr, channelName string) (*Channel, error) {
	channelPath := fmt.Sprintf("/channel/%s", channelName)
	u := url.URL{Scheme: "ws", Host: addr, Path: channelPath}
	logger.Debugf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}

	channel := &Channel{
		conn:    conn,
		Receive: make(chan string),
	}

	return channel, nil
}

// Connect sets up everything to receive and send messages to the server
func (c *Channel) Connect(stop chan bool) {
	done := make(chan struct{})

	// read messages
	go ReadMessages(c, done)

loop:
	for {
		select {
		case <-stop:
			err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Debug("write close:", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.conn.Close()
			break loop
		}
	}
}

// Send sends a text message through the active websocket connection
func (c *Channel) Send(text []byte) error {
	logger.Debugf("Sending action: %+v", string(text))
	err := c.conn.WriteMessage(websocket.TextMessage, text)
	if err != nil {
		return fmt.Errorf("Could not send message through channel: %v", err)
	}

	return nil
}

// Close closes channel allocated resources
func (c *Channel) Close() {
	c.conn.Close()
}

// ReadMessages reads messages from websocket
func ReadMessages(c *Channel, done chan struct{}) {
	defer c.Close()
	defer close(done)

	for {
		_, message, err := c.conn.ReadMessage()

		// TODO: Proper logging
		if err != nil {
			logger.Debug("read:", err)
			return
		}

		logger.Debugf("Receiving action: %+v", string(message))
		c.Receive <- string(message)
	}
}
