/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

// the url starts with '/'
func ParseURL(url string) (string, string, string) {
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
	rooms    map[string]*Room
	messages []*Message
}

func NewServer() *Server {
	return &Server{
		make(map[string]*Room),
		[]*Message{},
	}
}

func (s *Server) Save() {
	for id, room := range s.rooms {
		path := filepath.Join(RoomPath(), id)
		log.Printf("Saving %s", path)

		// TODO: the term "handshake" has become obsolete
		// as this is not the same process for client handshakes anymore
		evt := room.State.InitData("handshake")
		dataStruct := &LoadJSON{}
		s, _ := evt.Value.(string)
		err := json.Unmarshal([]byte(s), dataStruct)
		if err != nil {
			continue
		}

		dataStruct.Password = room.password
		data, err := json.Marshal(dataStruct)
		if err != nil {
			continue
		}

		err = ioutil.WriteFile(path, data, 0644)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) Load() {
	dir := RoomPath()
	sgfs, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range sgfs {
		log.Println(e)
		id := e.Name()
		path := filepath.Join(dir, id)
		data, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}

		load := &LoadJSON{}
		err = json.Unmarshal(data, load)
		if err != nil {
			continue
		}

		sgf, err := base64.StdEncoding.DecodeString(load.SGF)
		if err != nil {
			continue
		}

		state, err := FromSGF(string(sgf))
		if err != nil {
			continue
		}

		state.SetPrefs(load.Prefs)

		state.NextIndex = load.NextIndex
		state.InputBuffer = load.Buffer

		loc := load.Loc
		if loc != "" {
			dirs := strings.Split(loc, ",")
			for _ = range dirs {
				state.Right()
			}
		}

		log.Printf("Loading %s", path)

		r := NewRoom()
		r.password = load.Password
		r.State = state
		s.rooms[id] = r
		go s.Heartbeat(id)
	}
}

func (s *Server) Heartbeat(roomID string) {
	room, ok := s.rooms[roomID]
	if !ok {
		return
	}
	for {
		now := time.Now()
		diff := now.Sub(*room.timeLastEvent)
		log.Println(roomID, "Inactive for", diff)
		if diff.Seconds() > room.State.Timeout {
			room.open = false
			break
		}
		time.Sleep(3600 * time.Second)
	}
	log.Println("Cleaning up board due to inactivity:", roomID)

	// close all the client connections
	for _, conn := range room.conns {
		conn.Close()
	}

	// delete the room from the server map
	delete(s.rooms, roomID)

	// delete the saved file (if it exists)
	path := filepath.Join(RoomPath(), roomID)
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}
}

type MessageJSON struct {
	Text string `json:"message"`
	TTL  int    `json:"ttl"`
}

type Message struct {
	Text      string
	ExpiresAt *time.Time
	Notified  map[string]bool
}

