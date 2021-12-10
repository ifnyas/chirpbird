package wss

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	H = Hub{
		make(map[string]map[*Connection]bool),
		make(chan Message),
		make(chan Subscription),
		make(chan Subscription),
	}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	histories = []Message{}
	port      = ""
)

func Init() {
	if port == "" {
		setPort()
	}
	createWsRoute()
	H.Run()
}

func setPort() {
	freePort, err := GetFreePort()
	if err != nil {
		log.Println(err.Error())
	}
	port = ":" + strconv.Itoa(freePort)
}

func createWsRoute() {
	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		ServeWs(c.Writer, c.Request, c.Query("key"))
	})

	r.GET("/history", func(c *gin.Context) {
		roomHistories := loadMsg(c.Query("room"))
		array := []Response{}
		for _, history := range roomHistories {
			var res Response
			json.Unmarshal(history.Data, &res)
			array = append(array, res)
		}
		c.JSON(http.StatusOK, gin.H{"data": array})
	})

	go r.Run(port)
}

func WsProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme := c.Request.URL.Scheme
		if scheme == "" {
			scheme = "http"
		}

		host := c.Request.Host
		if strings.Contains(host, ":") {
			host = strings.Split(host, ":")[0]
		}

		director := func(req *http.Request) {
			req.URL.Scheme = scheme
			req.URL.Host = host + port
			req.Header["my-header"] = []string{c.Request.Header.Get("my-header")}
			delete(req.Header, "My-Header")
		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

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
		broadcast(s, "<i>("+"is offline"+")</i>", false)

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
	broadcast(s, "<i>("+"is online"+")</i>", false)

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

func broadcast(s Subscription, msg string, toSave bool) {
	firstMsg := modMsg([]byte(msg), s.User)
	m := Message{firstMsg, s.Room}
	H.Broadcast <- m

	if toSave {
		saveMsg(m)
	}
}

func (c *Connection) write(mt int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(mt, payload)
}

func saveMsg(m Message) {
	histories = append(histories, m)
}

func loadMsg(room string) []Message {
	selected := []Message{}
	for _, history := range histories {
		if history.Room == room {
			selected = append(selected, history)
		}
	}
	return selected
}

/*
	Utils
*/
func createSub(c *Connection, key string) Subscription {
	split := strings.SplitN(key, ":", 2)
	room := split[0]
	user := split[1]
	return Subscription{c, room, user}
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func modMsg(msg []byte, user string) []byte {
	res := Response{user, string(msg), time.Now().Format("15:04")}
	json, err := json.Marshal(res)
	if err != nil {
		json = []byte("")
	}
	return json
}
