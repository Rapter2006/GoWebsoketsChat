package main

type Room struct {
	// идентификатор комнаты
	id string
	// зарегистрированные соединения
	connections connectionMap
	// канал для рассылки броадкаст сообщений
	broadcast chan []byte
	// канал для отлова сообщений о регистрации
	register chan *Connection
	//  канал для отлова сообщений о разрыве соединений
	unregister chan *Connection
}

func (r *Room) Register(c *Connection) {
	r.register <- c
}

func (r *Room) Unregister(c *connection) {
	r.unregister <- c
}

func (r *Room) Broadcast(s []byte) {
	r.broadcast <- s
}

func (r *Room) Run() {
	for {
		select {
		case connection := <-r.register:
		//Add a connection
			go r.connect(connection)
		case connection := <-r.unregister:
		//Delete a connection
			go r.disconnect(connection)
		case message := <-r.broadcast:
		//We've received a message that is potentially supposed to be broadcast

		//If not a goroutine messages will be received by each client in order
		//(unless 1: there is a goroutine internally, or 2: hub.broadcast is unbuffered or is over its buffer)
		//If a goroutine, no guarantee about message order
			r.bcast(message)
		}
	}
}

func (r *Room) connect(connection *Connection) {
	r.connections.mu.Lock()
	r.connections.m[connection] = struct{}{}
	numCons := len(h.connections.m)
	r.connections.mu.Unlock()

	//Unless register and unregister have a buffer, make sure any messaging during these
	//processes is concurrent.
	go func() {
		//p, _ := Packetize("new_connection", fmt.Sprintf("%d clients currently connected to hub %s\n", numCons, h.id))
		//h.broadcast <- p
		//p, _ = Packetize("num_connections", numCons)
		//h.broadcast <- p
	}()
}