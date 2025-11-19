/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package integration_test

import (
	"encoding/base64"
	"encoding/json"
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

func TestSim3(t *testing.T) {
	// make a new simulator and add some clients
	sim, err := integration.NewSim()
	assert.NoError(t, err)

	// make sure there's no rooms
	assert.Equal(t, sim.Hub.RoomCount(), 0)

	roomID := "someboard"
	sim.AddClient(roomID)

	// connect all the clients
	sim.ConnectAll()

	_, err = sim.Hub.GetRoom(roomID)
	assert.NoError(t, err)

	// disconnect all the clients
	sim.DisconnectAll()

	// save then load
	sim.Hub.Save()
	sim.Hub.Load()

	// make sure we have 1 room loaded
	assert.Equal(t, sim.Hub.RoomCount(), 1)
}

// a mirror of the test in pkg/app, but just wanted to make sure it works
func TestPing(t *testing.T) {
	sim, err := integration.NewSim()
	assert.NoError(t, err)

	body, err := sim.SendGet("/api/ping")
	assert.NoError(t, err)

	pong := struct {
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(body, &pong)
	assert.NoError(t, err)
	assert.Equal(t, pong.Message, "pong")
}

func TestVersion(t *testing.T) {
	sim, err := integration.NewSim()
	assert.NoError(t, err)

	body, err := sim.SendGet("/api/version")
	assert.NoError(t, err)

	resp := struct {
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(body, &resp)
	assert.NoError(t, err)
	assert.Equal(t, resp.Message, "test")
}

func TestUpload(t *testing.T) {
	sim, err := integration.NewSim()
	assert.NoError(t, err)
	body, err := sim.SendGet("/ext/upload")
	assert.NoError(t, err)
	// because there's no fetcher attached to the hub's rooms
	// trying to upload doesn't do anything
	assert.Equal(t, string(body), "bad request\n")
	// still, a room will be allocated because of the request
	assert.Equal(t, sim.Hub.RoomCount(), 1)
}

func TestWebRouter(t *testing.T) {
	sim, err := integration.NewSim()
	assert.NoError(t, err)

	for _, route := range []string{"/", "/about", "/integrations"} {
		body, err := sim.SendGet(route)
		assert.NoError(t, err)
		assert.True(t, len(body) > 1000)
	}
}

func TestWebRouterStateful(t *testing.T) {
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.Scoring1))
	evt := event.NewEvent("upload_sgf", sgf)

	sim, err := integration.SimWithEvents("room123", []event.Event{evt})
	assert.NoError(t, err)

	body, err := sim.SendGet("/b/room123/sgf")
	assert.NoError(t, err)
	assert.Equal(t, len(body), 2091)

	body, err = sim.SendGet("/b/room123/sgfix")
	assert.NoError(t, err)
	assert.Equal(t, len(body), 4298)

	body, err = sim.SendGet("/b/room123/debug")
	assert.NoError(t, err)
	assert.Equal(t, len(body), 8316)
}

func TestScoring1(t *testing.T) {
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.Scoring1))
	evts := []event.Event{
		event.NewEvent("upload_sgf", sgf),
		event.NewEvent("fastforward", nil),
		event.NewEvent("markdead", []any{4.0, 3.0}),
		event.NewEvent("markdead", []any{5.0, 2.0}),
		event.NewEvent("markdead", []any{6.0, 3.0}),
		event.NewEvent("markdead", []any{13.0, 0.0}),
		event.NewEvent("markdead", []any{14.0, 1.0}),
		event.NewEvent("markdead", []any{16.0, 0.0}),
		event.NewEvent("markdead", []any{11.0, 5.0}),
		event.NewEvent("markdead", []any{18.0, 5.0}),
		event.NewEvent("markdead", []any{1.0, 15.0}),
		event.NewEvent("markdead", []any{3.0, 13.0}),
		event.NewEvent("score", nil),
	}
	sim, err := integration.SimWithEvents("room123", evts)
	assert.NoError(t, err)
	saved := sim.Clients[0].SavedEvents
	final, ok := (saved[len(saved)-1].Value()).(*core.Frame)
	assert.True(t, ok)
	assert.Equal(t, final.BlackCaps, 92)
	assert.Equal(t, final.WhiteCaps, 74)
	assert.Equal(t, len(final.BlackArea), 56)
	assert.Equal(t, len(final.BlackArea), 56)
	assert.Equal(t, len(final.Dame), 7)
}

func TestScoring3(t *testing.T) {
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.Scoring3))
	evts := []event.Event{
		event.NewEvent("upload_sgf", sgf),
		event.NewEvent("fastforward", nil),
		event.NewEvent("markdead", []any{3.0, 13.0}),
		event.NewEvent("markdead", []any{5.0, 16.0}),
		event.NewEvent("score", nil),
	}
	sim, err := integration.SimWithEvents("room123", evts)
	assert.NoError(t, err)
	saved := sim.Clients[0].SavedEvents
	final, ok := (saved[len(saved)-1].Value()).(*core.Frame)
	assert.True(t, ok)
	assert.Equal(t, final.BlackCaps, 29)
	assert.Equal(t, final.WhiteCaps, 90)
	assert.Equal(t, len(final.BlackArea), 27)
	assert.Equal(t, len(final.BlackArea), 27)
	assert.Equal(t, len(final.Dame), 11)
}

