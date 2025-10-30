/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub_test

import (
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/hub"
	"github.com/jarednogo/board/pkg/loader"
	"github.com/jarednogo/board/pkg/socket"
)

func TestHub1(t *testing.T) {
	h, err := hub.NewHubWithDB(loader.NewMemoryLoader())
	assert.NoError(t, err, "new hub")
	h.Load()

	mock := socket.NewMockRoomConn()
	roomID := "someboard"
	h.Handler(mock, roomID)

	h.Save()
}

func TestHub2(t *testing.T) {
	h, err := hub.NewHubWithDB(loader.NewMemoryLoader())
	assert.NoError(t, err, "new hub")

	mock := socket.NewMockRoomConn()
	mock.QueuedEvents = append(mock.QueuedEvents,
		&core.EventJSON{
			Event: "pass",
			Value: 1.0,
		},
	)
	roomID := "someboard"
	h.Handler(mock, roomID)

	// 3 events:
	// initial frame
	// connected users
	// pass event
	assert.Equal(t, len(mock.SavedEvents), 3, "mock.receivedEvents")
}
