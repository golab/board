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
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jarednogo/board/backend/core"
	"github.com/jarednogo/board/backend/server"
	"github.com/jarednogo/board/backend/twitch"
	"github.com/jarednogo/board/frontend"
)

// constants
const WSPORT = 9000
const WSHOST = "localhost"

func sanitize(s string) string {
	ok := []rune{}
	for _, c := range s {
		if (c >= '0' && c <= '9') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') {
			ok = append(ok, c)
		}
	}
	return string(ok)
}

func uuid4() string {
	r, _ := uuid.NewRandom()
	s := r.String()
	// remove hyphens
	return sanitize(s)
}

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
	templ.Execute(w, nil)
}

func includeCommon(w http.ResponseWriter, page string) {
	html := []string{"header.html", "menubar.html", "footer.html"}
	html = append([]string{page}, html...)
	templ := renderTemplate(html...)
	if templ == nil {
		return
	}
	templ.Execute(w, nil)
}

func twitchSubscribe(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	expiration := time.Now().Add(2 * time.Minute)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   true,
		Expires:  expiration,
		Path:     "/",
	})
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s/apps/twitch/callback&scope=%s&state=%s", twitch.ClientID(), core.MyURL(), "channel:bot", state)
	http.Redirect(w, r, url, http.StatusFound)
}

func twitchUnsubscribe(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	expiration := time.Now().Add(2 * time.Minute)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   true,
		Expires:  expiration,
		Path:     "/",
	})
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s/apps/twitch/callback&state=%s", twitch.ClientID(), core.MyURL(), state)
	http.Redirect(w, r, url, http.StatusFound)
}

func twitchCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")

	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != state {
		http.Error(w, "invalid state", http.StatusForbidden)
		return
	}

	if code != "" {

		// use the code to get an access token
		token, err := twitch.GetUserAccessToken(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		// use the user access token to get the user id
		user, err := twitch.GetUsers(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		// get an app access token (one could imagine putting this
		// in the subscribe function directly)
		token, err = twitch.GetAppAccessToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if scope == "" {
			// unsubscribe logic
			id, err := twitch.GetSubscription(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			// unsubscribe
			err = twitch.Unsubscribe(id, token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			log.Println("unsubscribing:", id, user)
		} else {
			// subscribe, get subscription id
			id, err := twitch.Subscribe(user, token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			log.Println("id of new subscription:", id, "for user:", user)
		}
	}

	//w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte(`{"message": "success"}`))
	w.Write([]byte("success"))
}

func about(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "about.html")
}

func index(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "index.html")
}

func twitchMain(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "twitch.html")
}

func board(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	if boardID != sanitize(boardID) {
		page400(w, r)
		return
	}
	renderSinglePage(w, "board.html")
}

func newBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.FormValue("board_id")
	boardID = sanitize(boardID)
	if len(strings.TrimSpace(boardID)) == 0 {
		boardID = uuid4()
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

func apiV1Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "pong"}`))
	})
	return r
}

func twitchRouter(s *server.Server) http.Handler {
	r := chi.NewRouter()
	r.Get("/subscribe", twitchSubscribe)
	r.Get("/unsubscribe", twitchUnsubscribe)
	r.Get("/callback", twitchCallback)
	r.Post("/callback", s.Twitch)
	return r
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
	r.Get("/board", board)
	r.Get("/twitch", twitchMain)
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

	r.Mount("/api/v1", apiV1Router())

	// see server package for routes
	r.Mount("/apps/twitch", twitchRouter(s))

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
