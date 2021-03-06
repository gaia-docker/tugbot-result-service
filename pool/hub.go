package pool

import (
	log "github.com/Sirupsen/logrus"
	"sync"
)

// Hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the Connections.
	broadcast chan []byte

	mutex *sync.Mutex
}

// NewHub creates Hub instance
func NewHub() *Hub {

	return &Hub{
		broadcast:   make(chan []byte),
		connections: make(map[*Connection]bool),
		mutex:       &sync.Mutex{},
	}
}

// Run hub channel loop, should be used with go routine (go hub.Run())
func (h *Hub) Run() {

	defer func() {
		for conn := range h.connections {
			h.closeConnection(conn)
		}
	}()
	for {
		select {
		case message, ok := <-h.broadcast:
			if !ok {
				log.Infof("Hub is shutting down..")
				return
			}
			log.Infof("Going to broadcast message to connections: %s", message)
			h.mutex.Lock()
			for conn := range h.connections {
				select {
				case conn.send <- message:
				default:
					h.closeConnection(conn)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// Register the given connection to the hub
func (h *Hub) Register(conn *Connection) {

	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.connections[conn] = true
	log.Infof("New connection registered. There are %d connections", len(h.connections))
}

// Unregister the given connection from the hub
func (h *Hub) Unregister(conn *Connection) {

	h.mutex.Lock()
	defer h.mutex.Unlock()
	if _, ok := h.connections[conn]; ok {
		h.closeConnection(conn)
		log.Infof("Unregistered connection. There are %d connections", len(h.connections))
	} else {
		log.Warningf("Connection <%+v> is not registered", conn)
	}
}

// Broadcast message to all hub connections
func (h *Hub) Broadcast(message *string) {

	h.broadcast <- []byte(*message)
}

// CloseBroadcastChannel closes broadcast channel
func (h *Hub) CloseBroadcastChannel() {

	close(h.broadcast)
}

func (h *Hub) closeConnection(conn *Connection) {

	log.Infof("Closing connection: %+v", conn)
	delete(h.connections, conn)
	close(conn.send)
}
