/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub

import (
	"encoding/base64"
	"log"
	"strings"
	"time"

	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/loader"
	"github.com/jarednogo/board/pkg/room"
	"github.com/jarednogo/board/pkg/socket"
	"github.com/jarednogo/board/pkg/state"
	"golang.org/x/net/websocket"
)

// the url starts with '/'
func ParseURL(url string) (string, string, string) {
	url = strings.TrimPrefix(url, "/socket")
	tokens := strings.Split(url, "/")
	if len(tokens) <= 1 {
		return "", "", ""
	}
	if len(tokens) == 2 {
		return tokens[1], "", ""
	}
	if len(tokens) == 3 {
		return tokens[1], tokens[2], ""
	}
	return tokens[1], tokens[2], tokens[3]
}

type Hub struct {
	rooms    map[string]*room.Room
	messages []*core.Message
	db       loader.Loader
}

func NewHub() *Hub {
	// get database setup
	db := loader.NewDefaultLoader()
	return NewHubWithDB(db)
}

func NewHubWithDB(db loader.Loader) *Hub {
	db.Setup()
	s := &Hub{
		rooms:    make(map[string]*room.Room),
		messages: []*core.Message{},
		db:       db,
	}

	// start message loop
	go s.MessageLoop()

	return s
}

func (h *Hub) Save() {
	for id, r := range h.rooms {
		log.Printf("Saving %s", id)

		stateJSON := r.State.CreateStateJSON()

		// embed stateJSON into save
		save := &loader.LoadJSON{}
		save.SGF = stateJSON.SGF
		save.Location = stateJSON.Location
		save.Prefs = stateJSON.Prefs
		save.Buffer = stateJSON.Buffer
		save.NextIndex = stateJSON.NextIndex

		// add on last fields owned by the room instead of the state
		save.Password = r.Password
		save.ID = id

		err := h.db.SaveRoom(id, save)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *Hub) Load() {
	rooms := h.db.LoadAllRooms()

	for _, load := range rooms {

		id := load.ID

		sgf, err := base64.StdEncoding.DecodeString(load.SGF)
		if err != nil {
			continue
		}

		st, err := state.FromSGF(string(sgf))
		if err != nil {
			continue
		}

		st.SetPrefs(load.Prefs)

		st.NextIndex = load.NextIndex
		st.InputBuffer = load.Buffer

		loc := load.Location
		if loc != "" {
			dirs := strings.Split(loc, ",")
			// don't need to assign to a variable if we don't use it
			for range dirs {
				st.Right()
			}
		}

		log.Printf("Loading %s", id)

		r := room.NewRoom()
		r.Password = load.Password
		r.State = st
		h.rooms[id] = r
		go h.Heartbeat(id)
	}
}

func (h *Hub) Heartbeat(roomID string) {
	r, ok := h.rooms[roomID]
	if !ok {
		return
	}
	for {
		now := time.Now()
		diff := now.Sub(*r.TimeLastEvent)
		log.Println(roomID, "Inactive for", diff)
		if diff.Seconds() > r.State.Timeout {
			break
		}
		time.Sleep(3600 * time.Second)
	}
	log.Println("Cleaning up board due to inactivity:", roomID)

	// close all the client connections
	for _, conn := range r.Conns {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}

	// delete the room from the server map
	delete(h.rooms, roomID)

	// delete it from the database
	err := h.db.DeleteRoom(roomID)
	if err != nil {
		log.Println(err)
	}
}

func (h *Hub) ReadMessages() {
	messages := h.db.LoadAllMessages()
	defer h.db.DeleteAllMessages()

	for _, msg := range messages {
		m := core.NewMessage(msg.Text, msg.TTL)
		h.messages = append(h.messages, m)
	}
}

func (h *Hub) SendMessages() {
	// go through each server message
	keep := []*core.Message{}
	for _, m := range h.messages {
		// check time
		now := time.Now()

		// skip the expired messages
		if m.ExpiresAt.Before(now) {
			continue
		}

		// keep the unexpired messages
		keep = append(keep, m)

		// go through each room
		for _, r := range h.rooms {
			r.BroadcastHubMessage(m)
		}
	}
	// save the unexpired messages
	h.messages = keep
}

func (h *Hub) MessageLoop() {
	for {
		// wait 5 seconds
		time.Sleep(5 * time.Second)

		h.ReadMessages()
		h.SendMessages()
	}
}

func (h *Hub) GetOrCreateRoom(roomID string) *room.Room {
	// if the room they want doesn't exist, create it
	//created := false
	if _, ok := h.rooms[roomID]; !ok {
		//created = true
		log.Println("New room:", roomID)
		r := room.NewRoom()
		h.rooms[roomID] = r
		go h.Heartbeat(roomID)
	}
	r := h.rooms[roomID]
	return r
}

func (h *Hub) HandlerWrapper(ws *websocket.Conn) {

	// first find the url they want
	url := ws.Request().URL.String()

	log.Println(url)

	// currently not using the prefix, but i may someday
	_, roomID, _ := ParseURL(url)

	// wrap the socket
	wrc := socket.NewWebsocketRoomConn(ws)
	h.Handler(wrc, roomID)
}

func (h *Hub) Handler(rc socket.RoomConn, roomID string) {
	// new connection

	// get or create the room
	r := h.GetOrCreateRoom(roomID)

	// send to the room for handling
	log.Println("Connecting:", rc.GetID())
	r.Handle(rc)
}
