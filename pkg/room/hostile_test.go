/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room_test

import (
	"fmt"
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/sgfsamples"
	"github.com/golab/board/pkg/event"
	"github.com/golab/board/pkg/logx"
	"github.com/golab/board/pkg/room"
)

// every event type is fed values of the wrong shape: none should panic
func TestHostileEventsDontPanic(t *testing.T) {
	events := []string{
		"isprotected", "checkpassword", "upload_sgf", "request_sgf", "trash", "update_nickname",
		"update_settings", "graft", "add_stone", "pass", "remove_stone", "triangle", "square",
		"letter", "number", "label", "remove_mark", "cut", "left", "right", "up", "down", "rewind",
		"fastforward", "goto_grid", "goto_coord", "comment", "draw", "erase_pen", "copy",
		"clipboard", "score", "markdead",
	}
	values := []any{
		nil, "x", 5.0, 1e18, -3.0,
		[]any{"a", "b"}, []any{30.0, 30.0}, []any{-1.0, 5.0}, []any{1.0, "b"},
		[]any{1.0, 2.0, 3.0}, []any{nil, nil},
		map[string]any{},
		map[string]any{"coords": []any{30.0, 30.0}, "color": 1.0, "letter": "A", "label": "x", "number": 1.0},
		map[string]any{"coords": []any{"a", "b"}, "color": 1.0},
		map[string]any{"coords": []any{3.0, 3.0}, "color": "x"},
		map[string]any{"buffer": "x", "size": 19.0},
		map[string]any{"buffer": 1.0, "size": 1e9, "nickname": "", "password": "", "black": "", "white": "", "komi": ""},
		"z99 a1", "t19",
	}
	for _, e := range events {
		for i, v := range values {
			t.Run(fmt.Sprintf("%s/%d", e, i), func(t *testing.T) {
				r := room.NewRoom("hostile")
				r.SetLogger(logx.NewRecorder(logx.LogLevelDebug))
				r.DisableBuffers()
				r.UploadSGF(sgfsamples.SimpleEightMoves)
				defer func() {
					if rec := recover(); rec != nil {
						t.Fatalf("panic: %v", rec)
					}
				}()
				r.HandleAny(event.NewEvent(e, v))
			})
		}
	}
}

func TestUpdateSettingsValidation(t *testing.T) {
	valid := func() map[string]any {
		return map[string]any{
			"buffer": 250.0, "size": 19.0, "nickname": "", "password": "",
			"black": "", "white": "", "komi": "",
		}
	}

	bad := valid()
	bad["size"] = 20.0
	missing := valid()
	delete(missing, "komi")

	for _, v := range []map[string]any{bad, missing} {
		r := room.NewRoom("")
		r.SetInputBuffer(500)
		evt := r.HandleAny(event.NewEvent("update_settings", v))
		assert.Equal(t, evt.Type(), "error")
		// nothing was applied
		assert.Equal(t, r.Size(), 19)
		assert.Equal(t, r.GetInputBuffer(), 500)
	}

	good := valid()
	good["size"] = 9.0
	r := room.NewRoom("")
	r.HandleAny(event.NewEvent("update_settings", good))
	assert.Equal(t, r.Size(), 9)
}
