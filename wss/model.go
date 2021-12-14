package wss

import "github.com/gorilla/websocket"

type Hub struct {
	Rooms      map[string]map[*Connection]bool
	Broadcast  chan Message
	Register   chan Subscription
	Unregister chan Subscription
}
type Subscription struct {
	Conn *Connection
	Room string
	User string
}

type Connection struct {
	Ws   *websocket.Conn
	Send chan []byte
}

type Message struct {
	Data []byte
	Room string
}

type Response struct {
	User string `json:"user"`
	Body string `json:"body"`
	Time string `json:"time"`
}
