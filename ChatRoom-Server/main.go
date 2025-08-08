package ChatRoom_Server

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

type Message struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message
}
