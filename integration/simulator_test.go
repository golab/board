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
	"github.com/jarednogo/board/pkg/event"
)

func TestSim1(t *testing.T) {
	// make a new simulator and add some clients
	sim, err := integration.NewSim()
	assert.NoError(t, err)

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
	assert.NoError(t, err)
	room.DisableBuffers()

	// simulate an event
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.PassWithTT))
	evt := event.NewEvent("upload_sgf", sgf)
	sim.Clients[0].SimulateEvent(evt)

	// let the event pass through all connections
	sim.FlushAll()

	// disconnect all the clients
	sim.DisconnectAll()

	// observe effects
	save := room.Save()
	assert.Equal(t, len(save.SGF), 6132)
}

func TestSim2(t *testing.T) {
	// make a new simulator and add some clients
	sim, err := integration.NewSim()
	assert.NoError(t, err)

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
	assert.NoError(t, err)
	room.DisableBuffers()

	p := core.NewParser(sgfsamples.Scoring2)
	root, err := p.Parse()
	assert.NoError(t, err)
	cur := root.Down[0]
	for len(cur.Down) != 0 {
		// simulate an event
		value := make(map[string]any)
		color := 1.0
		key := "B"
		if _, ok := cur.Fields["W"]; ok {
			key = "W"
			color = 2.0
		}
		value["color"] = color

		coord := core.LettersToCoord(cur.Fields[key][0])

		value["coords"] = []any{float64(coord.X), float64(coord.Y)}
		evt := event.NewEvent("add_stone", value)
		sim.Clients[0].SimulateEvent(evt)

		// let the event pass through all connections
		sim.FlushAll()

		cur = cur.Down[0]
	}

	// disconnect all the clients
	sim.DisconnectAll()

	// observe effects
	save := room.Save()
	assert.Equal(t, len(save.SGF), 4024)

}
