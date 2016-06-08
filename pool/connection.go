package pool

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

// Connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type closeHandler func(*Connection)

// NewConnection creates new instance of Connection
func NewConnection(ws *websocket.Conn, send chan []byte) *Connection {

	return &Connection{send: make(chan []byte, 256), ws: ws}
}

// Listen pumps messages from the hub to the websocket connection.
func (c *Connection) Listen(onClose closeHandler) {

	defer func() {
		c.ws.Close()
		onClose(c)
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				log.Infof("Closing websocket connection...")
				return
			}
			log.Debug("Going to publish message to websocket...")
			if err := c.write(websocket.TextMessage, message); err != nil {
				log.Errorf("Failed writing to websocket: %v", err)
				return
			}
		}
	}
}

func (c *Connection) String() string {

	return fmt.Sprintf("Connection: %+v", c.ws)
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {

	return c.ws.WriteMessage(mt, payload)
}
