package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"strings"
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
	paths := []string{}
	for _, fname := range fnames {
		path := fmt.Sprintf("%s/%s", "html", fname)
		paths = append(paths, path)
	}
	templ, err := template.ParseFiles(paths...)
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

func about(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "about.html")
}

func index(w http.ResponseWriter, r *http.Request) {
	includeCommon(w, "index.html")
}

func board(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	if boardID != sanitize(boardID) {
		page400(w, r)
		return
	}
	renderSinglePage(w, "board.html")
}

func suffixOp(w http.ResponseWriter, r *http.Request, suffix string) {
	boardID := chi.URLParam(r, "boardID")
	wsURL := fmt.Sprintf("ws://%s:%d/b/%s/%s", WSHOST, WSPORT, boardID, suffix)
	ws, err := websocket.Dial(wsURL, "", "http://localhost")
	if err != nil {
		return
	}

	dataLen := make([]byte, 4)
	ws.Read(dataLen)
	length := binary.LittleEndian.Uint32(dataLen)

	data, err := readBytes(ws, int(length))
	//data := make([]byte, length)
	//ws.Read(data)

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return
	}
	w.Write(decoded)
}

func sgf(w http.ResponseWriter, r *http.Request) {
	suffixOp(w, r, "sgf")
}

func sgfix(w http.ResponseWriter, r *http.Request) {
	suffixOp(w, r, "sgfix")
}

func debug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	suffixOp(w, r, "debug")
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

func upload(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	sgf := r.URL.Query().Get("sgf")
	boardID := r.URL.Query().Get("board_id")
	boardID = sanitize(boardID)
	if len(strings.TrimSpace(boardID)) == 0 {
		boardID = uuid4()
	}
	if url != "" {
		requestSGF(boardID, url)
	} else if sgf != "" {
		uploadSGF(boardID, sgf)
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
	path := fmt.Sprintf("%s/%s", "static", fname)
	// TODO: make sure this is locked down
	// TODO: make sure this uses the 404 page when the file isn't found
	http.ServeFile(w, r, path)
}

func favicon(w http.ResponseWriter, r *http.Request) {
	serveStatic(w, r, "favicon.svg")
}

func image(w http.ResponseWriter, r *http.Request) {
	image := chi.URLParam(r, "image")
	serveStatic(w, r, image)
}

// socket stuff

type EventJSON struct {
	Event string `json:"event"`
	Value string `json:"value"`
}

func websocketSend(e *EventJSON, boardID string) {
	route := fmt.Sprintf("/b/%s", boardID)
	wsURL := fmt.Sprintf("ws://%s:%d%s", WSHOST, WSPORT, route)

	payload, err := json.Marshal(e)
	if err != nil {
		return
	}
	length := uint32(len(payload))
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, length)
	// send to websocket
	ws, err := websocket.Dial(wsURL, "", "http://localhost")
	if err != nil {
		return
	}
	ws.Write(buf)
	ws.Write(payload)
	ws.Close()
}

func readBytes(ws *websocket.Conn, size int) ([]byte, error) {
	chunkSize := 64
	message := []byte{}
	for {
		if len(message) >= size {
			break
		}
		l := size - len(message)
		if l > chunkSize {
			l = chunkSize
		}
		temp := make([]byte, l)
		n, err := ws.Read(temp)
		if err != nil {
			return nil, err
		}
		message = append(message, temp[:n]...)
	}
	return message, nil
}

func requestSGF(boardID, url string) {
	e := &EventJSON{
		Event: "request_sgf",
		Value: url}
	websocketSend(e, boardID)
}

func uploadSGF(boardID, sgf string) {
	e := &EventJSON{
		Event: "upload_sgf",
		Value: sgf,
	}
	websocketSend(e, boardID)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)

	r.Get("/", index)
	r.Get("/about", about)
	r.Get("/board", board)
	r.Get("/favicon.ico", favicon)
	r.Post("/new", newBoard)
	r.Get("/upload", upload)

	r.Get("/static/{image}", image)

	r.Get("/b/{boardID}", board)
	r.Get("/b/{boardID}/sgf", sgf)
	r.Get("/b/{boardID}/sgfix", sgfix)
	r.Get("/b/{boardID}/debug", debug)

	r.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	r.NotFound(page404)

	port := "8080"
	host := "localhost"
	url := fmt.Sprintf("%s:%s", host, port)
	log.Println("Listening on", url)
	http.ListenAndServe(url, r)
}
