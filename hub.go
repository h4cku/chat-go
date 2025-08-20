package main

import (
	"fmt"
	"log"
	"sync"
)

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			clientCount := len(h.clients)
			h.mutex.Unlock()

			log.Printf("User %s connected. Total users: %d", client.username, clientCount)

			// Send welcome message to the new client
			welcomeMsg := fmt.Sprintf(`{"username":"System","content":"Welcome to the chat, %s! There are %d users online.","type":"system"}`, client.username, clientCount)
			select {
			case client.send <- []byte(welcomeMsg):
			default:
				close(client.send)
				delete(h.clients, client)
			}

			// Notify others that user joined (but not the user themselves)
			joinMsg := fmt.Sprintf(`{"username":"System","content":"%s joined the chat","type":"join"}`, client.username)
			h.mutex.RLock()
			for c := range h.clients {
				if c != client { // Don't send to the newly joined client
					select {
					case c.send <- []byte(joinMsg):
					default:
						close(c.send)
						delete(h.clients, c)
					}
				}
			}
			h.mutex.RUnlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				clientCount := len(h.clients)

				log.Printf("User %s disconnected. Total users: %d", client.username, clientCount)

				// Notify others that user left
				leaveMsg := fmt.Sprintf(`{"username":"System","content":"%s left the chat","type":"leave"}`, client.username)
				for c := range h.clients {
					select {
					case c.send <- []byte(leaveMsg):
					default:
						close(c.send)
						delete(h.clients, c)
					}
				}
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}
