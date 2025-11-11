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

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/event"
	"github.com/jarednogo/board/pkg/fetch"
	"github.com/jarednogo/board/pkg/message"
	"github.com/jarednogo/board/pkg/room"
	"github.com/jarednogo/board/pkg/room/plugin"
)

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

	message := message.NewMessage("foobar", 30)
	r.BroadcastHubMessage(message)

	assert.Equal(t, len(mock1.SavedEvents), 2)
	assert.Equal(t, len(mock2.SavedEvents), 2)
}

func TestPlugin(t *testing.T) {
	r := room.NewRoom("")
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
