/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub

import (
	"log"
	"strings"
	"time"

	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/loader"
	"github.com/jarednogo/board/pkg/room"
	"github.com/jarednogo/board/pkg/socket"
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

func NewHub() (*Hub, error) {
	// get database setup
	db := loader.NewDefaultLoader()
	return NewHubWithDB(db)
}

func NewHubWithDB(db loader.Loader) (*Hub, error) {
	err := db.Setup()
	if err != nil {
		return nil, err
	}
	s := &Hub{
		rooms:    make(map[string]*room.Room),
		messages: []*core.Message{},
		db:       db,
	}

	// start message loop
	go s.MessageLoop()

	return s, nil
}

func (h *Hub) RoomCount() int {
	return len(h.rooms)
}

func (h *Hub) MessageCount() int {
	return len(h.messages)
}

func (h *Hub) Save() {
	for id, r := range h.rooms {
		log.Printf("Saving %s", id)

		save := r.Save()

		err := h.db.SaveRoom(id, save)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *Hub) Load() {
	rooms, err := h.db.LoadAllRooms()
	if err != nil {
		log.Println(err)
		return
	}

	for _, load := range rooms {
		r, err := room.Load(load)
		if err != nil {
			log.Println(err)
			continue
		}

		id := r.ID()
		log.Printf("Loading %s", id)
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
		diff := now.Sub(*r.LastActive())
		log.Println(roomID, "Inactive for", diff)
		if diff.Seconds() > r.Timeout() {
			break
		}
		time.Sleep(3600 * time.Second)
	}
	log.Println("Cleaning up board due to inactivity:", roomID)

	// close the room down
	err := r.Close()
	if err != nil {
		log.Println("errors while closing:", err)
	}

	// delete the room from the server map
	delete(h.rooms, roomID)

	// delete it from the database
	err = h.db.DeleteRoom(roomID)
	if err != nil {
		log.Println(err)
	}
}

func (h *Hub) ReadMessages() {
	messages, err := h.db.LoadAllMessages()
	if err != nil {
		log.Println(err)
		return
	}
	defer h.db.DeleteAllMessages() //nolint:errcheck

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
	if _, ok := h.rooms[roomID]; !ok {
		log.Println("New room:", roomID)
		r := room.NewRoom(roomID)
		h.rooms[roomID] = r
		go h.Heartbeat(roomID)
	}
	r := h.rooms[roomID]
	return r
}

func (h *Hub) HandlerWrapper(ws *websocket.Conn) {
	// first find the url they want
	url := ws.Request().URL.String()

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
	log.Printf(
		"[*] New connection: %s to %s",
		rc.ID(),
		r.ID(),
	)
	log.Printf(
		"[-] Disconnection: %s from %s (%v)",
		rc.ID(),
		r.ID(),
		r.Handle(rc),
	)
}
