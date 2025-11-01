package integration_test

import (
	"sync"
	"testing"

	"github.com/jarednogo/board/integration"
	"github.com/jarednogo/board/internal/assert"
)

func TestSim(t *testing.T) {
	sim, err := integration.NewSim(10)
	assert.NoError(t, err, "test sim")

	var wg sync.WaitGroup
	wg.Add(len(sim.Clients))

	roomID := "someboard"

	// connect all the clients
	for _, client := range sim.Clients {
		go func() {
			defer wg.Done()
			sim.Hub.Handler(client, roomID)
		}()
	}

	// block until all the clients are connected
	for _, client := range sim.Clients {
		<-client.Ready()
	}

	// disconnect all the clients
	for _, client := range sim.Clients {
		client.Disconnect()
	}

	// waits until all the clients are disconnected
	wg.Wait()
}
