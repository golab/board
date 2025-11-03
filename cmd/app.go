/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jarednogo/board/pkg/config"
	"github.com/jarednogo/board/pkg/hub"
)

var version = "dev"

func Setup(cfg *config.Config) (http.Handler, error) {
	// make a new hub
	h, err := hub.NewHub(cfg)
	if err != nil {
		return nil, err
	}
	h.Load()

	// http server setup

	// initialize router and middlewares
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)

	// web router
	r.Mount("/", h.WebRouter())

	// extension router
	r.Mount("/ext", h.ExtRouter())

	// api routers
	r.Mount("/api", hub.ApiRouter(version))
	r.Mount("/api/v1", hub.ApiV1Router())

	// see server package for routes
	r.Mount("/apps/twitch", h.TwitchRouter())

	// mount websocket
	r.Mount("/socket", h.SocketRouter())

	return r, nil
}
