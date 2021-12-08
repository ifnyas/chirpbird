package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(w http.ResponseWriter, r *http.Request, room string, user string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := &Connection{ws, make(chan []byte, 256)}
	s := Subscription{c, room, user}

	H.Register <- s

	go s.writePump()
	go s.readPump()
}

func (s Subscription) readPump() {
	c := s.Conn
	defer func() {
		H.Unregister <- s
		c.Ws.Close()
	}()

	c.Ws.SetReadLimit(maxMessageSize)
	c.Ws.SetReadDeadline(time.Now().Add(pongWait))
	c.Ws.SetPongHandler(func(string) error {
		c.Ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.Ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg = modMsg(msg, s.User)
		m := Message{msg, s.Room}

		H.Broadcast <- m
	}
}

func (s *Subscription) writePump() {
	c := s.Conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			err := c.write(websocket.TextMessage, message)
			if err != nil {
				return
			}
		case <-ticker.C:
			err := c.write(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}
}

func (c *Connection) write(mt int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(mt, payload)
}

func modMsg(msg []byte, user string) []byte {
	res := Response{user, string(msg), time.Now().Format("15:04")}
	json, err := json.Marshal(res)
	if err != nil {
		json = []byte("")
	}
	return json
}
