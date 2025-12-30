/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/fetch"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/internal/sgfsamples"
	"github.com/golab/board/pkg/core"
	"github.com/golab/board/pkg/event"
	"github.com/golab/board/pkg/logx"
	"github.com/golab/board/pkg/room"
	"github.com/golab/board/pkg/room/plugin"
)

func TestHandleIsProtected(t *testing.T) {
	r := room.NewRoom("")
	evt := event.NewEvent("isprotected", nil)
	evt = r.HandleAny(evt)
	assert.False(t, evt.Value().(bool))

	r.SetPassword("abcdef")
	evt = event.NewEvent("isprotected", nil)
	evt = r.HandleAny(evt)
	assert.True(t, evt.Value().(bool))
}

func TestHandleCheckPassword(t *testing.T) {
	r := room.NewRoom("")
	r.SetPassword(core.Hash("somepassword"))

	evt := event.NewEvent("checkpassword", "abcdef")
	evt = r.HandleAny(evt)
	assert.Equal(t, evt.Value().(string), "")

	evt = event.NewEvent("checkpassword", "somepassword")
	evt = r.HandleAny(evt)
	assert.Equal(t, evt.Value().(string), "somepassword")
}

func TestHandleUploadSGF1(t *testing.T) {
	r := room.NewRoom("")
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.SimpleEightMoves))

	evt := event.NewEvent("upload_sgf", sgf)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 113)
}

func TestHandleUploadSGF2(t *testing.T) {
	r := room.NewRoom("")
	sgf := base64.StdEncoding.EncodeToString(sgfsamples.SimpleZip)

	evt := event.NewEvent("upload_sgf", sgf)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 105)
}

func TestHandleUploadSGF3(t *testing.T) {
	r := room.NewRoom("")
	sgf1 := base64.StdEncoding.EncodeToString([]byte(sgfsamples.SimpleFourMoves))
	sgf2 := base64.StdEncoding.EncodeToString([]byte(sgfsamples.SimpleEightMoves))
	sgfs := []any{sgf1, sgf2}

	evt := event.NewEvent("upload_sgf", sgfs)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 213)
}

func TestHandleUploadSGF4(t *testing.T) {
	r := room.NewRoom("")
	data := make([]byte, 1<<20)
	data = append([]byte("(;GM[1]SZ[19];B[aa])"), data...)
	sgf := base64.StdEncoding.EncodeToString(data)

	evt := event.NewEvent("upload_sgf", sgf)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 65)
}

func TestHandleUploadSGF5(t *testing.T) {
	r := room.NewRoom("")
	data := make([]byte, 1<<19)
	data = append([]byte("(;GM[1]SZ[19];B[aa])"), data...)
	sgf := base64.StdEncoding.EncodeToString(data)

	evt := event.NewEvent("upload_sgf", sgf)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 20)
}

func TestHandleTrash(t *testing.T) {
	r := room.NewRoom("")
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.SimpleEightMoves))

	evt := event.NewEvent("upload_sgf", sgf)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 113)
	evt = event.NewEvent("trash", nil)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 65)
}

func TestHandleUpdateNickname(t *testing.T) {
	r := room.NewRoom("")
	evt := event.NewEvent("update_nickname", "mynick")
	evt.SetUser("user_123")
	r.HandleAny(evt)
	n, ok := r.GetNick("user_123")
	assert.True(t, ok)
	assert.Equal(t, n, "mynick")
}

func TestHandleUpdateSettings(t *testing.T) {
	r := room.NewRoom("")
	value := make(map[string]any)
	value["buffer"] = 500.0
	value["size"] = 13.0
	value["nickname"] = "mynick"
	value["password"] = "somepassword"
	value["black"] = "black123"
	value["white"] = "white456"
	value["komi"] = "10.5"

	evt := event.NewEvent("update_settings", value)
	evt.SetUser("user_123")
	r.HandleAny(evt)

	assert.True(
		t,
		core.CorrectPassword("somepassword", r.GetPassword()),
	)
	assert.Equal(t, r.GetInputBuffer(), 500)
	assert.Equal(t, r.Size(), 13)

	n, ok := r.GetNick("user_123")
	assert.True(t, ok)
	assert.Equal(t, n, "mynick")
}

func TestHandleRequestSGF1(t *testing.T) {
	r := room.NewRoom("")
	r.SetFetcher(fetch.NewMockFetcher(sgfsamples.Empty))

	evt := event.NewEvent("request_sgf", "http://www.gokifu.com/somefile.sgf")
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 65)
}

func TestHandleRequestSGF2(t *testing.T) {
	r := room.NewRoom("")
	r.SetFetcher(fetch.NewMockFetcher(sgfsamples.Empty))

	evt := event.NewEvent("request_sgf", "https://online-go.com/game/1")
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 65)
}

func TestHandleAddStone(t *testing.T) {
	r := room.NewRoom("")

	val := make(map[string]any)
	val["coords"] = []any{9.0, 9.0}
	val["color"] = 1.0
	evt := event.NewEvent("add_stone", val)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 71)
}

func TestSlow(t *testing.T) {
	r := room.NewRoom("")

	val := make(map[string]any)
	val["coords"] = []any{9.0, 9.0}
	val["color"] = 1.0
	evt := event.NewEvent("add_stone", val)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 71)

	// adding a stone with the buffers in place shouldn't do anything
	val["coords"] = []any{10.0, 10.0}
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 71)

	// disable the buffers then add a stone, and there should be a change
	r.DisableBuffers()
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF()), 77)
}

func TestLogUploadSGF(t *testing.T) {
	l := logx.NewRecorder(logx.LogLevelInfo)
	r := room.NewRoom("")
	r.SetLogger(l)

	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.SimpleEightMoves))
	evt := event.NewEvent("upload_sgf", sgf)
	r.HandleAny(evt)

	require.Equal(t, len(l.Lines()), 2)
	log := struct {
		EventType string `json:"event_type"`
	}{}
	err := json.Unmarshal([]byte(l.Lines()[0]), &log)
	assert.NoError(t, err)
	assert.Equal(t, log.EventType, "upload_sgf")
}

func TestLogRequestSGF(t *testing.T) {
	l := logx.NewRecorder(logx.LogLevelInfo)
	r := room.NewRoom("")
	r.SetLogger(l)
	r.SetFetcher(fetch.NewMockFetcher(sgfsamples.Empty))

	evt := event.NewEvent("request_sgf", "http://www.gokifu.com/somefile.sgf")
	r.HandleAny(evt)
	require.Equal(t, len(l.Lines()), 3)
	log := struct {
		EventType string `json:"event_type"`
	}{}
	err := json.Unmarshal([]byte(l.Lines()[0]), &log)
	assert.NoError(t, err)
	assert.Equal(t, log.EventType, "request_sgf")
}

func TestLogRequestOGS(t *testing.T) {
	l := logx.NewRecorder(logx.LogLevelInfo)
	r := room.NewRoom("")
	r.SetLogger(l)
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r.SetFetcher(f)

	r.SetConnector("ogs", plugin.NewMockOGSPlugin)

	evt := event.NewEvent("request_sgf", "http://online-go.com/review/42")
	r.HandleAny(evt)

	p := r.GetPlugin("ogs")
	require.NotNil(t, p)
}
