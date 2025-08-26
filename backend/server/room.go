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
	"log"

	"github.com/jarednogo/board/backend/core"
	"github.com/jarednogo/board/backend/ogs"
	"github.com/jarednogo/board/backend/socket"
	"github.com/jarednogo/board/backend/state"
	"golang.org/x/net/websocket"
	"time"
)

type Room struct {
	Conns         map[string]*websocket.Conn
	State         *state.State
	TimeLastEvent *time.Time
	LastUser      string
	lastMessages  map[string]*time.Time
	OGSLink       *ogs.OGSConnector
	Password      string
	auth          map[string]bool
	Nicks         map[string]string
}

func NewRoom() *Room {
	conns := make(map[string]*websocket.Conn)
	s := state.NewState(19, true)
	now := time.Now()
	msgs := make(map[string]*time.Time)
	auth := make(map[string]bool)
	nicks := make(map[string]string)
	return &Room{conns, s, &now, "", msgs, nil, "", auth, nicks}
}

func (r *Room) HasPassword() bool {
	return r.Password != ""
}

func (r *Room) SendTo(id string, evt *core.EventJSON) {
	if ws, ok := r.Conns[id]; ok {
		socket.SendEvent(ws, evt)
	}
}

func (r *Room) Broadcast(evt *core.EventJSON, setTime bool) {
	if evt.Event == "nop" {
		return
	}
	// augment event with connection id
	id := evt.UserID

	// marshal event back into data
	data, err := json.Marshal(evt)
	if err != nil {
		log.Println(id, err)
		return
	}

	// rebroadcast message
	for _, conn := range r.Conns {
		conn.Write(data)
	}

	if setTime {
		// set last user information
		r.LastUser = id
		now := time.Now()
		r.TimeLastEvent = &now
	}
}

func (r *Room) UploadSGF(sgf string) *core.EventJSON {
	s, err := state.FromSGF(sgf)
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("Error parsing SGF: %s", err)
		return core.ErrorJSON(msg)
	}
	r.State = s

	// replace evt with initdata
	frame := r.State.GenerateFullFrame(core.Full)
	return core.FrameJSON(frame)
}

func (r *Room) SendUserList() {
	// send list of currently connected users
	evt := &core.EventJSON{
		"connected_users",
		r.Nicks,
		0,
		"",
	}

	r.Broadcast(evt, false)
}
