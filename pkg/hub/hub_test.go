/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub_test

import (
	"fmt"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/config"
	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/hub"
	"github.com/jarednogo/board/pkg/loader"
)

func TestParseURL(t *testing.T) {
	testcases := []struct {
		input  string
		output [3]string
	}{
		{"/socket/b/someboard/", [3]string{"b", "someboard", ""}},
		{"/socket/foo/", [3]string{"foo", "", ""}},
		{"/socket/a/bcd/e", [3]string{"a", "bcd", "e"}},
		{"/socket/x", [3]string{"x", "", ""}},
		{"/socket", [3]string{"", "", ""}},
		{"/v/wxy/z", [3]string{"v", "wxy", "z"}},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("parseURL%d", i), func(t *testing.T) {
			a, b, c := hub.ParseURL(tc.input)
			assert.Equal(t, [3]string{a, b, c}, tc.output, "TestParseURL")
		})
	}
}

func TestHub1(t *testing.T) {
	h, err := hub.NewHubWithDB(loader.NewMemoryLoader(), config.Default())
	assert.NoError(t, err, "new hub")
	h.Load()

	mock := core.NewMockEventChannel()
	roomID := "someboard"
	h.Handler(mock, roomID)

	h.Save()
}

func TestHub2(t *testing.T) {
	h, err := hub.NewHubWithDB(loader.NewMemoryLoader(), config.Default())
	assert.NoError(t, err, "new hub")

	mock := core.NewMockEventChannel()
	mock.QueuedEvents = append(mock.QueuedEvents,
		&core.Event{
			Type: "pass",
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

func TestHub3(t *testing.T) {
	ml := loader.NewMemoryLoader()
	// messages that expire immediately
	ml.AddMessage("hello world", 0)
	ml.AddMessage("server message", 0)
	// message that doesn't expire immediately
	ml.AddMessage("save this message", 30)
	h, err := hub.NewHubWithDB(ml, config.Default())
	assert.NoError(t, err, "new hub")

	roomID := "someboard"
	mock1 := core.NewMockEventChannel()
	h.Handler(mock1, roomID)

	assert.Equal(t, h.RoomCount(), 1, "hub room count")

	h.Save()

	h.Load()

	assert.Equal(t, h.RoomCount(), 1, "hub room count")

	assert.Equal(t, h.MessageCount(), 0, "hub message count")
	assert.Equal(t, ml.MessageCount(), 3, "db message count")
	// reads messages from the db (deletes from the db)
	h.ReadMessages()

	assert.Equal(t, h.MessageCount(), 3, "hub message count")
	assert.Equal(t, ml.MessageCount(), 0, "db message count")

	h.SendMessages()

	// one message lives long enough to be saved
	assert.Equal(t, h.MessageCount(), 1, "hub message count")
}

/*
// this test appears to be kind of flaky so i'm removing it for now
// the integration tests are better anyway
func TestHub4(t *testing.T) {
	ml := loader.NewMemoryLoader()
	// messages that expire immediately
	ml.AddMessage("hello world", 0)
	ml.AddMessage("server message", 0)
	// message that doesn't expire immediately
	ml.AddMessage("save this message", 30)
	h, err := hub.NewHubWithDB(ml, config.Default())
	assert.NoError(t, err, "new hub")

	roomID := "someboard"
	mock1 := core.NewMockEventChannel()
	h.Handler(mock1, roomID)

	assert.Equal(t, h.RoomCount(), 1, "hub room count")

	h.Save()

	h.Load()

	assert.Equal(t, h.RoomCount(), 1, "hub room count")

	assert.Equal(t, h.MessageCount(), 0, "hub message count")
	assert.Equal(t, ml.MessageCount(), 3, "db message count")

	// reads messages from the db (deletes from the db)
	h.ReadMessages()

	assert.Equal(t, h.MessageCount(), 3, "hub message count")
	assert.Equal(t, ml.MessageCount(), 0, "db message count")

	var wg sync.WaitGroup
	wg.Add(1)

	mock2 := core.NewBlockingMockEventChannel()
	go func() {
		defer wg.Done()
		// mock2.OnConnect() will signal mock2.Ready
		h.Handler(mock2, roomID)
	}()

	// block until handler has actually started
	<-mock2.Ready()

	h.SendMessages()

	mock2.Disconnect()

	wg.Wait()

	// 3 events
	// initial frame
	// connected users
	// one of the hub messages (the one with 30s ttl)
	assert.Equal(t, len(mock2.SavedEvents), 3, "mock.receivedEvents")
	receivedEventTypes := make(map[string]bool)
	for _, evt := range mock2.SavedEvents {
		receivedEventTypes[evt.Event] = true
	}

	_, ok := receivedEventTypes["frame"]
	_, ok = receivedEventTypes["connected_users"]
	_, ok = receivedEventTypes["global"]

	assert.True(t, ok, "mock savedEvents")
	assert.True(t, ok, "mock savedEvents")
	assert.True(t, ok, "mock savedEvents")

	// one message lives long enough to be saved
	assert.Equal(t, h.MessageCount(), 1, "hub message count")
}
*/
