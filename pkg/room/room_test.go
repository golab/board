/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room_test

import (
	"io"
	"testing"
	"time"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/fetch"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
	"github.com/golab/board/pkg/event"
	"github.com/golab/board/pkg/logx"
	"github.com/golab/board/pkg/message"
	"github.com/golab/board/pkg/room"
	"github.com/golab/board/pkg/room/plugin"
	"github.com/golab/board/pkg/state"
)

func TestNumConns(t *testing.T) {
	r := room.NewRoom("")
	mock1 := event.NewMockEventChannel()
	mock2 := event.NewMockEventChannel()
	mock3 := event.NewMockEventChannel()

	assert.Equal(t, r.NumConns(), 0)
	id1 := r.RegisterConnection(mock1)
	assert.Equal(t, r.NumConns(), 1)
	id2 := r.RegisterConnection(mock2)
	assert.Equal(t, r.NumConns(), 2)
	id3 := r.RegisterConnection(mock3)
	assert.Equal(t, r.NumConns(), 3)

	r.DeregisterConnection(id1)
	assert.Equal(t, r.NumConns(), 2)
	r.DeregisterConnection(id2)
	assert.Equal(t, r.NumConns(), 1)
	r.DeregisterConnection(id3)
	assert.Equal(t, r.NumConns(), 0)
}

func TestBroadcast(t *testing.T) {
	r := room.NewRoom("")
	mock1 := event.NewMockEventChannel()
	mock2 := event.NewMockEventChannel()

	r.RegisterConnection(mock1)
	r.RegisterConnection(mock2)

	assert.Equal(t, len(mock1.SavedEvents), 1)
	assert.Equal(t, len(mock2.SavedEvents), 1)

	evt := event.EmptyEvent()
	r.Broadcast(evt)

	assert.Equal(t, len(mock1.SavedEvents), 2)
	assert.Equal(t, len(mock2.SavedEvents), 2)
}

func TestBroadcastMessage(t *testing.T) {
	r := room.NewRoom("")
	mock1 := event.NewMockEventChannel()
	mock2 := event.NewMockEventChannel()

	r.RegisterConnection(mock1)
	r.RegisterConnection(mock2)

	message := message.New("foobar", 30)
	r.BroadcastHubMessage(message)

	assert.Equal(t, len(mock1.SavedEvents), 2)
	assert.Equal(t, len(mock2.SavedEvents), 2)
}

func TestBroadcastFullFrame(t *testing.T) {
	r := room.NewRoom("")
	mock1 := event.NewMockEventChannel()
	r.RegisterConnection(mock1)

	r.BroadcastFullFrame()

	require.Equal(t, len(mock1.SavedEvents), 2)
	assert.Equal(t, mock1.SavedEvents[1].Type(), "frame")
}

func TestBroadcastTreeOnly(t *testing.T) {
	r := room.NewRoom("")
	mock1 := event.NewMockEventChannel()
	r.RegisterConnection(mock1)

	r.BroadcastTreeOnly()

	require.Equal(t, len(mock1.SavedEvents), 2)
	assert.Equal(t, mock1.SavedEvents[1].Type(), "frame")
}

func TestBroadcastNop(t *testing.T) {
	r := room.NewRoom("")
	mock1 := event.NewMockEventChannel()
	r.RegisterConnection(mock1)

	evt := event.NopEvent()
	r.Broadcast(evt)

	assert.Equal(t, len(mock1.SavedEvents), 1)
}

func TestPlugin(t *testing.T) {
	r := room.NewRoom("")
	logger := logx.NewRecorder(logx.LogLevelInfo)
	r.SetLogger(logger)
	mp := plugin.NewMockPlugin()
	args := make(map[string]any)
	args["key"] = "mock"
	r.RegisterPlugin(mp, args)
	ok := r.HasPlugin("mock")
	assert.True(t, ok)

	r.DeregisterPlugin("mock")
	ok = r.HasPlugin("mock")
	assert.True(t, !ok)
}

func TestSendUserList(t *testing.T) {
	r := room.NewRoom("")
	mock := event.NewMockEventChannel()
	r.RegisterConnection(mock)
	r.SendUserList()
	// first event is the frame immediately on connecting
	if len(mock.SavedEvents) != 2 {
		t.Fatalf("failed to send user list")
	}

	assert.Equal(t, mock.SavedEvents[1].Type(), "connected_users")
}

func TestSendTo(t *testing.T) {
	r := room.NewRoom("")
	mock := event.NewMockEventChannel()
	id := r.RegisterConnection(mock)
	evt := event.EmptyEvent()
	r.SendTo(id, evt)

	assert.Equal(t, len(mock.SavedEvents), 2)
}

/*
func TestHandlers(t *testing.T) {
	r := room.NewRoom("")
	mock := core.NewMockEventChannel()
	id := r.RegisterConnection(mock)
	handlers := r.CreateHandlers()
	evt := core.EmptyEvent()
	handlers["ping"](evt)

	evt.UserID = id
	handlers["isprotected"](evt)

	evt.Value = "foo"
	resp := handlers["checkpassword"](evt)
	assert.Zero(t, resp.Value.(string), "password string should be empty")
}
*/

