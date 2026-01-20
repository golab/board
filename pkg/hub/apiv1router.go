/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golab/board/pkg/event"
)

func (h *Hub) handler(w http.ResponseWriter, r *http.Request) {
	board := chi.URLParam(r, "board")
	room := h.GetOrCreateRoom(board)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf(`{"success": false, "error": "%s"}`, err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	evt, err := event.EventFromJSON(data)
	if err != nil {
		msg := fmt.Sprintf(`{"success": false, "error": "%s"}`, err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	evt = room.HandleAny(evt)
	if evt.Type() == "error" {
		v, ok := evt.Value().(string)
		if !ok {
			v = "unknown error occurred"
		}
		msg := fmt.Sprintf(`{"success": false, "error": "%s"}`, v)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	data, _ = json.Marshal(evt)
	fmt.Fprintf(w, `{"success": true, "output": %s}`, string(data)) // nolint:errcheck
}

func (h *Hub) ApiV1Router() http.Handler {
	r := chi.NewRouter()
	r.Post("/room/{board}", h.handler)
	return r
}
