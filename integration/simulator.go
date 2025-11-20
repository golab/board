/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package integration

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/jarednogo/board/pkg/app"
	"github.com/jarednogo/board/pkg/config"
	"github.com/jarednogo/board/pkg/event"
	"github.com/jarednogo/board/pkg/hub"
	"github.com/jarednogo/board/pkg/logx"
)

type Sim struct {
	Hub     *hub.Hub
	Clients []*event.TwoWayMockEventChannel
	wg      sync.WaitGroup
	router  http.Handler
	Logger  *logx.Recorder
}

func NewSim() (*Sim, error) {
	l := logx.NewRecorder(logx.LogLevelInfo)
	a, err := app.New(config.Test(), l)
	if err != nil {
		return nil, err
	}

	sim := &Sim{
		Hub:    a.Hub,
		router: a.Router,
		Logger: l,
	}

	return sim, nil
}

func (s *Sim) SendGet(route string) ([]byte, error) {
	req := httptest.NewRequest("GET", route, nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	return io.ReadAll(rec.Body)
}

func (s *Sim) SendPost(route string, body []byte) ([]byte, error) {
	req := httptest.NewRequest("POST", route, bytes.NewBuffer(body))
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	return io.ReadAll(rec.Body)
}

func (s *Sim) AddClient(roomID string) {
	client := event.NewTwoWayMockEventChannel()
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

func (s *Sim) FlushAll() {
	for _, client := range s.Clients {
		client.Flush()
	}
}

func SimWithEvents(roomID string, evts []event.Event) (*Sim, error) {
	// make a new simulator and add some clients
	sim, err := NewSim()
	if err != nil {
		return nil, err
	}

	sim.AddClient(roomID)

	// connect all the clients
	sim.ConnectAll()

	room, err := sim.Hub.GetRoom(roomID)
	if err != nil {
		return nil, err
	}
	room.DisableBuffers()

	for _, evt := range evts {
		sim.Clients[0].SimulateEvent(evt)

		// let the event pass through all connections
		sim.FlushAll()
	}

	// disconnect all the clients
	sim.DisconnectAll()

	// let the hub close by waiting for all the room handlers to finish
	sim.Hub.Close()

	return sim, nil
}
