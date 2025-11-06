/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core

import (
	"io"
	"sync"

	"github.com/google/uuid"
)

type MockEventChannel struct {
	QueuedEvents []*Event
	index        int
	SavedEvents  []*Event
	roomID       string
	id           string
	Closed       bool
	mu           sync.Mutex
}

func NewMockEventChannel() *MockEventChannel {
	id := uuid.New().String()
	return &MockEventChannel{id: id}
}

func (mcr *MockEventChannel) SetRoomID(s string) {
	mcr.roomID = s
}

func (mcr *MockEventChannel) GetRoomID() string {
	return mcr.roomID
}

func (mcr *MockEventChannel) OnConnect() {
}

func (mcr *MockEventChannel) SendEvent(evt *Event) error {
	mcr.mu.Lock()
	defer mcr.mu.Unlock()
	mcr.SavedEvents = append(mcr.SavedEvents, evt)
	return nil
}

func (mcr *MockEventChannel) ReceiveEvent() (*Event, error) {
	if mcr.index >= len(mcr.QueuedEvents) {
		return nil, io.EOF
	}
	i := mcr.index
	mcr.index++
	return mcr.QueuedEvents[i], nil
}

func (mcr *MockEventChannel) Close() error {
	mcr.Closed = true
	return nil
}

func (mcr *MockEventChannel) ID() string {
	return mcr.id
}

type BlockingMockEventChannel struct {
	conn  chan bool
	ready chan bool
	*MockEventChannel
}

func NewBlockingMockEventChannel() *BlockingMockEventChannel {
	return &BlockingMockEventChannel{
		make(chan bool),
		make(chan bool),
		NewMockEventChannel(),
	}
}

func (mcr *BlockingMockEventChannel) Ready() <-chan bool {
	return mcr.ready
}

func (mcr *BlockingMockEventChannel) OnConnect() {
	close(mcr.ready)
}

func (mcr *BlockingMockEventChannel) Disconnect() {
	mcr.conn <- true
}

func (mcr *BlockingMockEventChannel) ReceiveEvent() (*Event, error) {
	if mcr.index >= len(mcr.QueuedEvents) {
		// blocks until there's a value from mcr.conn
		<-mcr.conn
	}
	return mcr.MockEventChannel.ReceiveEvent()
}

func (mcr *BlockingMockEventChannel) SendEvent(evt *Event) error {
	// signals to external caller that the room conn is ready
	return mcr.MockEventChannel.SendEvent(evt)
}
