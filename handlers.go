package main

import (
	"log"
	"net/http"
)

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous"
	}

	client := &Client{
		conn:     conn,
		username: username,
		hub:      hub,
		send:     make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Simple Go Chat</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        #messages { border: 1px solid #ccc; height: 400px; padding: 10px; overflow-y: scroll; margin-bottom: 10px; }
        #messageInput { width: 70%; padding: 10px; }
        #sendBtn { width: 25%; padding: 10px; }
        .message { margin: 5px 0; }
        .join { color: green; font-style: italic; }
        .leave { color: red; font-style: italic; }
        .username { font-weight: bold; color: #0066cc; }
    </style>
</head>
<body>
    <h1>Simple Go Chat</h1>
    <div id="messages"></div>
    <input type="text" id="messageInput" placeholder="Type your message..." />
    <button id="sendBtn">Send</button>

    <script>
        const username = prompt("Enter your username:") || "Anonymous";
        const ws = new WebSocket("ws://localhost:8888/ws?username=" + encodeURIComponent(username));
        const messages = document.getElementById("messages");
        const messageInput = document.getElementById("messageInput");
        const sendBtn = document.getElementById("sendBtn");

        ws.onmessage = function(event) {
            const data = JSON.parse(event.data);
            const messageDiv = document.createElement("div");
            messageDiv.className = "message";
            
            if (data.type === "join") {
                messageDiv.innerHTML = '<span class="join">🟢 ' + data.content + '</span>';
            } else if (data.type === "leave") {
                messageDiv.innerHTML = '<span class="leave">🔴 ' + data.content + '</span>';
            } else if (data.type === "system") {
                messageDiv.innerHTML = '<span class="join">ℹ️ ' + data.content + '</span>';
            } else {
                messageDiv.innerHTML = '<span class="username">' + data.username + ':</span> ' + data.content;
            }
            
            messages.appendChild(messageDiv);
            messages.scrollTop = messages.scrollHeight;
        };

        function sendMessage() {
            const message = messageInput.value.trim();
            if (message) {
                ws.send(message);
                messageInput.value = "";
            }
        }

        sendBtn.onclick = sendMessage;
        messageInput.addEventListener("keypress", function(e) {
            if (e.key === "Enter") {
                sendMessage();
            }
        });

        ws.onopen = function() {
            console.log("Connected to chat server");
        };

        ws.onclose = function() {
            console.log("Disconnected from chat server");
            messages.innerHTML += '<div class="message"><span class="leave">Disconnected from server</span></div>';
        };
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
