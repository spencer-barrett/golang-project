package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// stores websock connections
var clients = make(map[*websocket.Conn]bool)

// converts html to WebSocket request
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handles incoming websocket requests
func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	clients[conn] = true
	log.Printf("New client connected. Total clients: %d\n", len(clients))

	for {
		_, message, _ := conn.ReadMessage()
		broadcast(message)
	}
}

// sends message to clients.
func broadcast(message []byte) {
	for client := range clients {
		client.WriteMessage(websocket.TextMessage, message)
	}
}

// serves html file
func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client.html")
}

// main function
func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", serveWebSocket)

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
