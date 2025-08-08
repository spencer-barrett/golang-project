package main

import (
	"fmt"
	"log"
	"net/http"
	"sync" // We'll use a mutex to ensure our map is safe for concurrent access.

	"github.com/gorilla/websocket"
)

// Message is a struct to represent a message with content and an author.
type Message struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

// We'll use a global map to store all active WebSocket connections.
var clients = make(map[*websocket.Conn]bool)

// We'll use a mutex to ensure our map is safe for concurrent access.
var mutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin.
	},
}

// serves a WebSocket connection, handling incoming messages.
func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	// Add the new connection to our global map, with a lock to prevent race conditions.
	mutex.Lock()
	clients[conn] = true
	log.Printf("New client connected. Total clients: %d\n", len(clients))
	mutex.Unlock()

	// Remove the connection from the map when the function returns.
	defer func() {
		mutex.Lock()
		delete(clients, conn)
		log.Printf("Client disconnected. Total clients: %d\n", len(clients))
		mutex.Unlock()
	}()

	// Infinite loop to listen for messages from this client.
	for {
		// Read a message from the client as a JSON object.
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			// Handle client disconnection.
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v", err)
			}
			break
		}

		// Broadcast the message to all other clients.
		broadcast(msg)
	}
}

// broadcast sends a message to all connected clients.
func broadcast(message Message) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Printf("write error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// homePage now serves the static HTML file from the disk.
func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client.html")
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", serveWebSocket)

	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Server is running...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	log.Println(http.ListenAndServe(":8080", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			delete(clients, conn)
			return
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
