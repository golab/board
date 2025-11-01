package integration

import (
	"sync"

	"github.com/jarednogo/board/pkg/hub"
	"github.com/jarednogo/board/pkg/loader"
	"github.com/jarednogo/board/pkg/socket"
)

type Sim struct {
	Hub     *hub.Hub
	Clients []*socket.BlockingMockRoomConn
	wg      sync.WaitGroup
}

func NewSim() (*Sim, error) {
	ml := loader.NewMemoryLoader()
	h, err := hub.NewHubWithDB(ml)
	if err != nil {
		return nil, err
	}

	sim := &Sim{
		Hub: h,
	}

	return sim, nil
}

func (s *Sim) AddClient(roomID string) {
	client := socket.NewBlockingMockRoomConn()
	client.SetRoomID(roomID)
	s.Clients = append(s.Clients, client)
}

func (s *Sim) ConnectAll() {
	for _, client := range s.Clients {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.Hub.Handler(client, client.GetRoomID())
		}()
	}

	// block until all the clients are connected
	for _, client := range s.Clients {
		<-client.Ready()
	}
}

func (s *Sim) DisconnectAll() {
	// disconnect all the clients
	for _, client := range s.Clients {
		client.Disconnect()
	}

	// waits until all the clients are disconnected
	s.wg.Wait()
}
