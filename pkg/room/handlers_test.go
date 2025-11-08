/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room_test

import (
	"encoding/base64"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/internal/sgfsamples"
	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/event"
	"github.com/jarednogo/board/pkg/room"
)

func TestHandleIsProtected(t *testing.T) {
	r := room.NewRoom("")
	evt := event.NewEvent("isprotected", nil)
	evt = r.HandleAny(evt)
	assert.False(t, evt.Value().(bool), "isprotected")

	r.SetPassword("abcdef")
	evt = event.NewEvent("isprotected", nil)
	evt = r.HandleAny(evt)
	assert.True(t, evt.Value().(bool), "isprotected")
}

func TestHandleCheckPassword(t *testing.T) {
	r := room.NewRoom("")
	r.SetPassword(core.Hash("somepassword"))

	evt := event.NewEvent("checkpassword", "abcdef")
	evt = r.HandleAny(evt)
	assert.Equal(t, evt.Value().(string), "", "checkpassword")

	evt = event.NewEvent("checkpassword", "somepassword")
	evt = r.HandleAny(evt)
	assert.Equal(t, evt.Value().(string), "somepassword", "checkpassword")
}

func TestHandleUploadSGF(t *testing.T) {
	r := room.NewRoom("")
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.SimpleEightMoves))

	evt := event.NewEvent("upload_sgf", sgf)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF(false)), 113, "upload_sgf")
}

func TestHandleTrash(t *testing.T) {
	r := room.NewRoom("")
	sgf := base64.StdEncoding.EncodeToString([]byte(sgfsamples.SimpleEightMoves))

	evt := event.NewEvent("upload_sgf", sgf)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF(false)), 113, "trash")
	evt = event.NewEvent("trash", nil)
	r.HandleAny(evt)
	assert.Equal(t, len(r.ToSGF(false)), 65, "trash")
}

func TestHandleUpdateNickname(t *testing.T) {
	r := room.NewRoom("")
	evt := event.NewEvent("update_nickname", "mynick")
	evt.SetUser("user_123")
	r.HandleAny(evt)
	nicks := r.Nicks()
	n, ok := nicks["user_123"]
	assert.True(t, ok, "update_nickname")
	assert.Equal(t, n, "mynick", "update_nickname")
}

func TestHandleUpdateSettings(t *testing.T) {
	r := room.NewRoom("")
	value := make(map[string]any)
	value["buffer"] = 500.0
	value["size"] = 13.0
	value["nickname"] = "mynick"
	value["password"] = "somepassword"

	evt := event.NewEvent("update_settings", value)
	evt.SetUser("user_123")
	r.HandleAny(evt)

	assert.True(
		t,
		core.CorrectPassword("somepassword", r.GetPassword()),
		"update_settings",
	)
	assert.Equal(t, r.GetInputBuffer(), 500, "update_settings")
	assert.Equal(t, r.Size(), 13, "update_settings")

	nicks := r.Nicks()
	n, ok := nicks["user_123"]
	assert.True(t, ok, "update_settings")
	assert.Equal(t, n, "mynick", "update_settings")
}
