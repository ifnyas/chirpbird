package wss

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var (
	H = Hub{
		make(map[string]map[*Connection]bool),
		make(chan Message),
		make(chan Subscription),
		make(chan Subscription),
	}
)

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.Register:
			connections := h.Rooms[s.Room]
			if connections == nil {
				connections = make(map[*Connection]bool)
				h.Rooms[s.Room] = connections
			}
			h.Rooms[s.Room][s.Conn] = true
		case s := <-h.Unregister:
			connections := h.Rooms[s.Room]
			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					delete(connections, s.Conn)
					close(s.Conn.Send)
					if len(connections) == 0 {
						delete(h.Rooms, s.Room)
					}
				}
			}
		case m := <-h.Broadcast:
			connections := h.Rooms[m.Room]
			for c := range connections {
				select {
				case c.Send <- m.Data:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Rooms, m.Room)
					}
				}
			}
		}
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request, key string) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &Connection{ws, make(chan []byte, 256)}
	s := createSub(c, key)

	H.Register <- s

	go s.writePump()
	go s.readPump()
}

func (s Subscription) readPump() {
	c := s.Conn

	defer func() {
		// broadcast user is offline
		broadcast(s, offlineMsg, false)

		H.Unregister <- s
		c.Ws.Close()
	}()

	c.Ws.SetReadLimit(maxMessageSize)
	c.Ws.SetReadDeadline(time.Now().Add(pongWait))
	c.Ws.SetPongHandler(func(string) error {
		c.Ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// broadcast user is online
	broadcast(s, onlineMsg, false)

	for {
		_, msg, err := c.Ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		broadcast(s, string(msg), true)
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

func createSub(c *Connection, key string) Subscription {
	split := strings.SplitN(key, ":", 2)
	room := split[0]
	user := split[1]
	return Subscription{c, room, user}
}

func broadcast(s Subscription, msg string, toSave bool) {
	firstMsg := modMsg([]byte(msg), s.User)
	m := Message{firstMsg, s.Room}
	H.Broadcast <- m

	if toSave {
		saveMsg(m)
	}
}