func (s *Server) ReadMessages() {
	// iterate through all files in the message path
	files, err := os.ReadDir(MessagePath())
	if err != nil {
		log.Fatal(err)
	}

	// check messages
	for _, file := range files {
		// might do something someday with nested directories
		if file.IsDir() {
			continue
		}

		// read each file
		path := filepath.Join(MessagePath(), file.Name())
		data, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}

		// convert json to struct
		msg := &MessageJSON{}
		err = json.Unmarshal(data, msg)
		if err != nil {
			continue
		}

		// remove the file
		os.Remove(path)

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
		evt := &EventJSON{
			"global",
			m.Text,
			0,
			"",
		}

		// go through each room
		for _, room := range s.rooms {
			// go through each client connection
			for id, conn := range room.conns {
				// check to see if we've already sent this message
				// to this connection
				if m.Notified[id] {
					continue
				}
				// otherwise, send and record
				SendEvent(conn, evt)
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

func (s *Server) SendMessagesToOne(ws *websocket.Conn, id string) {
	// send messages
	for _, m := range s.messages {
		// make a new event to send
		evt := &EventJSON{
			"global",
			m.Text,
			0,
			"",
		}

		SendEvent(ws, evt)
		m.Notified[id] = true
	}
}

func (s *Server) GetOrCreateRoom(roomID string) (*Room, bool) {
	// if the room they want doesn't exist, create it
	created := false
	if _, ok := s.rooms[roomID]; !ok {
		created = true
		log.Println("New room:", roomID)
		r := NewRoom()
		s.rooms[roomID] = r
		go s.Heartbeat(roomID)
	}
	room := s.rooms[roomID]
	return room, created
}

func (room *Room) NewConnection(ws *websocket.Conn, first bool) string {
	// assign the new connection a new id
	id := uuid.New().String()

	// set the last user
	room.lastUser = id

	// store the new connection by id
	room.conns[id] = ws

	// save current user
	room.nicks[id] = ""

	// send initial state if it's not the first connection
	if !first {
		frame := room.State.GenerateFullFrame(Full)
		evt := FrameJSON(frame)
		SendEvent(ws, evt)
	}

	return id
}

// these operate outside the main websocket loop
// see frontend/main.go suffixOp for corresponding receiver
func (s *Server) HandleOp(ws *websocket.Conn, op, roomID string) {
	data := ""
	room, ok := s.rooms[roomID]
	if !ok {
		EncodeSend(ws, data)
		return
	}
	switch op {
	case "sgf":
		// if the room doesn't exist, send empty string
		data = room.State.ToSGF(false)
	case "sgfix":
		// basically do the same thing but include indexes
		data = room.State.ToSGF(true)
	case "debug":
		// send debug info
		evt := room.State.InitData("handshake")
		data, _ = evt.Value.(string)
	}
	EncodeSend(ws, data)
}

// Echo the data received on the WebSocket.
func (s *Server) Handler(ws *websocket.Conn) {
	// new connection

	// first find the url they want
	url := ws.Request().URL.String()

	// currently not using the prefix, but i may someday
	_, roomID, op := ParseURL(url)

	// check for op suffix
	if op != "" {
		s.HandleOp(ws, op, roomID)
		return
	}

	// get or create the room
	room, first := s.GetOrCreateRoom(roomID)

	// assign id to the new connection
	id := room.NewConnection(ws, first)
	log.Println(url, "Connecting:", id)
	s.SendMessagesToOne(ws, id)

	// defer removing the client
	defer delete(room.conns, id)

	// send list of currently connected users
	room.SendUserList()

	// send disconnection notification
	// golang deferrals are called in LIFO order
	defer room.SendUserList()
	defer delete(room.nicks, id)

	handlers := map[string]EventHandler{
		"isprotected":   room.HandleIsProtected,
		"checkpassword": room.HandleCheckPassword,
		"debug":         HandleDebug,
		"ping":          HandlePing,

		"upload_sgf": Chain(
			room.HandleUploadSGF,
			room.OutsideBuffer,
			room.Authorized,
			room.CloseOGS,
			room.BroadcastAfter(false)),
		"request_sgf": Chain(
			room.HandleRequestSGF,
			room.OutsideBuffer,
			room.Authorized,
			room.CloseOGS,
			room.BroadcastAfter(false)),
		"trash": Chain(
			room.HandleTrash,
			room.OutsideBuffer,
			room.Authorized,
			room.CloseOGS,
			room.BroadcastAfter(false)),
		"update_nickname": Chain(
			room.HandleUpdateNickname,
			room.BroadcastAfter(false)),
		"update_settings": Chain(
			room.HandleUpdateSettings,
			room.Authorized,
			room.BroadcastConnectedUsersAfter,
			room.BroadcastAfter(false),
			room.BroadcastFullFrameAfter),
		"add_stone": Chain(
			room.HandleEvent,
			room.OutsideBuffer,
			room.Authorized,
			room.Slow,
			room.BroadcastAfter(true)),
		"_": Chain(
			room.HandleEvent,
			room.OutsideBuffer,
			room.Authorized,
			room.BroadcastAfter(true)),
	}

	// main loop
	for {
		// receive the event
		evt, err := ReceiveEvent(ws)
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
