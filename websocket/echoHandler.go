package websocket

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gaia-docker/tugbot-result-service/pool"
	"github.com/gorilla/websocket"
	"net/http"
)

type EchoHandler struct {
	hub *pool.Hub
}

var upgrader = websocket.Upgrader{} // use default options

// NewUploadHandler creates UploadHandler instance
func NewEchoHandler(hub *pool.Hub) *EchoHandler {

	return &EchoHandler{
		hub: hub,
	}
}

func (eh *EchoHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Errorf("Failed upgrading connection to web socket %+v", err)
		return
	}
	log.Infof("New websocket connection established")
	conn := pool.NewConnection(ws)
	eh.hub.Register(conn)
	go conn.Listen(eh.hub.Unregister)
}
