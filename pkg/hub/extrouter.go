/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golab/board/pkg/core"
	"github.com/golab/board/pkg/event"
)

func (h *Hub) Upload(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	boardID := r.FormValue("board_id")
	boardID = core.Sanitize(boardID)
	if len(strings.TrimSpace(boardID)) == 0 {
		boardID = core.UUID4()
	}
	newroom := h.GetOrCreateRoom(boardID)

	var evt event.Event
	if url != "" {
		evt = event.NewEvent("request_sgf", url)
	} else {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	newroom.HandleAny(evt)

	redirect := fmt.Sprintf("/b/%s", boardID)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (h *Hub) ExtRouter() http.Handler {
	r := chi.NewRouter()

	// stateful endpoints
	r.Get("/upload", h.Upload)

	return r
}
