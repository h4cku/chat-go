package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Client represents a connected user
type Client struct {
	conn     *websocket.Conn
	username string
	hub      *Hub
	send     chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		// Format message with username
		formattedMsg := fmt.Sprintf(`{"username":"%s","content":"%s","type":"message"}`, c.username, string(message))
		c.hub.broadcast <- []byte(formattedMsg)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