func TestFetcher(t *testing.T) {
	r := room.NewRoom("")
	r.SetFetcher(fetch.NewMockFetcher(""))
}

func TestSaveLoad(t *testing.T) {
	r := room.NewRoom("")
	r.SetPassword("foobar")
	l := r.Save()
	r2, err := room.Load(l)
	assert.NoError(t, err)
	assert.True(t, r2.HasPassword())
	assert.Equal(t, r2.GetPassword(), "foobar")
}

func TestClose(t *testing.T) {
	r := room.NewRoom("")
	mock1 := event.NewMockEventChannel()
	mock2 := event.NewMockEventChannel()
	mock3 := event.NewMockEventChannel()
	r.RegisterConnection(mock1)
	r.RegisterConnection(mock2)
	r.RegisterConnection(mock3)
	err := r.Close()
	assert.NoError(t, err)
	assert.True(t, mock1.Closed)
	assert.True(t, mock2.Closed)
	assert.True(t, mock3.Closed)
}

func TestHandle(t *testing.T) {
	r := room.NewRoom("")
	// mock initializes with zero events so automatically sends an error
	mock := event.NewMockEventChannel()
	err := r.Handle(mock)
	assert.ErrorIs(t, err, io.EOF)
}

func TestID(t *testing.T) {
	r := room.NewRoom("abc")
	id := r.ID()
	assert.Equal(t, id, "abc")
}

func TestGetLastActive(t *testing.T) {
	r := room.NewRoom("")
	tm := r.GetLastActive()
	s := time.Since(tm)
	assert.True(t, s > 0)
}

func TestSaveState(t *testing.T) {
	r := room.NewRoom("")
	s := r.SaveState()
	assert.Equal(t, len(s.SGF), 96)
}

func TestTimeout(t *testing.T) {
	r := room.NewRoom("")
	r.SetTimeout(3.14)
	assert.Equal(t, r.GetTimeout(), 3.14)
}

func TestGenerateTreeOnly(t *testing.T) {
	r := room.NewRoom("")
	f := r.GenerateTreeOnly(state.Full)
	assert.Equal(t, f.Type, state.DiffFrame)
}

func TestAddStones(t *testing.T) {
	r := room.NewRoom("")
	stones := []*coord.Stone{}
	stones = append(stones, coord.NewStone(2, 3, color.Black))
	stones = append(stones, coord.NewStone(14, 3, color.White))
	r.AddStones(stones)
	s := r.GetState()
	assert.Equal(t, s.GetNextIndex(), 3)
}

func TestAddStonesToTrunk(t *testing.T) {
	r := room.NewRoom("")
	stones := []*coord.Stone{}
	stones = append(stones, coord.NewStone(2, 3, color.Black))
	stones = append(stones, coord.NewStone(14, 3, color.White))
	stones = append(stones, coord.NewStone(3, 14, color.Black))
	stones = append(stones, coord.NewStone(14, 14, color.White))
	r.AddStones(stones)

	tstones := []*coord.Stone{}
	tstones = append(tstones, coord.NewStone(4, 14, color.Black))
	tstones = append(tstones, coord.NewStone(15, 14, color.White))
	r.AddStonesToTrunk(2, tstones)

	s := r.GetState()
	assert.Equal(t, s.GetNextIndex(), 7)
}

func TestGetColorAt(t *testing.T) {
	r := room.NewRoom("")
	r.PushHead(2, 3, color.Black)
	r.PushHead(14, 3, color.White)
	assert.Equal(t, r.GetColorAt(2), color.White)
}

func TestHeadColor(t *testing.T) {
	r := room.NewRoom("")
	assert.Equal(t, r.HeadColor(), color.Empty)
	r.PushHead(2, 3, color.Black)
	assert.Equal(t, r.HeadColor(), color.Black)
	r.PushHead(14, 3, color.White)
	assert.Equal(t, r.HeadColor(), color.White)
}

func TestToSGFIX(t *testing.T) {
	r := room.NewRoom("")
	r.PushHead(2, 3, color.Black)
	r.PushHead(14, 3, color.White)
	r.PushHead(3, 14, color.Black)
	r.PushHead(14, 14, color.White)
	assert.Equal(t, len(r.ToSGF()), 89)
	assert.Equal(t, len(r.ToSGFIX()), 114)
}

func TestBoard(t *testing.T) {
	r := room.NewRoom("")
	r.PushHead(2, 3, color.Black)
	r.PushHead(14, 3, color.White)
	r.PushHead(3, 14, color.Black)
	r.PushHead(14, 14, color.White)
	b := r.Board()
	assert.Equal(t, b.Get(coord.NewCoord(2, 3)), color.Black)
	assert.Equal(t, b.Get(coord.NewCoord(14, 3)), color.White)
	assert.Equal(t, b.Get(coord.NewCoord(3, 14)), color.Black)
	assert.Equal(t, b.Get(coord.NewCoord(14, 14)), color.White)
}
