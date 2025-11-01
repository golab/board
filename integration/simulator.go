package integration

import (
	"github.com/jarednogo/board/pkg/hub"
	"github.com/jarednogo/board/pkg/loader"
	"github.com/jarednogo/board/pkg/socket"
)

type Sim struct {
	Hub     *hub.Hub
	Clients []*socket.BlockingMockRoomConn
}

func NewSim(numClients int) (*Sim, error) {
	ml := loader.NewMemoryLoader()
	h, err := hub.NewHubWithDB(ml)
	if err != nil {
		return nil, err
	}

	clients := []*socket.BlockingMockRoomConn{}
	for i := 0; i < numClients; i++ {
		clients = append(clients, socket.NewBlockingMockRoomConn())
	}

	sim := &Sim{
		Hub:     h,
		Clients: clients,
	}

	return sim, nil
}
