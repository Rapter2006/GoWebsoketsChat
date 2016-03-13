package main


type hub struct {
	id string
	connections map[*Connection]bool
	broadcast chan []byte
	register chan *Connection
	unregister chan *Connection
}



func NewHub(id string) *hub {
	h := hub {
		id: id,
		broadcast:   make(chan []byte),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
	}
	return &h
}

func (h *hub) disconnect(connection *Connection) {

	// Отпишим соединение от комнаты
	h.unregister <- connection

	// Закроем соединение
	connection.mu.Lock()
	close(connection.send)
	connection.ws.Close()
	connection.mu.Unlock()

	// TODO продумать удаление хаба, когда тим не осталось соединений
}


func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}