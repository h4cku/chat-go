package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Message structure for JSON communication
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Type     string `json:"type"` // "message", "join", "leave"
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in this example
	},
}

func main() {
	hub := newHub()
	go hub.run()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/ui/", http.StripPrefix("/ui/", fs))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(hub, w, r)
	})

	fmt.Println("Chat server starting on :8888")
	fmt.Println("Open http://localhost:8888 in multiple browser tabs to test")

	log.Fatal(http.ListenAndServe(":8888", nil))
}
