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
