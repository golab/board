/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/internal/fetch"
	"github.com/jarednogo/board/internal/sgfsamples"
	"github.com/jarednogo/board/pkg/config"
	"github.com/jarednogo/board/pkg/hub"
	"github.com/jarednogo/board/pkg/logx"
)

func TestApiRouter(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	_, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)
	r := chi.NewRouter()
	r.Mount("/api", hub.ApiRouter("version"))

	req := httptest.NewRequest("GET", "/api/version", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)

	msg := struct {
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(body, &msg)
	assert.NoError(t, err)

	assert.Equal(t, msg.Message, "version")
}

func TestExtRouter(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)

	room := h.GetOrCreateRoom("someboard")
	room.SetFetcher(fetch.NewMockFetcher(sgfsamples.SimpleEightMoves))

	r := chi.NewRouter()
	r.Mount("/ext", h.ExtRouter())

	v := url.Values{}
	v.Set("url", "https://online-go.com/game/1")
	v.Set("board_id", "someboard")
	path := "/ext/upload?" + v.Encode()

	req := httptest.NewRequest("GET", path, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, len(room.ToSGF()), 113)
}

func TestSocketRouter(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)

	r := chi.NewRouter()
	r.Mount("/socket", h.SocketRouter())
}
