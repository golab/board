/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package socket

import (
	"io"

	"github.com/google/uuid"
	"github.com/jarednogo/board/pkg/core"
)

type MockRoomConn struct {
	QueuedEvents []*core.EventJSON
	index        int
	SavedEvents  []*core.EventJSON
	id           string
	Closed       bool
}

func NewMockRoomConn() *MockRoomConn {
	id := uuid.New().String()
	return &MockRoomConn{id: id}
}

func (mcr *MockRoomConn) SendEvent(evt *core.EventJSON) error {
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
