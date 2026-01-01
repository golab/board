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
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Hub) stats(w http.ResponseWriter, r *http.Request) {
	s := struct {
		Rooms       int `json:"rooms"`
		Connections int `json:"connections"`
	}{
		Rooms:       h.RoomCount(),
		Connections: h.ConnCount(),
	}
	data, err := json.Marshal(s)
	if err != nil {
		return
	}
	w.Write(data) //nolint:errcheck
}

func (h *Hub) ApiRouter(version string) http.Handler {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "pong"}`)) //nolint:errcheck
	})
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf(`{"message": "%s"}`, version)
		w.Write([]byte(msg)) //nolint:errcheck
	})
	r.Get("/stats", h.stats)
	return r
}
