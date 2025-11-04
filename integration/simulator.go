/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package integration

import (
	"sync"

	"github.com/jarednogo/board/pkg/config"
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
	h, err := hub.NewHubWithDB(ml, config.Default())
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
