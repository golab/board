/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package integration_test

import (
	"encoding/base64"
	"testing"

	"github.com/jarednogo/board/integration"
	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/internal/sgfsamples"
	"github.com/jarednogo/board/pkg/core"
)

func TestSim(t *testing.T) {
	// make a new simulator and add some clients
	sim, err := integration.NewSim()
	assert.NoError(t, err, "test sim")

	roomID := "someboard"
	for i := 0; i < 10; i++ {
		sim.AddClient(roomID)
	}

	// connect all the clients
	sim.ConnectAll()

	// clear out all events
	sim.FlushAll()

	// reduce input buffer to 0 for this test
	// otherwise the buffer causes events to get dropped
	// by the room handler
	room, err := sim.Hub.GetRoom(roomID)
	assert.NoError(t, err, "get room")
	room.SetInputBuffer(0)

	// simulate an event
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.PassWithTT))
	evt := &core.Event{
		Type:   "upload_sgf",
		Value:  sgf,
		UserID: "",
	}
	sim.Clients[0].SimulateEvent(evt)

	// let the event pass through all connections
	sim.FlushAll()

	// disconnect all the clients
	sim.DisconnectAll()

	// observe effects
	save := room.Save()
	assert.Equal(t, len(save.SGF), 6132, "len(save.SGF)")
}
