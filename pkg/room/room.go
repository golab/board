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
	"sync"
	"time"

	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/event"
	"github.com/jarednogo/board/pkg/fetch"
	"github.com/jarednogo/board/pkg/loader"
	"github.com/jarednogo/board/pkg/message"
	"github.com/jarednogo/board/pkg/room/plugin"
	"github.com/jarednogo/board/pkg/state"
)

type engine = state.State

type Room struct {
	conns map[string]event.EventChannel
	*engine
	lastActive   *time.Time
	lastUser     string
	lastMessages map[string]*time.Time
	plugins      map[string]plugin.Plugin
	password     string
	auth         map[string]bool
	nicks        map[string]string
	fetcher      fetch.Fetcher
	id           string
	handlers     map[string]EventHandler
	inputBuffer  int64
	userBuffer   int64
	timeout      float64
	mu           sync.Mutex
}

func NewRoom(id string) *Room {
	conns := make(map[string]event.EventChannel)
	s := state.NewState(19, true)
	now := time.Now()
	msgs := make(map[string]*time.Time)
	auth := make(map[string]bool)
	nicks := make(map[string]string)
	plugins := make(map[string]plugin.Plugin)
	// default input buffer of 250
	// default timeout of 86400

	r := &Room{
		conns:        conns,
		engine:       s,
		lastActive:   &now,
		lastUser:     "",
		lastMessages: msgs,
		plugins:      plugins,
		password:     "",
		auth:         auth,
		nicks:        nicks,
		fetcher:      fetch.NewDefaultFetcher(),
		id:           id,
		inputBuffer:  250,
		userBuffer:   100,
		timeout:      86400,
	}
	r.initHandlers()

	return r
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

	loc := load.Location
	st.SetLocation(loc)
	r := NewRoom(id)
	r.password = load.Password
	r.inputBuffer = load.Buffer
	r.setState(st)

	return r, nil
}

func (r *Room) ID() string {
	return r.id
}

func (r *Room) LastActive() *time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lastActive
}

func (r *Room) SaveState() *state.StateJSON {
	return r.engine.Save()
}

func (r *Room) GetInputBuffer() int64 {
	return r.inputBuffer
}

func (r *Room) SetInputBuffer(i int64) {
	r.inputBuffer = i
}

func (r *Room) SetUserBuffer(i int64) {
	r.userBuffer = i
}

// used for testing
func (r *Room) DisableBuffers() {
	r.SetInputBuffer(0)
	r.SetUserBuffer(0)
}

func (r *Room) GetTimeout() float64 {
	return r.timeout
}

func (r *Room) SetTimeout(f float64) {
	r.timeout = f
}

func (r *Room) Nicks() map[string]string {
	return r.nicks
}

func (r *Room) Save() *loader.LoadJSON {
	stateJSON := r.engine.Save()

	// embed stateJSON into save
	save := &loader.LoadJSON{}
	save.SGF = stateJSON.SGF
	save.Location = stateJSON.Location
	save.Prefs = stateJSON.Prefs
	save.NextIndex = stateJSON.NextIndex

	// add on last fields owned by the room instead of the state
	save.Password = r.password
	save.Buffer = r.inputBuffer
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

func (r *Room) setState(s *state.State) {
	r.engine = s
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

func (r *Room) SendTo(id string, evt event.Event) {
	if ec, ok := r.conns[id]; ok {
		ec.SendEvent(evt) //nolint:errcheck
	}
}

func (r *Room) Broadcast(evt event.Event) {
	if evt.Type() == "nop" {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	// rebroadcast message
	for _, conn := range r.conns {
		conn.SendEvent(evt) //nolint:errcheck
	}
}

func (r *Room) BroadcastHubMessage(m *message.Message) {
	// make a new event to broadcast
	evt := event.NewEvent("global", m.Text)

	r.mu.Lock()
	defer r.mu.Unlock()
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

func (r *Room) UploadSGF(sgf string) event.Event {
	s, err := state.FromSGF(sgf)
	if err != nil {
		msg := fmt.Sprintf("Error parsing SGF: %s", err)
		return event.ErrorEvent(msg)
	}
	r.engine = s

	// replace evt with frame data
	frame := r.GenerateFullFrame(core.Full)
	return event.FrameEvent(frame)
}

func (r *Room) SendUserList() {
	// send list of currently connected users
	evt := event.NewEvent("connected_users", r.nicks)

	r.Broadcast(evt)
}

func (r *Room) RegisterConnection(ec event.EventChannel) string {
	// currently a no-op, but useful for testing
	ec.OnConnect()

	// the room connection generates its own id
	id := ec.ID()

	r.mu.Lock()

	// set the last user
	r.lastUser = id

	// store the new connection by id
	r.conns[id] = ec

	// save current user
	r.nicks[id] = ""

	r.mu.Unlock()

	// send initial state
	frame := r.GenerateFullFrame(core.Full)
	evt := event.FrameEvent(frame)
	ec.SendEvent(evt) //nolint:errcheck

	return id
}

func (r *Room) DeregisterConnection(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.conns, id)
}

func (r *Room) Handle(ec event.EventChannel) error {
	// assign id to the new connection
	id := r.RegisterConnection(ec)

	// defer removing the client
	defer r.SendUserList()
	defer r.DeregisterConnection(id)

	// send list of currently connected users
	r.SendUserList()

	// send disconnection notification
	defer func() {
		r.mu.Lock()
		delete(r.nicks, id)
		r.mu.Unlock()
	}()

	// main loop
	for {
		// receive the event
		evt, err := ec.ReceiveEvent()
		if err != nil {
			return err
		}

		// augment with user id
		evt.SetUser(id)

		// handle the event
		r.HandleAny(evt)
	}
}

func (r *Room) RegisterPlugin(p plugin.Plugin, args map[string]any) {
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
