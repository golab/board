/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/jarednogo/board/backend/core"
	"github.com/jarednogo/board/backend/loader"
	"github.com/jarednogo/board/backend/twitch"
)

func (s *Server) ApiV1Router() http.Handler {
	r := chi.NewRouter()
	r.Post("/twitch", s.Twitch)
	return r
}

func (s *Server) Twitch(w http.ResponseWriter, r *http.Request) {
	// read the body into a []byte
	body, _ := io.ReadAll(r.Body)

	// try to read the body into a TwitchJSON struct
	var req twitch.TwitchJSON
	err := json.Unmarshal(body, &req)
	if err != nil {
		log.Println(err)
		return
	}

	// on subscriptions, twitch sends a challenge that we need to respond to
	if req.Challenge != "" {
		w.Write([]byte(req.Challenge))
		return
	}

	// Grab headers for verification
	msgid := r.Header.Get("Twitch-Eventsub-Message-Id")
	timestamp := r.Header.Get("Twitch-Eventsub-Message-Timestamp")
	signature := r.Header.Get("Twitch-Eventsub-Message-Signature")

	// concat for verification
	message := msgid + timestamp + string(body)

	// do verification
	v := twitch.Verify(message, signature)
	if !v {
		log.Println("unverified message")
		return
	}

	// try to pull out the event
	evt := req.Event
	if evt == nil {
		log.Println("no event parsed")
		return
	}

	// try to pull out the message
	if evt.Message == nil {
		log.Println("no message parsed")
		return
	}

	// get broadcaster and chatter
	broadcaster := evt.BroadcasterUserID
	chatter := evt.ChatterUserID

	// extract the message in chat
	text := evt.Message.Text
	chat, err := twitch.Parse(text)
	if err != nil {
		//log.Println(err)
		return
	}

	// only care about the relevant commands
	if chat.Command != "branch" && chat.Command != "setboard" {
		log.Println("invalid command:", chat.Command)
		return
	}

	log.Println(chat.Command, chat.Body)

	// make sure only the broadcaster can set the room
	if chat.Command == "setboard" {
		if broadcaster == chatter {
			tokens := strings.Split(chat.Body, " ")
			if len(tokens) == 0 {
				return
			}
			roomID := tokens[0]

			log.Println("setting roomid", broadcaster, roomID)
			loader.TwitchSetRoom(broadcaster, roomID)
		} else {
			log.Println("unauthorized user tried to setboard")
		}
	} else if chat.Command == "branch" {
		log.Println("grafting:", chat.Body)
		branch := strings.ToLower(chat.Body)

		// create the event
		e := &core.EventJSON{
			Event: "graft",
			Value: branch,
		}

		roomID := loader.TwitchGetRoom(broadcaster)
		if roomID == "" {
			log.Println("room not set for", broadcaster)
			return
		}
		log.Println("room found for", broadcaster, roomID)
		room := s.GetOrCreateRoom(roomID)

		handler := Chain(
			room.HandleEvent,
			room.BroadcastFullFrameAfter)
		handler(e)
	}

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
		evt = &core.EventJSON{
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
		evt = &core.EventJSON{
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
