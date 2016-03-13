package main

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"log"
)

type Connection struct {
	ws *websocket.Conn
	send chan []byte
	dead bool
	mu sync.RWMutex
}

func NewConnection(ws *websocket.Conn, send chan []byte) *Connection {
	return &Connection{
		ws:   ws,
		send: send,
	}
}

func (c *Connection) Send(message []byte, fin chan struct{}, r *Room) {
	defer func() {
		//Tell the calling function that this goroutine is done sending
		fin <- struct{}{}
	}()

	c.mu.RLock()
	if c.dead {
		//Channel is already dead, we cannot send on it anymore and we must exit
		c.mu.RUnlock()
		return
	}
	c.mu.RUnlock()

	//We don't want to try to send over the channel if another
	//goroutine has closed this channel in the meantime. Thus, we
	//must block writing before we send over this channel.
	c.mu.Lock()
	select {
	//If the message is sent over the websocket, unlock this connection and continue
	case c.send <- message:
		c.mu.Unlock()

	//If we cannot send, this means that the user's buffer is full. At this point we basically
	//assume that the user disconnected or is just stuck.
	default:
	//Unlock before unregistering since the act of unregistering triggers changes in c
		c.mu.Unlock()
		r.unregister <- c
	}
}


func (c *Connection) Reader(r *Room) {
	c.ws.SetReadLimit(cfg.maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(cfg.readWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(cfg.pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		r.broadcast <- message
	}
}

func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(cfg.writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *Connection) Writer() {
	ticker := time.NewTicker(cfg.pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
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