func TestAddStone(t *testing.T) {
	evts := []event.Event{
		event.NewEvent("add_stone", map[string]any{
			"color":  1.0,
			"coords": []any{2.0, 2.0},
		}),
	}

	sim, err := integration.SimWithEvents("room123", evts)
	assert.NoError(t, err)
	room, err := sim.Hub.GetRoom("room123")
	assert.NoError(t, err)
	assert.True(t, room.Current().XY.Equal(&core.Coord{X: 2, Y: 2}))
}

func TestPass(t *testing.T) {
	evts := []event.Event{
		event.NewEvent("pass", 1.0),
	}

	sim, err := integration.SimWithEvents("room123", evts)
	assert.NoError(t, err)
	room, err := sim.Hub.GetRoom("room123")
	assert.NoError(t, err)
	assert.Equal(t, room.Current().XY, nil)
}

func TestRemoveStone(t *testing.T) {
	evts := []event.Event{
		event.NewEvent("add_stone", map[string]any{
			"color":  1.0,
			"coords": []any{2.0, 2.0},
		}),
		event.NewEvent("remove_stone", []any{2.0, 2.0}),
	}

	sim, err := integration.SimWithEvents("room123", evts)
	assert.NoError(t, err)
	room, err := sim.Hub.GetRoom("room123")
	assert.NoError(t, err)
	assert.Equal(t, room.Current().XY, nil)
}

func TestMarks(t *testing.T) {
	evts := []event.Event{
		event.NewEvent("triangle", []any{2.0, 2.0}),
		event.NewEvent("square", []any{3.0, 3.0}),
		event.NewEvent("letter", map[string]any{
			"coords": []any{4.0, 4.0},
			"letter": "A",
		}),
		event.NewEvent("number", map[string]any{
			"coords": []any{5.0, 5.0},
			"number": 1.0,
		}),
	}

	sim, err := integration.SimWithEvents("room123", evts)
	assert.NoError(t, err)
	room, err := sim.Hub.GetRoom("room123")
	assert.NoError(t, err)
	node := room.Current()
	assert.Equal(t, node.Fields["TR"][0], "cc")
	assert.Equal(t, node.Fields["SQ"][0], "dd")
	assert.Equal(t, node.Fields["LB"][0], "ee:A")
	assert.Equal(t, node.Fields["LB"][1], "ff:1")
}

func TestRemoveMark(t *testing.T) {
	evts := []event.Event{
		event.NewEvent("triangle", []any{2.0, 2.0}),
		event.NewEvent("remove_mark", []any{2.0, 2.0}),
	}

	sim, err := integration.SimWithEvents("room123", evts)
	assert.NoError(t, err)
	room, err := sim.Hub.GetRoom("room123")
	assert.NoError(t, err)
	node := room.Current()
	assert.Equal(t, len(node.Fields["TR"]), 0)
}

func TestNav(t *testing.T) {
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.Resignation1))
	evts := []event.Event{
		event.NewEvent("upload_sgf", sgf),
		event.NewEvent("fastforward", nil),
		event.NewEvent("rewind", nil),
		event.NewEvent("right", nil),
		event.NewEvent("right", nil),
		event.NewEvent("right", nil),
		event.NewEvent("right", nil),
		event.NewEvent("right", nil),
		event.NewEvent("left", nil),
		event.NewEvent("add_stone", map[string]any{
			"color":  1.0,
			"coords": []any{2.0, 2.0},
		}),
		event.NewEvent("left", nil),
		event.NewEvent("up", nil),
		event.NewEvent("down", nil),
		event.NewEvent("right", nil),
	}
	sim, err := integration.SimWithEvents("room123", evts)
	assert.NoError(t, err)
	room, err := sim.Hub.GetRoom("room123")
	assert.NoError(t, err)
	node := room.Current()
	assert.True(t, node.XY.Equal(&core.Coord{X: 2.0, Y: 2.0}))
}

func TestGotoGrid(t *testing.T) {
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.Resignation1))
	evts := []event.Event{
		event.NewEvent("upload_sgf", sgf),
		event.NewEvent("goto_grid", 170.0),
	}
	sim, err := integration.SimWithEvents("room123", evts)
	assert.NoError(t, err)
	room, err := sim.Hub.GetRoom("room123")
	assert.NoError(t, err)
	node := room.Current()
	assert.True(t, node.XY.Equal(&core.Coord{X: 1.0, Y: 1.0}))
}

func TestGotoCoord(t *testing.T) {
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.Resignation1))
	evts := []event.Event{
		event.NewEvent("upload_sgf", sgf),
		event.NewEvent("goto_coord", []any{9.0, 9.0}),
	}
	sim, err := integration.SimWithEvents("room123", evts)
	assert.NoError(t, err)
	room, err := sim.Hub.GetRoom("room123")
	assert.NoError(t, err)
	node := room.Current()
	assert.Equal(t, node.Index, 78)
}
