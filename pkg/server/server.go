/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package server

import (
	"encoding/base64"
	"encoding/json"
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

type Server struct {
	rooms    map[string]*room.Room
	messages []*Message
	db       loader.Loader
}

func NewServer() *Server {
	// get database setup
	db := loader.NewDefaultLoader()
	db.Setup()

	s := &Server{
		rooms:    make(map[string]*room.Room),
		messages: []*Message{},
		db:       db,
	}

	// start message loop
	go s.MessageLoop()

	return s
}

func (s *Server) Save() {
	for id, r := range s.rooms {
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

		err := s.db.SaveRoom(id, save)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) Load() {
	rooms := s.db.LoadAllRooms()

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
		s.rooms[id] = r
		go s.Heartbeat(id)
	}
}

func (s *Server) Heartbeat(roomID string) {
	r, ok := s.rooms[roomID]
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
	delete(s.rooms, roomID)

	// delete it from the database
	err := s.db.DeleteRoom(roomID)
	if err != nil {
		log.Println(err)
	}
}

type Message struct {
	Text      string
	ExpiresAt *time.Time
	Notified  map[string]bool
}

func (s *Server) ReadMessages() {

	messages := s.db.LoadAllMessages()
	defer s.db.DeleteAllMessages()

	for _, msg := range messages {
		// calculate the expiration time using TTL
		now := time.Now()
		expiresAt := now.Add(time.Duration(msg.TTL) * time.Second)

		// add to server messages
		m := &Message{msg.Text, &expiresAt, make(map[string]bool)}
		s.messages = append(s.messages, m)
	}
}

func (s *Server) SendMessages() {
	// go through each server message
	keep := []*Message{}
	for _, m := range s.messages {
		// check time
		now := time.Now()

		// skip the expired messages
		if m.ExpiresAt.Before(now) {
			continue
		}

		// keep the unexpired messages
		keep = append(keep, m)

		// make a new event to broadcast
		evt := &core.EventJSON{
			Event:  "global",
			Value:  m.Text,
			Color:  0,
			UserID: "",
		}

		// go through each room
		for _, r := range s.rooms {
			// go through each client connection
			for id, conn := range r.Conns {
				// check to see if we've already sent this message
				// to this connection
				if m.Notified[id] {
					continue
				}
				// otherwise, send and record
				conn.SendEvent(evt)
				m.Notified[id] = true
			}
		}
	}
	// save the unexpired messages
	s.messages = keep
}

func (s *Server) MessageLoop() {
	for {
		// wait 5 seconds
		time.Sleep(5 * time.Second)

		s.ReadMessages()
		s.SendMessages()
	}
}

func (s *Server) SendMessagesToOne(rc socket.RoomConn, id string) {
	// send messages
	for _, m := range s.messages {
		// make a new event to send
		evt := &core.EventJSON{
			Event:  "global",
			Value:  m.Text,
			Color:  0,
			UserID: "",
		}

		rc.SendEvent(evt)
		m.Notified[id] = true
	}
}

func (s *Server) GetOrCreateRoom(roomID string) *room.Room {
	// if the room they want doesn't exist, create it
	//created := false
	if _, ok := s.rooms[roomID]; !ok {
		//created = true
		log.Println("New room:", roomID)
		r := room.NewRoom()
		s.rooms[roomID] = r
		go s.Heartbeat(roomID)
	}
	r := s.rooms[roomID]
	//return r, created
	return r
}

// these operate outside the main websocket loop
// see frontend/main.go suffixOp for corresponding receiver
func (s *Server) HandleOp(op, roomID string) string {
	data := ""
	r, ok := s.rooms[roomID]
	if !ok {
		return ""
	}
	switch op {
	case "sgf":
		// if the room doesn't exist, send empty string
		data = r.State.ToSGF(false)
	case "sgfix":
		// basically do the same thing but include indexes
		data = r.State.ToSGF(true)
	case "debug":
		// send debug info
		stateJSON := r.State.CreateStateJSON()
		dataBytes, _ := json.Marshal(stateJSON)
		data = string(dataBytes)
	}
	return data
}

// Echo the data received on the WebSocket.
func (s *Server) Handler(ws *websocket.Conn) {
	// new connection

	// first find the url they want
	url := ws.Request().URL.String()

	log.Println(url)

	// currently not using the prefix, but i may someday
	_, roomID, _ := ParseURL(url)

	// get or create the room
	r := s.GetOrCreateRoom(roomID)

	// assign id to the new connection
	wrc := socket.NewWebsocketRoomConn(ws)
	id := r.RegisterConnection(wrc)
	log.Println(url, "Connecting:", id)
	s.SendMessagesToOne(wrc, id)

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
		evt, err := wrc.ReceiveEvent()
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
