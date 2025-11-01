/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package socket

import (
	"io"
	"sync"

	"github.com/google/uuid"
	"github.com/jarednogo/board/pkg/core"
)

type MockRoomConn struct {
	QueuedEvents []*core.EventJSON
	index        int
	SavedEvents  []*core.EventJSON
	roomID       string
	id           string
	Closed       bool
	mu           sync.Mutex
}

func NewMockRoomConn() *MockRoomConn {
	id := uuid.New().String()
	return &MockRoomConn{id: id}
}

func (mcr *MockRoomConn) SetRoomID(s string) {
	mcr.roomID = s
}

func (mcr *MockRoomConn) GetRoomID() string {
	return mcr.roomID
}

func (mcr *MockRoomConn) OnConnect() {
}

func (mcr *MockRoomConn) SendEvent(evt *core.EventJSON) error {
	mcr.mu.Lock()
	defer mcr.mu.Unlock()
	mcr.SavedEvents = append(mcr.SavedEvents, evt)
	return nil
}

func (mcr *MockRoomConn) ReceiveEvent() (*core.EventJSON, error) {
	if mcr.index >= len(mcr.QueuedEvents) {
		return nil, io.EOF
	}
	i := mcr.index
	mcr.index++
	return mcr.QueuedEvents[i], nil
}

func (mcr *MockRoomConn) Close() error {
	mcr.Closed = true
	return nil
}

func (mcr *MockRoomConn) ID() string {
	return mcr.id
}

type BlockingMockRoomConn struct {
	conn  chan bool
	ready chan bool
	*MockRoomConn
}

func NewBlockingMockRoomConn() *BlockingMockRoomConn {
	return &BlockingMockRoomConn{
		make(chan bool),
		make(chan bool),
		NewMockRoomConn(),
	}
}

func (mcr *BlockingMockRoomConn) Ready() <-chan bool {
	return mcr.ready
}

func (mcr *BlockingMockRoomConn) OnConnect() {
	close(mcr.ready)
}

func (mcr *BlockingMockRoomConn) Disconnect() {
	mcr.conn <- true
}

func (mcr *BlockingMockRoomConn) ReceiveEvent() (*core.EventJSON, error) {
	if mcr.index >= len(mcr.QueuedEvents) {
		// blocks until there's a value from mcr.conn
		<-mcr.conn
	}
	return mcr.MockRoomConn.ReceiveEvent()
}

func (mcr *BlockingMockRoomConn) SendEvent(evt *core.EventJSON) error {
	// signals to external caller that the room conn is ready
	return mcr.MockRoomConn.SendEvent(evt)
}
