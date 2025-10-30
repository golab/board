/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/fetch"
	"github.com/jarednogo/board/pkg/loader"
	"github.com/jarednogo/board/pkg/room/plugin"
	"github.com/jarednogo/board/pkg/socket"
	"github.com/jarednogo/board/pkg/state"
)

type Room struct {
	conns        map[string]socket.RoomConn
	state        *state.State
	lastActive   *time.Time
	lastUser     string
	lastMessages map[string]*time.Time
	plugins      map[string]plugin.Plugin
	password     string
	auth         map[string]bool
	nicks        map[string]string
	fetcher      fetch.Fetcher
	id           string
}

func NewRoom(id string) *Room {
	conns := make(map[string]socket.RoomConn)
	s := state.NewState(19, true)
	now := time.Now()
	msgs := make(map[string]*time.Time)
	auth := make(map[string]bool)
	nicks := make(map[string]string)
	plugins := make(map[string]plugin.Plugin)
	return &Room{
		conns:        conns,
		state:        s,
		lastActive:   &now,
		lastUser:     "",
		lastMessages: msgs,
		plugins:      plugins,
		password:     "",
		auth:         auth,
		nicks:        nicks,
		fetcher:      fetch.NewDefaultFetcher(),
		id:           id,
	}
}

func Load(load *loader.LoadJSON) (*Room, error) {
	id := load.ID

	sgf, err := base64.StdEncoding.DecodeString(load.SGF)
	if err != nil {
		return nil, err
	}

	st, err := state.FromSGF(string(sgf))
	if err != nil {
		return nil, err
	}

	st.SetPrefs(load.Prefs)

	st.SetNextIndex(load.NextIndex)
	st.SetInputBuffer(load.Buffer)

	loc := load.Location
	if loc != "" {
		dirs := strings.Split(loc, ",")
		// don't need to assign to a variable if we don't use it
		for range dirs {
			st.Right()
		}
	}
	r := NewRoom(id)
	r.password = load.Password
	r.state = st

	return r, nil
}

func (r *Room) ID() string {
	return r.id
}

func (r *Room) Timeout() float64 {
	return r.state.GetTimeout()
}

func (r *Room) LastActive() *time.Time {
	return r.lastActive
}

func (r *Room) Save() *loader.LoadJSON {
	stateJSON := r.state.CreateStateJSON()

	// embed stateJSON into save
	save := &loader.LoadJSON{}
	save.SGF = stateJSON.SGF
	save.Location = stateJSON.Location
	save.Prefs = stateJSON.Prefs
	save.Buffer = stateJSON.Buffer
	save.NextIndex = stateJSON.NextIndex

	// add on last fields owned by the room instead of the state
	save.Password = r.password
	save.ID = r.id
	return save
}

func (r *Room) Close() error {
	// close all the client connections
	errs := []error{}
	for _, conn := range r.conns {
		err := conn.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

func (r *Room) SetFetcher(f fetch.Fetcher) {
	r.fetcher = f
}

func (r *Room) HasPassword() bool {
	return r.password != ""
}

func (r *Room) GetPassword() string {
	return r.password
}

func (r *Room) SetPassword(p string) {
	r.password = p
}

func (r *Room) SendTo(id string, evt *core.EventJSON) {
	if rc, ok := r.conns[id]; ok {
		rc.SendEvent(evt) //nolint:errcheck
	}
}

func (r *Room) Broadcast(evt *core.EventJSON) {
	if evt.Event == "nop" {
		return
	}

	// rebroadcast message
	for _, conn := range r.conns {
		conn.SendEvent(evt) //nolint:errcheck
	}
}

func (r *Room) BroadcastHubMessage(m *core.Message) {
	// make a new event to broadcast
	evt := &core.EventJSON{
		Event:  "global",
		Value:  m.Text,
		UserID: "",
	}

	// go through each client connection
	for id, conn := range r.conns {
		// check to see if we've already sent this message
		// to this connection
		if m.IsNotified(id) {
			continue
		}
		// otherwise, send and record
		conn.SendEvent(evt) //nolint:errcheck
		m.MarkNotified(id)
	}
}

func (r *Room) UploadSGF(sgf string) *core.EventJSON {
	s, err := state.FromSGF(sgf)
	if err != nil {
		msg := fmt.Sprintf("Error parsing SGF: %s", err)
		return core.ErrorEvent(msg)
	}
	r.state = s

	// replace evt with frame data
	frame := r.state.GenerateFullFrame(core.Full)
	return core.FrameEvent(frame)
}

func (r *Room) SendUserList() {
	// send list of currently connected users
	evt := &core.EventJSON{
		Event:  "connected_users",
		Value:  r.nicks,
		UserID: "",
	}

	r.Broadcast(evt)
}

func (r *Room) RegisterConnection(rc socket.RoomConn) string {
	// the room connection generates its own id
	id := rc.ID()

	// set the last user
	r.lastUser = id

	// store the new connection by id
	r.conns[id] = rc

	// save current user
	r.nicks[id] = ""

	// send initial state
	frame := r.state.GenerateFullFrame(core.Full)
	evt := core.FrameEvent(frame)
	rc.SendEvent(evt) //nolint:errcheck

	return id
}

func (r *Room) DeregisterConnection(id string) {
	delete(r.conns, id)
}

func (r *Room) Handle(rc socket.RoomConn) error {
	// assign id to the new connection
	id := r.RegisterConnection(rc)

	// defer removing the client
	defer r.DeregisterConnection(id)

	// send list of currently connected users
	r.SendUserList()

	// send disconnection notification
	// golang deferrals are called in LIFO order
	defer r.SendUserList()
	defer delete(r.nicks, id)

	handlers := r.CreateHandlers()

	// main loop
	for {
		// receive the event
		evt, err := rc.ReceiveEvent()
		if err != nil {
			return err
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
	r.plugins[key] = p
}

func (r *Room) DeregisterPlugin(key string) {
	if p, ok := r.plugins[key]; ok {
		p.End()
		delete(r.plugins, key)
	}
}

func (r *Room) GetPlugin(key string) plugin.Plugin {
	if p, ok := r.plugins[key]; ok {
		return p
	}
	return nil
}

func (r *Room) HasPlugin(key string) bool {
	_, ok := r.plugins[key]
	return ok
}
