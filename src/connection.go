package main


import (
	"github.com/gorilla/websocket"
	"log"
	"time"
	"net/http"
	"fmt"
	"sync"
)

const (
// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Connection struct {
	ws *websocket.Conn
	userId *string
	send chan []byte
	mu sync.RWMutex
}


func ConnectionInitialize(ws *websocket.Conn, roomId, userId *string) {
	// Создадим новое соединение
	c := NewConnection(ws, userId)
	// Получим или создадим новый хаб
	h := GetHub(*roomId)
	h.register <- c
	// Запустим поток на запись
	go c.writePump()
	// И будем продолжать слушать в этом потоке
	c.readPump(h)
}


func NewConnection(ws *websocket.Conn, userId *string) *Connection{
	return &Connection{send: make(chan []byte, 256), ws: ws, userId: userId}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Connection) readPump(h *hub) {
	defer func() {
		(*h).unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.ws.ReadMessage()

		if err != nil {
			//h.disconnect(c)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		print_binary(message)
		(*h).broadcast <- message
	}
}

func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				fmt.Println("a")
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
