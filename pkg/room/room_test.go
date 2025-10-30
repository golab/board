/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room_test

import (
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/fetch"
	"github.com/jarednogo/board/pkg/room"
	"github.com/jarednogo/board/pkg/room/plugin"
	"github.com/jarednogo/board/pkg/socket"
)

func TestBroadcast(t *testing.T) {
	r := room.NewRoom("")
	mock1 := socket.NewMockRoomConn()
	mock2 := socket.NewMockRoomConn()

	r.RegisterConnection(mock1)
	r.RegisterConnection(mock2)

	assert.Equal(t, len(mock1.SavedEvents), 1, "expected mock1 to receive first frame event")
	assert.Equal(t, len(mock2.SavedEvents), 1, "expected mock2 to receive first frame event")

	evt := core.EmptyEvent()
	r.Broadcast(evt)

	assert.Equal(t, len(mock1.SavedEvents), 2, "expected mock1 to recieve a test event")
	assert.Equal(t, len(mock2.SavedEvents), 2, "expected mock2 to recieve a test event")
}

func TestPlugin(t *testing.T) {
	r := room.NewRoom("")
	mp := plugin.NewMockPlugin()
	args := make(map[string]interface{})
	args["key"] = "mock"
	r.RegisterPlugin(mp, args)
	ok := r.HasPlugin("mock")
	assert.True(t, ok, "failed to register mock plugin")

	r.DeregisterPlugin("mock")
	ok = r.HasPlugin("mock")
	assert.True(t, !ok, "failed to deregister mock plugin")
}

func TestSendUserList(t *testing.T) {
	r := room.NewRoom("")
	mock := socket.NewMockRoomConn()
	r.RegisterConnection(mock)
	r.SendUserList()
	// first event is the frame immediately on connecting
	if len(mock.SavedEvents) != 2 {
		t.Fatalf("failed to send user list")
	}

	assert.Equal(t, mock.SavedEvents[1].Event, "connected_users", "failed to send correct event")
}

func TestSendTo(t *testing.T) {
	r := room.NewRoom("")
	mock := socket.NewMockRoomConn()
	id := r.RegisterConnection(mock)
	evt := core.EmptyEvent()
	r.SendTo(id, evt)

	assert.Equal(t, len(mock.SavedEvents), 2, "expected client to receive a test event")
}

func TestHandlers(t *testing.T) {
	r := room.NewRoom("")
	mock := socket.NewMockRoomConn()
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

func TestFetcher(t *testing.T) {
	r := room.NewRoom("")
	r.SetFetcher(fetch.NewMockFetcher())
}
