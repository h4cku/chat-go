const username = prompt("Enter your username:") || "Anonymous";
const ws = new WebSocket(
  "ws://localhost:8888/ws?username=" + encodeURIComponent(username),
);
const messages = document.getElementById("messages");
const messageInput = document.getElementById("messageInput");
const sendBtn = document.getElementById("sendBtn");

ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  const messageDiv = document.createElement("div");
  messageDiv.className = "message";

  if (data.type === "join") {
    messageDiv.innerHTML = '<span class="join">🟢 ' + data.content + "</span>";
  } else if (data.type === "leave") {
    messageDiv.innerHTML = '<span class="leave">🔴 ' + data.content + "</span>";
  } else if (data.type === "system") {
    messageDiv.innerHTML = '<span class="join">ℹ️ ' + data.content + "</span>";
  } else {
    messageDiv.innerHTML =
      '<span class="username">' + data.username + ":</span> " + data.content;
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
  messages.innerHTML +=
    '<div class="message"><span class="leave">Disconnected from server</span></div>';
};
