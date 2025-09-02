/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package server

import (
	"fmt"
	"log"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/jarednogo/board/backend/core"
)

func (s *Server) ApiV1Router() http.Handler {
	r := chi.NewRouter()
	r.Post("/twitch", s.Twitch)
	return r
}

func (s *Server) Twitch(w http.ResponseWriter, r *http.Request) {
	roomID := "testroom"
	room := s.GetOrCreateRoom(roomID)

	data, _ := io.ReadAll(r.Body)

	evt := &core.EventJSON {
		Event: "graft",
		Value: string(data),
	}

	handler := Chain(
		room.HandleEvent,
		room.BroadcastFullFrameAfter)
	handler(evt)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) Debug(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	data := s.HandleOp("debug", boardID)
	w.Write([]byte(data))
}

func (s *Server) Sgfix(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	data := s.HandleOp("sgfix", boardID)
	w.Write([]byte(data))
}

func (s *Server) Sgf(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	data := s.HandleOp("sgf", boardID)
	w.Write([]byte(data))
}

func (s *Server) Upload(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	sgf := r.URL.Query().Get("sgf")
	log.Println("url:", url)
	log.Println("sgf:", sgf)
	boardID := r.URL.Query().Get("board_id")
	boardID = sanitize(boardID)
	if len(strings.TrimSpace(boardID)) == 0 {
		boardID = uuid4()
	}
	room := s.GetOrCreateRoom(boardID)

	var handler EventHandler
	var evt *core.EventJSON
	if url != "" {
		handler = Chain(
			room.HandleRequestSGF,
			room.OutsideBuffer,
			room.Authorized,
			room.CloseOGS,
			room.BroadcastAfter(false))
		evt = &core.EventJSON {
			Event: "request_sgf",
			Value: url,
		}
	} else if sgf != "" {
		handler = Chain(
			room.HandleUploadSGF,
			room.OutsideBuffer,
			room.Authorized,
			room.CloseOGS,
			room.BroadcastAfter(false))
		evt = &core.EventJSON {
			Event: "upload_sgf",
			Value: sgf,
		}
	} else {
		return
	}
	handler(evt)

	redirect := fmt.Sprintf("/b/%s", boardID)
	http.Redirect(w, r, redirect, http.StatusFound)
}

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


