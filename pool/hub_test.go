package pool

import (
	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {

	var hub = NewHub()
	log.Info("TestRegister")
	go hub.Run()
	hub.Register(NewConnection(nil, nil))
	assert.Len(t, hub.connections, 1, "Number of connections should be 1")
}

func TestUnregister(t *testing.T) {

	var hub = NewHub()
	log.Info("TestUnregister")
	go hub.Run()
	connection1 := NewConnection(nil, nil)
	connection2 := NewConnection(nil, nil)
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
	go hub.Run()
	connection1 := NewConnection(nil, nil)
	hub.Unregister(connection1)
	assert.Len(t, hub.connections, 0, "Number of connections should be 0")
}
