package integration_test

import (
	"testing"

	"github.com/jarednogo/board/integration"
	"github.com/jarednogo/board/internal/assert"
)

func TestSim(t *testing.T) {
	// make a new simulator and add some clients
	sim, err := integration.NewSim()
	assert.NoError(t, err, "test sim")

	roomID := "someboard"
	for i := 0; i < 10; i++ {
		sim.AddClient(roomID)
	}

	// connect all the clients
	sim.ConnectAll()

	/*
		do stuff with the clients
	*/

	// disconnect all the clients
	sim.DisconnectAll()
}
