package pool

import log "github.com/Sirupsen/logrus"

// Hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the Connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection
}

// Created Hub instance
func NewHub() *Hub {

	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
	}
}

// Run hub channel loop, should be used with go routine (go hub.Run())
func (h *Hub) Run() {

	for {
		select {
		case conn := <-h.register:
			h.connections[conn] = true
			log.Infof("New connection registered. There are %d connections", len(h.connections))
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				h.closeConnection(conn)
				log.Infof("Unregistered connection. There are %d connections", len(h.connections))
			}
		case message := <-h.broadcast:
			log.Debugf("Going to broadcast message to connections: %s", message)
			for conn := range h.connections {
				select {
				case conn.send <- message:
				default:
					h.closeConnection(conn)
				}
			}
		}
	}
}

// Registers the given connection to the hub
func (h *Hub) Register(conn *Connection) {

	h.register <- conn
}

// Unregisters the given connection from the hub
func (h *Hub) Unregister(conn *Connection) {

	h.unregister <- conn
}

// Broadcast message to all hub connections
func (h *Hub) Broadcast(message *string) {

	h.broadcast <- []byte(*message)
}

func (h *Hub) closeConnection(conn *Connection) {

	delete(h.connections, conn)
	close(conn.send)
}
