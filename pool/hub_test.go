package pool

import (
	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {

	var hub = NewHub()
	log.Info("TestRegister")
	hub.Register(NewConnection(nil))
	assert.Len(t, hub.connections, 1, "Number of connections should be 1")
}

func TestUnregister(t *testing.T) {

	var hub = NewHub()
	log.Info("TestUnregister")
	connection1 := NewConnection(nil)
	connection2 := NewConnection(nil)
	hub.Register(connection1)
	hub.Register(connection2)
	assert.Len(t, hub.connections, 2, "Number of connections should be 2")
	hub.Unregister(connection1)
	assert.Len(t, hub.connections, 1, "Number of connections should be 1")
	hub.Unregister(connection2)
	assert.Len(t, hub.connections, 0, "Number of connections should be 0")
}

func TestUnregisterNoConnections(t *testing.T) {

	var hub = NewHub()
	log.Info("TestUnregisterNoConnections")
	connection1 := NewConnection(nil)
	hub.Unregister(connection1)
	assert.Len(t, hub.connections, 0, "Number of connections should be 0")
}

func TestBroadcast(t *testing.T) {

	var hub = NewHub()
	ws, err := GetWebsocketConnMock()
	if err != nil {
		t.Fatal(err)
	}
	log.Info("TestBroadcast")
	go hub.Run()
	connection1 := NewConnection(ws)
	go connection1.Listen(hub.Unregister)
	hub.Register(connection1)
	message := "hello world"
	hub.Broadcast(&message)
	hub.Unregister(connection1)
}
