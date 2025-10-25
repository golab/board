/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room_test

import (
	"testing"

	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/fetch"
	"github.com/jarednogo/board/pkg/room"
	"github.com/jarednogo/board/pkg/room/plugin"
	"github.com/jarednogo/board/pkg/socket"
)

func TestBroadcast(t *testing.T) {
	r := room.NewRoom()
	mock1 := socket.NewMockRoomConn()
	mock2 := socket.NewMockRoomConn()

	r.RegisterConnection(mock1)
	r.RegisterConnection(mock2)

	if len(mock1.SavedEvents) != 1 || len(mock2.SavedEvents) != 1 {
		t.Errorf("expected both clients to receive the first frame event")
	}

	evt := core.EmptyEvent()
	r.Broadcast(evt)

	if len(mock1.SavedEvents) != 2 || len(mock2.SavedEvents) != 2 {
		t.Errorf("expected both clients to receive a test event")
	}
}

func TestPlugin(t *testing.T) {
	r := room.NewRoom()
	mp := plugin.NewMockPlugin()
	args := make(map[string]interface{})
	args["key"] = "mock"
	r.RegisterPlugin(mp, args)
	if _, ok := r.Plugins["mock"]; !ok {
		t.Fatalf("failed to register mock plugin")
	}

	r.DeregisterPlugin("mock")
	if _, ok := r.Plugins["mock"]; ok {
		t.Errorf("failed to deregister mock plugin")
	}
}

func TestSendUserList(t *testing.T) {
	r := room.NewRoom()
	mock := socket.NewMockRoomConn()
	r.RegisterConnection(mock)
	r.SendUserList()
	// first event is the frame immediately on connecting
	if len(mock.SavedEvents) != 2 {
		t.Fatalf("failed to send user list")
	}

	if mock.SavedEvents[1].Event != "connected_users" {
		t.Errorf("failed to send correct event")
	}
}

func TestSendTo(t *testing.T) {
	r := room.NewRoom()
	mock := socket.NewMockRoomConn()
	id := r.RegisterConnection(mock)
	evt := core.EmptyEvent()
	r.SendTo(id, evt)

	if len(mock.SavedEvents) != 2 {
		t.Errorf("expected client to receive a test event after SendTo")
	}
}

func TestHandlers(t *testing.T) {
	r := room.NewRoom()
	mock := socket.NewMockRoomConn()
	id := r.RegisterConnection(mock)
	handlers := r.CreateHandlers()
	evt := core.EmptyEvent()
	handlers["ping"](evt)

	evt.UserID = id
	handlers["isprotected"](evt)
}

func TestFetcher(t *testing.T) {
	r := room.NewRoom()
	r.SetFetcher(fetch.NewMockFetcher())
}
