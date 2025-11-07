/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package event

import (
	"fmt"
	"io"
	"sync"

	"github.com/google/uuid"
)

type MockEventChannel struct {
	QueuedEvents []Event
	index        int
	SavedEvents  []Event
	roomID       string
	id           string
	Closed       bool
	mu           sync.Mutex
}

func NewMockEventChannel() *MockEventChannel {
	id := uuid.New().String()
	return &MockEventChannel{id: id}
}

func (ec *MockEventChannel) SetRoomID(s string) {
	ec.roomID = s
}

func (ec *MockEventChannel) GetRoomID() string {
	return ec.roomID
}

func (ec *MockEventChannel) OnConnect() {
}

func (ec *MockEventChannel) SendEvent(evt Event) error {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.SavedEvents = append(ec.SavedEvents, evt)
	return nil
}

func (ec *MockEventChannel) ReceiveEvent() (Event, error) {
	if ec.index >= len(ec.QueuedEvents) {
		return nil, io.EOF
	}
	i := ec.index
	ec.index++
	return ec.QueuedEvents[i], nil
}

func (ec *MockEventChannel) Close() error {
	ec.Closed = true
	return nil
}

func (ec *MockEventChannel) ID() string {
	return ec.id
}

// BlockingMockEventChannel essentially does the same
// thing as a MockEventChannel except when it's out of QueuedEvents
// then it just blocks until Disconnect() is called.
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

func (ec *BlockingMockEventChannel) Ready() <-chan bool {
	return ec.ready
}

func (ec *BlockingMockEventChannel) OnConnect() {
	close(ec.ready)
}

func (ec *BlockingMockEventChannel) Disconnect() {
	close(ec.conn)
}

func (ec *BlockingMockEventChannel) ReceiveEvent() (Event, error) {
	if ec.index >= len(ec.QueuedEvents) {
		// blocks until there's a value from ec.conn (or ec.conn closes)
		<-ec.conn
	}
	return ec.MockEventChannel.ReceiveEvent()
}

type TwoWayMockEventChannel struct {
	sentEvents     chan Event
	receivedEvents chan Event
	*BlockingMockEventChannel
}

func NewTwoWayMockEventChannel() *TwoWayMockEventChannel {
	return &TwoWayMockEventChannel{
		make(chan Event, 50),
		make(chan Event),
		NewBlockingMockEventChannel(),
	}
}

func (ec *TwoWayMockEventChannel) SendEvent(evt Event) error {
	go func() {
		ec.sentEvents <- evt
	}()
	return ec.BlockingMockEventChannel.SendEvent(evt)
}

func (ec *TwoWayMockEventChannel) Disconnect() {
	close(ec.receivedEvents)
	ec.BlockingMockEventChannel.Disconnect()
}

func (ec *TwoWayMockEventChannel) ReceiveEvent() (Event, error) {
	evt, ok := <-ec.receivedEvents
	if !ok {
		return nil, fmt.Errorf("channel closed")
	}
	return evt, nil
}

func (ec *TwoWayMockEventChannel) SimulateEvent(evt Event) {
	ec.receivedEvents <- evt
}

func (ec *TwoWayMockEventChannel) Flush() {
	for {
		select {
		case <-ec.sentEvents:
		default:
			return
		}
	}
}
