/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golab/board/pkg/core"
	"github.com/golab/board/pkg/frontend"
)

// stateful

func (h *Hub) HandleOp(op, roomID string) string {
	roomID = strings.ToLower(roomID)
	data := ""
	r, err := h.GetRoom(roomID)
	if err != nil {
		return ""
	}
	switch op {
	case "sgf":
		// if the room doesn't exist, send empty string
		data = r.ToSGF()
	case "sgfix":
		// basically do the same thing but include indexes
		data = r.ToSGFIX()
	case "debug":
		// send debug info
		stateJSON := r.SaveState()
		dataBytes, _ := json.Marshal(stateJSON)
		data = string(dataBytes)
	}
	return data
}

func (h *Hub) Debug(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	data := h.HandleOp("debug", boardID)
	_, err := w.Write([]byte(data))
	if err != nil {
		h.logger.Error("/debug", "err", err)
	}
}

func (h *Hub) Sgfix(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	data := h.HandleOp("sgfix", boardID)
	_, err := w.Write([]byte(data))
	if err != nil {
		h.logger.Error("/sgfix", "err", err)
	}
}

func (h *Hub) Sgf(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	data := h.HandleOp("sgf", boardID)
	_, err := w.Write([]byte(data))
	if err != nil {
		h.logger.Error("/sgf", "err", err)
	}
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
	templ.Execute(w, nil) //nolint:errcheck
}

func includeCommon(w http.ResponseWriter, page string) {
	html := []string{"header.html", "menubar.html", "footer.html"}
	html = append([]string{page}, html...)
	templ := renderTemplate(html...)
	if templ == nil {
		return
	}
	templ.Execute(w, nil) //nolint:errcheck
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

func integrationsDisabled(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "integrations_disabled.html")
}

func board(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	boardID = strings.ToLower(boardID)
	if boardID != core.Sanitize(boardID) {
		page400(w, r)
		return
	}
	renderSinglePage(w, "board.html") //nolint:errcheck
}

func newBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.FormValue("board_id")
	boardID = core.Sanitize(boardID)
	if len(strings.TrimSpace(boardID)) == 0 {
		//boardID = core.UUID4()
		boardID = core.RandomBoardName()
	}
	redirect := fmt.Sprintf("/b/%s", boardID)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func page400(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "400.html")
}

func page404(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "404.html")
}

func page500(w http.ResponseWriter, r *http.Request) { //nolint:unused
	includeCommon(w, "500.html")
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

func static(w http.ResponseWriter, r *http.Request) {
	static := chi.URLParam(r, "static")
	serveStatic(w, r, static)
}

func (h *Hub) WebRouter() http.Handler {
	r := chi.NewRouter()

	// stateful endpoints
	r.Get("/b/{boardID}/sgf", h.Sgf)
	r.Get("/b/{boardID}/sgfix", h.Sgfix)
	r.Get("/b/{boardID}/debug", h.Debug)

	// pure web endpoints
	r.Get("/", index)
	r.Get("/about", about)
	if h.cfg.TwitchEnabled() {
		r.Get("/integrations", integrations)
	} else {
		r.Get("/integrations", integrationsDisabled)
	}
	r.Get("/favicon.ico", favicon)
	r.Post("/new", newBoard)

	r.Get("/static/{static}", static)

	r.Get("/b/{boardID}", board)
	jsFS, _ := fs.Sub(frontend.Content, "js")
	r.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.FS(jsFS))))

	r.NotFound(page404)

	return r
}
