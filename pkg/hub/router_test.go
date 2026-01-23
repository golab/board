/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/fetch"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/internal/sgfsamples"
	"github.com/golab/board/pkg/config"
	"github.com/golab/board/pkg/event"
	"github.com/golab/board/pkg/hub"
	"github.com/golab/board/pkg/logx"
)

func TestApiRouter(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)
	r := chi.NewRouter()
	r.Mount("/api", h.ApiRouter("version123"))

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

	assert.Equal(t, msg.Message, "version123")
}

func TestApiRouterStats(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)

	// populate rooms
	r1 := h.GetOrCreateRoom("room1")
	r2 := h.GetOrCreateRoom("room2")
	mock1 := event.NewMockEventChannel()
	mock2 := event.NewMockEventChannel()
	mock3 := event.NewMockEventChannel()

	r1.RegisterConnection(mock1)
	r1.RegisterConnection(mock2)
	r2.RegisterConnection(mock3)

	// mount router
	r := chi.NewRouter()
	r.Mount("/api", h.ApiRouter("dev"))

	req := httptest.NewRequest("GET", "/api/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)

	msg := struct {
		Rooms       int `json:"rooms"`
		Connections int `json:"connections"`
	}{}

	err = json.Unmarshal(body, &msg)
	require.NoError(t, err)

	assert.Equal(t, msg.Rooms, 2)
	assert.Equal(t, msg.Connections, 3)
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

func TestWebRouterDebug(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)
	r := chi.NewRouter()
	r.Mount("/", h.WebRouter())

	// send a request to /debug (should be empty)
	req := httptest.NewRequest("GET", "/b/someboard/debug", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, len(body), 0)

	// now create the room
	h.GetOrCreateRoom("someboard")

	// resend the request (shouldn't be empty)
	req = httptest.NewRequest("GET", "/b/someboard/debug", nil)
	r.ServeHTTP(rec, req)

	body, err = io.ReadAll(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, len(body), 146)
}

func TestWebRouterSGF(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)
	r := chi.NewRouter()
	r.Mount("/", h.WebRouter())

	// send a request to /sgf (should be empty)
	req := httptest.NewRequest("GET", "/b/someboard/sgf", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, len(body), 0)

	// now create the room
	h.GetOrCreateRoom("someboard")

	// resend the request (shouldn't be empty)
	req = httptest.NewRequest("GET", "/b/someboard/sgf", nil)
	r.ServeHTTP(rec, req)

	body, err = io.ReadAll(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, len(body), 65)
}

func TestWebRouterSGFIX(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)
	r := chi.NewRouter()
	r.Mount("/", h.WebRouter())

	// send a request to /sgfix (should be empty)
	req := httptest.NewRequest("GET", "/b/someboard/sgfix", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, len(body), 0)

	// now create the room
	h.GetOrCreateRoom("someboard")

	// resend the request (shouldn't be empty)
	req = httptest.NewRequest("GET", "/b/someboard/sgfix", nil)
	r.ServeHTTP(rec, req)

	body, err = io.ReadAll(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, len(body), 70)
}

func getStatusCode(method, endpoint string) int {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	if err != nil {
		return 0
	}
	r := chi.NewRouter()
	r.Mount("/", h.WebRouter())

	// send a request to /sgfix (should be empty)
	req := httptest.NewRequest(method, endpoint, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec.Code
}

func TestWebRouterEndpoints(t *testing.T) {
	endpointTests := []struct {
		method   string
		endpoint string
		code     int
	}{
		{"GET", "/", 200},
		{"GET", "/integrations", 200},
		{"GET", "/about", 200},
		{"POST", "/new", 302},
		{"GET", "/b/someboard", 200},
		{"GET", "/static/js/init.js", 200},
		{"GET", "/static/foobar.js", 404},
		{"GET", "/favicon.ico", 200},
	}
	for i, tt := range endpointTests {
		t.Run(fmt.Sprintf("endpoint%d", i), func(t *testing.T) {
			code := getStatusCode(tt.method, tt.endpoint)
			assert.Equal(t, code, tt.code)
		})
	}
}

func TestApiV1Router1(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)
	r := chi.NewRouter()
	r.Mount("/api/v1", h.ApiV1Router())

	payload := []byte(`{"event": "add_stone", "value": {"coords": [0, 0], "color": 1}}`)
	req := httptest.NewRequest("POST", "/api/v1/room/abc", bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	output := struct {
		Success bool               `json:"success"`
		Event   event.DefaultEvent `json:"output"`
	}{}

	err = json.Unmarshal(body, &output)
	require.NoError(t, err)
	evt := output.Event

	assert.Equal(t, output.Success, true)
	assert.Equal(t, evt.Type(), "frame")
}

func TestApiV1Router2(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)
	r := chi.NewRouter()
	r.Mount("/api/v1", h.ApiV1Router())

	for i := 0; i < 2; i++ {
		payload := []byte(`{"event": "add_stone", "value": {"coords": [0, 0], "color": 1}}`)
		req := httptest.NewRequest("POST", "/api/v1/room/abc", bytes.NewBuffer(payload))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		body, err := io.ReadAll(rec.Body)
		require.NoError(t, err)

		output := struct {
			Success bool               `json:"success"`
			Event   event.DefaultEvent `json:"output"`
		}{}

		err = json.Unmarshal(body, &output)
		require.NoError(t, err)
		evt := output.Event
		assert.Equal(t, output.Success, true)

		// trying to add a stone on a coordinate where there's already a stone
		// is illegal, so instead of getting a "frame" event back, we just
		// get back the original event
		if i == 0 {
			assert.Equal(t, evt.Type(), "frame")
		} else {
			assert.Equal(t, evt.Type(), "add_stone")
		}
	}
}

func TestApiV1Router3(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)
	r := chi.NewRouter()
	r.Mount("/api/v1", h.ApiV1Router())

	payload := []byte(`{"invalidjson`)
	req := httptest.NewRequest("POST", "/api/v1/room/abc", bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	output := struct {
		Success bool               `json:"success"`
		Event   event.DefaultEvent `json:"output"`
	}{}

	err = json.Unmarshal(body, &output)
	require.NoError(t, err)

	assert.Equal(t, output.Success, false)
}

func TestApiV1Router4(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := hub.NewHub(config.Test(), logger)
	assert.NoError(t, err)
	r := chi.NewRouter()
	r.Mount("/api/v1", h.ApiV1Router())

	payload := []byte(`{"event": "upload_sgf", "value": ")invalidsgf"}`)
	req := httptest.NewRequest("POST", "/api/v1/room/abc", bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	output := struct {
		Success bool               `json:"success"`
		Event   event.DefaultEvent `json:"output"`
	}{}

	err = json.Unmarshal(body, &output)
	require.NoError(t, err)

	assert.Equal(t, output.Success, false)
}
