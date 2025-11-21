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

type Room struct {
	conns        map[string]event.EventChannel
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
	r.SetPassword(load.Password)
	r.SetInputBuffer(load.Buffer)
	r.SetState(st)

	return r, nil
}

func (r *Room) ID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.id
}

func (r *Room) GetLastMessages(user string) (time.Time, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var tnil time.Time
	t, ok := r.lastMessages[user]
	if !ok {
		return tnil, ok
	}
	return *t, ok
}

func (r *Room) SetLastMessages(user string, t *time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastMessages[user] = t
}

func (r *Room) GetLastActive() time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	return *r.lastActive
}

func (r *Room) SetLastActive(t *time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastActive = t
}

func (r *Room) GetLastUser() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lastUser
}

func (r *Room) SetLastUser(user string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastUser = user
}

func (r *Room) SaveState() *state.StateJSON {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state.Save()
}

func (r *Room) GetInputBuffer() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.inputBuffer
}

func (r *Room) SetInputBuffer(i int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.inputBuffer = i
}

func (r *Room) SetUserBuffer(i int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.userBuffer = i
}

func (r *Room) GetUserBuffer() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.userBuffer
}

func (r *Room) SetAuth(user string, val bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.auth[user] = val
}

func (r *Room) GetAuth(user string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.auth[user]
	return ok
}

func (r *Room) SetAuthAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for connID := range r.conns {
		r.auth[connID] = true
	}
}

func (r *Room) GetTimeout() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.timeout
}

func (r *Room) SetTimeout(f float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.timeout = f
}

func (r *Room) Nicks() map[string]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.nicks
}

func (r *Room) GetNick(user string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.nicks[user]
	return n, ok
}

func (r *Room) SetNick(id, nick string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nicks[id] = nick
}

func (r *Room) Save() *loader.LoadJSON {
	r.mu.Lock()
	defer r.mu.Unlock()

	stateJSON := r.state.Save()

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
	r.mu.Lock()
	defer r.mu.Unlock()
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

func (r *Room) SetState(s *state.State) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state = s
}

func (r *Room) SetFetcher(f fetch.Fetcher) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fetcher = f
}

func (r *Room) GetFetcher() fetch.Fetcher {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.fetcher
}

func (r *Room) HasPassword() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.password != ""
}

func (r *Room) GetPassword() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.password
}

func (r *Room) SetPassword(p string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.password = p
}

func (r *Room) SendTo(id string, evt event.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ec, ok := r.conns[id]; ok {
		ec.SendEvent(evt) //nolint:errcheck
	}
}

func (r *Room) Broadcast(evt event.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if evt.Type() == "nop" {
		return
	}

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
	r.SetState(s)

	// replace evt with frame data
	frame := r.GenerateFullFrame(core.Full)
	return event.FrameEvent(frame)
}

// used for testing
func (r *Room) DisableBuffers() {
	r.SetInputBuffer(0)
	r.SetUserBuffer(0)
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

func (room *Room) GetHandler(t string) EventHandler {
	room.mu.Lock()
	defer room.mu.Unlock()
	h, ok := room.handlers[t]
	if ok {
		return h
	}
	return room.handlers["_"]
}

func (r *Room) RegisterPlugin(p plugin.Plugin, args map[string]any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := args["key"].(string)
	p.Start(args)
	r.plugins[key] = p
}

func (r *Room) DeregisterPlugin(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.plugins[key]; ok {
		p.End()
		delete(r.plugins, key)
	}
}

func (r *Room) GetPlugin(key string) plugin.Plugin {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.plugins[key]; ok {
		return p
	}
	return nil
}

func (r *Room) HasPlugin(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.plugins[key]
	return ok
}

func (r *Room) Execute(cmd state.Command) (*core.Frame, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return cmd.Execute(r.state)
}

// wrappers around the state
func (r *Room) Size() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state.Size()
}

func (r *Room) GenerateFullFrame(t core.TreeJSONType) *core.Frame {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state.GenerateFullFrame(t)
}

func (r *Room) EditPlayerBlack(s string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state.EditPlayerBlack(s)
}

func (r *Room) EditPlayerWhite(s string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state.EditPlayerWhite(s)
}

func (r *Room) EditKomi(s string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state.EditKomi(s)
}

func (r *Room) AddStones(moves []*core.Stone) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state.AddStones(moves)
}

func (r *Room) HeadColor() core.Color {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state.HeadColor()
}

func (r *Room) PushHead(x, y int, col core.Color) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state.PushHead(x, y, col)
}

func (r *Room) ToSGF() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state.ToSGF()
}

func (r *Room) ToSGFIX() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state.ToSGFIX()
}

func (r *Room) Current() *core.TreeNode {
	r.mu.Lock()
	defer r.mu.Unlock()
	c := *r.state.Current()
	return &c
}

func (r *Room) Board() *core.Board {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state.Board().Copy()
}
