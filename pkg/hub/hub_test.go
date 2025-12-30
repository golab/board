/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/internal/sgfsamples"
	"github.com/golab/board/pkg/config"
	"github.com/golab/board/pkg/event"
	"github.com/golab/board/pkg/hub"
	"github.com/golab/board/pkg/loader"
	"github.com/golab/board/pkg/logx"
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
			assert.Equal(t, [3]string{a, b, c}, tc.output)
		})
	}
}

func TestHub1(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHubWithDB(loader.NewMemoryLoader(), config.Default(), logger)
	assert.NoError(t, err)
	h.Load()

	mock := event.NewMockEventChannel()
	roomID := "someboard"
	h.Handler(mock, roomID)

	h.Save()
}

func TestHub2(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHubWithDB(loader.NewMemoryLoader(), config.Default(), logger)
	assert.NoError(t, err)

	mock := event.NewMockEventChannel()
	evt := event.NewEvent("pass", 1.0)
	mock.QueuedEvents = append(mock.QueuedEvents, evt)
	roomID := "someboard"
	h.Handler(mock, roomID)

	// 3 events:
	// initial frame
	// connected users
	// pass event
	assert.Equal(t, len(mock.SavedEvents), 3)
}

func TestHub3(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	ml := loader.NewMemoryLoader()
	// messages that expire immediately
	ml.AddMessage("hello world", 0)
	ml.AddMessage("server message", 0)
	// message that doesn't expire immediately
	ml.AddMessage("save this message", 30)
	h, err := hub.NewHubWithDB(ml, config.Default(), logger)
	assert.NoError(t, err)

	roomID := "someboard"
	mock1 := event.NewMockEventChannel()
	h.Handler(mock1, roomID)

	assert.Equal(t, h.RoomCount(), 1)

	h.Save()

	h.Load()

	assert.Equal(t, h.RoomCount(), 1)

	assert.Equal(t, h.MessageCount(), 0)
	assert.Equal(t, ml.MessageCount(), 3)
	// reads messages from the db (deletes from the db)
	h.ReadMessages()

	assert.Equal(t, h.MessageCount(), 3)
	assert.Equal(t, ml.MessageCount(), 0)

	h.SendMessages()

	// one message lives long enough to be saved
	assert.Equal(t, h.MessageCount(), 1)
}

func TestGetRoom(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHubWithDB(loader.NewMemoryLoader(), config.Default(), logger)
	assert.NoError(t, err)
	roomID := "room123"
	_, err = h.GetRoom(roomID)
	assert.NotNil(t, err)
	_ = h.GetOrCreateRoom(roomID)
	room, err := h.GetRoom(roomID)
	assert.NoError(t, err)
	assert.Equal(t, room.ID(), roomID)
}

func TestRoomLogger(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHubWithDB(loader.NewMemoryLoader(), config.Default(), logger)
	assert.NoError(t, err)
	roomID := "room123"
	h.GetOrCreateRoom(roomID)

	mock := event.NewMockEventChannel()
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.SimpleEightMoves))
	evt := event.NewEvent("upload_sgf", sgf)
	mock.QueuedEvents = append(mock.QueuedEvents, evt)
	h.Handler(mock, roomID)

	require.Equal(t, len(logger.Lines()), 5)
	upload := logger.Lines()[2]
	log := struct {
		RoomID string `json:"room_id"`
	}{}
	err = json.Unmarshal([]byte(upload), &log)
	assert.NoError(t, err)
	assert.Equal(t, log.RoomID, roomID)
}

func TestHubConnCount(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHubWithDB(loader.NewMemoryLoader(), config.Default(), logger)
	assert.NoError(t, err)

	assert.Equal(t, h.ConnCount(), 0)

	r1 := h.GetOrCreateRoom("room1")
	r2 := h.GetOrCreateRoom("room2")
	mock1 := event.NewMockEventChannel()
	mock2 := event.NewMockEventChannel()
	mock3 := event.NewMockEventChannel()

	r1.RegisterConnection(mock1)
	r1.RegisterConnection(mock2)
	r2.RegisterConnection(mock3)

	assert.Equal(t, h.ConnCount(), 3)
}
