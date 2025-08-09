package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WireMessage struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

// stores websock connections
var clients = make(map[*websocket.Conn]bool)

// converts html to WebSocket request
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func broadcastCount() {
	payload := []byte(fmt.Sprintf(`{"type":"count","count":%d}`, len(clients)))
	broadcast(payload)
}

// handles incoming websocket requests
func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	name := r.URL.Query().Get("name")

	if name == "" {
		name = "anonymous_user"
	}

	clients[conn] = true
	log.Printf("%s connected. Total clients: %d\n", name, len(clients))

	msg := WireMessage{
		Author:  "system",
		Content: fmt.Sprintf("%s joined the chat\t%s", name, time.Now().Format(time.RFC1123)),
	}

	b, _ := json.Marshal(msg)

	broadcast(b)
	broadcastCount()

	defer func() {
		delete(clients, conn)
		err := conn.Close()
		if err != nil {
			return
		}

		msg := WireMessage{
			Author:  "system",
			Content: fmt.Sprintf("%s left the chat\t%s", name, time.Now().Format(time.RFC1123)),
		}

		b, _ := json.Marshal(msg)

		broadcast(b)
		broadcastCount()

		log.Printf("%s disconnected. Total clients: %d\n", name, len(clients))
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		broadcast(message)
	}
}

// sends message to clients.
func broadcast(message []byte) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

// main function
func main() {
	mux := http.NewServeMux()

	// WS endpoint
	mux.HandleFunc("/ws", serveWebSocket)

	// Static files from ./public
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs) // serves public/index.html at "/"

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", mux)
}
