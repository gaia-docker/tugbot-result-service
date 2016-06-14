package pool

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
)

var upgrader = websocket.Upgrader{}

func getWebsocketConnMock() (*websocket.Conn, error) {

	var ret *websocket.Conn
	var err error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ret, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot upgrade: %v", err), http.StatusInternalServerError)
		}
	}))
	wsURL := strings.Replace(ts.URL, "http", "ws", 1)
	log.Infof("wsURL: %s", wsURL)
	_, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("cannot make websocket connection: %v", err)
	}

	return ret, err
}
