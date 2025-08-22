/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
	"time"
)

type Room struct {
	conns         map[string]*websocket.Conn
	State         *State
	timeLastEvent *time.Time
	lastUser      string
	lastMessages  map[string]*time.Time
	open          bool
	OGSLink       *OGSConnector
	password      string
	auth          map[string]bool
	nicks         map[string]string
}

func NewRoom() *Room {
	conns := make(map[string]*websocket.Conn)
	state := NewState(19, true)
	now := time.Now()
	msgs := make(map[string]*time.Time)
	auth := make(map[string]bool)
	nicks := make(map[string]string)
	return &Room{conns, state, &now, "", msgs, true, nil, "", auth, nicks}
}

func (r *Room) HasPassword() bool {
	return r.password != ""
}

func (r *Room) SendTo(id string, evt *EventJSON) {
	if ws, ok := r.conns[id]; ok {
		SendEvent(ws, evt)
	}
}

func (r *Room) Broadcast(evt *EventJSON, setTime bool) {
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
	for _, conn := range r.conns {
		conn.Write(data)
	}

	if setTime {
		// set last user information
		r.lastUser = id
		now := time.Now()
		r.timeLastEvent = &now
	}
}

func (r *Room) PushHead(x, y, col int) *EventJSON {
	r.State.PushHead(x, y, col)
	evt := &EventJSON{
		Event:  "push_head",
		Value:  []int{x, y},
		Color:  col,
		UserID: "",
	}
	return evt
}

func (r *Room) UploadSGF(sgf string) *EventJSON {
	state, err := FromSGF(sgf)
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("Error parsing SGF: %s", err)
		return ErrorJSON(msg)
	}
	r.State = state

	// replace evt with initdata
	frame := r.State.GenerateFullFrame(Full)
	return FrameJSON(frame)
}

func (r *Room) SendUserList() {
	// send list of currently connected users
	evt := &EventJSON{
		"connected_users",
		r.nicks,
		0,
		"",
	}

	r.Broadcast(evt, false)
}
