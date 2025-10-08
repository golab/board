/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"io/fs"

	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jarednogo/board/backend/core"
	"github.com/jarednogo/board/backend/server"
	"github.com/jarednogo/board/frontend"
)

// constants
const WSPORT = 9000
const WSHOST = "localhost"

// html

func renderTemplate(fnames ...string) *template.Template {
	htmlFS, _ := fs.Sub(frontend.Content, "html")

	templ, err := template.ParseFS(htmlFS, fnames...)
	if err != nil {
		return nil
	}
	return templ
}

func renderSinglePage(w http.ResponseWriter, page string) {
	templ := renderTemplate(page)
	if templ == nil {
		return
	}
	err := templ.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func includeCommon(w http.ResponseWriter, page string) {
	html := []string{"header.html", "menubar.html", "footer.html"}
	html = append([]string{page}, html...)
	templ := renderTemplate(html...)
	if templ == nil {
		return
	}
	err := templ.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "about.html")
}

func index(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "index.html")
}

func integrations(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "integrations.html")
}

func board(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	if boardID != core.Sanitize(boardID) {
		page400(w, r)
		return
	}
	renderSinglePage(w, "board.html")
}

func newBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.FormValue("board_id")
	boardID = core.Sanitize(boardID)
	if len(strings.TrimSpace(boardID)) == 0 {
		boardID = core.UUID4()
	}
	redirect := fmt.Sprintf("/b/%s", boardID)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func page400(w http.ResponseWriter, r *http.Request) {
	log.Printf("400 - %s\n", r.URL)
	includeCommon(w, "400.html")
}

func page404(w http.ResponseWriter, r *http.Request) {
	log.Printf("404 - %s\n", r.URL)
	includeCommon(w, "404.html")
}

//static

func serveStatic(w http.ResponseWriter, r *http.Request, fname string) {
	staticFS, _ := fs.Sub(frontend.Content, "static")

	//path := fmt.Sprintf("%s/%s", "static", fname)
	// TODO: make sure this is locked down
	// TODO: make sure this uses the 404 page when the file isn't found
	http.ServeFileFS(w, r, staticFS, fname)
}

func favicon(w http.ResponseWriter, r *http.Request) {
	serveStatic(w, r, "favicon.svg")
}

func image(w http.ResponseWriter, r *http.Request) {
	image := chi.URLParam(r, "image")
	serveStatic(w, r, image)
}

func main() {
	// websocket server setup
	cfg := websocket.Config{}
	s := server.NewServer()
	s.Load()
	defer s.Save()

	// create new websocket server
	ws := websocket.Server{
		Config:    cfg,
		Handshake: nil,
		Handler:   s.Handler,
	}

	// http server setup

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)

	r.Get("/", index)
	r.Get("/about", about)
	r.Get("/integrations", integrations)
	r.Get("/favicon.ico", favicon)
	r.Post("/new", newBoard)

	r.Get("/upload", s.Upload)

	r.Get("/static/{image}", image)

	r.Get("/b/{boardID}", board)
	r.Get("/b/{boardID}/sgf", s.Sgf)
	r.Get("/b/{boardID}/sgfix", s.Sgfix)
	r.Get("/b/{boardID}/debug", s.Debug)

	jsFS, _ := fs.Sub(frontend.Content, "js")
	r.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.FS(jsFS))))

	r.NotFound(page404)

	r.Mount("/api/v1", server.ApiV1Router())

	// see server package for routes
	r.Mount("/apps/twitch", s.TwitchRouter())

	// mount websocket
	r.Get("/socket/b/{boardID}", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeHTTP(w, r)
	})

	// start everything
	port := "8080"
	host := "localhost"
	url := fmt.Sprintf("%s:%s", host, port)

	// get ready to catch signals
	cancelChan := make(chan os.Signal, 1)

	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	// run server
	go func() {
		log.Println("Listening on", url)
		log.Fatal(http.ListenAndServe(url, r))
	}()

	// catch cancel signal
	sig := <-cancelChan

	log.Printf("Caught signal %v", sig)
	log.Println("Shutting down gracefully")

}
