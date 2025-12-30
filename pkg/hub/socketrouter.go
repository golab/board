/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golab/board/pkg/config"
	"golang.org/x/net/websocket"
)

type socketServer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type MockSocketServer struct {
}

func (ms *MockSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func NewMockSocketServer() *MockSocketServer {
	return &MockSocketServer{}
}

func DefaultSocketServer(handler websocket.Handler) socketServer {
	cfg := websocket.Config{}
	ws := websocket.Server{
		Config:    cfg,
		Handshake: nil,
		Handler:   handler,
	}
	return ws
}

func (h *Hub) SocketRouter() http.Handler {
	r := chi.NewRouter()

	// create new websocket server
	var ws socketServer
	if h.cfg.Mode == config.ModeTest {
		ws = NewMockSocketServer()
	} else {
		ws = DefaultSocketServer(h.HandlerWrapper)
	}

	r.Get("/b/{boardID}", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeHTTP(w, r)
	})

	return r
}
