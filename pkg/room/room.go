/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room

import (
	"fmt"
	"log"
	"time"

	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/fetch"
	"github.com/jarednogo/board/pkg/room/plugin"
	"github.com/jarednogo/board/pkg/socket"
	"github.com/jarednogo/board/pkg/state"
)

type Room struct {
	Conns         map[string]socket.RoomConn
	State         *state.State
	TimeLastEvent *time.Time
	LastUser      string
	lastMessages  map[string]*time.Time
	Plugins       map[string]plugin.Plugin
	Password      string
	auth          map[string]bool
	Nicks         map[string]string
	fetcher       fetch.Fetcher
}

func NewRoom() *Room {
	conns := make(map[string]socket.RoomConn)
	s := state.NewState(19, true)
	now := time.Now()
	msgs := make(map[string]*time.Time)
	auth := make(map[string]bool)
	nicks := make(map[string]string)
	plugins := make(map[string]plugin.Plugin)
	return &Room{
		Conns:         conns,
		State:         s,
		TimeLastEvent: &now,
		LastUser:      "",
		lastMessages:  msgs,
		Plugins:       plugins,
		Password:      "",
		auth:          auth,
		Nicks:         nicks,
		fetcher:       fetch.NewDefaultFetcher(),
	}
}

func (r *Room) SetFetcher(f fetch.Fetcher) {
	r.fetcher = f
}

func (r *Room) HasPassword() bool {
	return r.Password != ""
}

func (r *Room) SendTo(id string, evt *core.EventJSON) {
	if rc, ok := r.Conns[id]; ok {
		rc.SendEvent(evt)
	}
}

func (r *Room) Broadcast(evt *core.EventJSON) {
	if evt.Event == "nop" {
		return
	}

	// rebroadcast message
	for _, conn := range r.Conns {
		conn.SendEvent(evt)
	}
}

func (r *Room) BroadcastHubMessage(m *core.Message) {
	// make a new event to broadcast
	evt := &core.EventJSON{
		Event:  "global",
		Value:  m.Text,
		Color:  0,
		UserID: "",
	}

	// go through each client connection
	for id, conn := range r.Conns {
		// check to see if we've already sent this message
		// to this connection
		if m.IsNotified(id) {
			continue
		}
		// otherwise, send and record
		conn.SendEvent(evt)
		m.MarkNotified(id)
	}
}

func (r *Room) UploadSGF(sgf string) *core.EventJSON {
	s, err := state.FromSGF(sgf)
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("Error parsing SGF: %s", err)
		return core.ErrorEvent(msg)
	}
	r.State = s

	// replace evt with frame data
	frame := r.State.GenerateFullFrame(core.Full)
	return core.FrameEvent(frame)
}

func (r *Room) SendUserList() {
	// send list of currently connected users
	evt := &core.EventJSON{
		Event:  "connected_users",
		Value:  r.Nicks,
		Color:  0,
		UserID: "",
	}

	r.Broadcast(evt)
}

func (r *Room) RegisterConnection(rc socket.RoomConn) string {
	// the room connection generates its own id
	id := rc.GetID()

	// set the last user
	r.LastUser = id

	// store the new connection by id
	r.Conns[id] = rc

	// save current user
	r.Nicks[id] = ""

	// send initial state
	frame := r.State.GenerateFullFrame(core.Full)
	evt := core.FrameEvent(frame)
	rc.SendEvent(evt)

	return id
}

func (r *Room) DeregisterConnection(id string) {
	delete(r.Conns, id)
}

func (r *Room) Handle(rc socket.RoomConn) {
	// assign id to the new connection
	id := r.RegisterConnection(rc)

	// defer removing the client
	defer r.DeregisterConnection(id)

	// send list of currently connected users
	r.SendUserList()

	// send disconnection notification
	// golang deferrals are called in LIFO order
	defer r.SendUserList()
	defer delete(r.Nicks, id)

	handlers := r.CreateHandlers()

	// main loop
	for {
		// receive the event
		evt, err := rc.ReceiveEvent()
		if err != nil {
			log.Println(id, err)
			break
		}

		// augment with user id
		evt.UserID = id

		// handle the event
		if handler, ok := handlers[evt.Event]; ok {
			handler(evt)
		} else {
			handlers["_"](evt)
		}
	}
}

func (r *Room) RegisterPlugin(p plugin.Plugin, args map[string]interface{}) {
	key := args["key"].(string)
	p.Start(args)
	r.Plugins[key] = p
}

func (r *Room) DeregisterPlugin(key string) {
	if p, ok := r.Plugins[key]; ok {
		p.End()
		delete(r.Plugins, key)
	}
}
