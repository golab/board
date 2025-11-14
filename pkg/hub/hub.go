/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/jarednogo/board/pkg/config"
	"github.com/jarednogo/board/pkg/event"
	"github.com/jarednogo/board/pkg/fetch"
	"github.com/jarednogo/board/pkg/loader"
	"github.com/jarednogo/board/pkg/logx"
	"github.com/jarednogo/board/pkg/message"
	"github.com/jarednogo/board/pkg/room"
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
	messages []*message.Message
	db       loader.Loader
	mu       sync.Mutex
	cfg      *config.Config
	logger   logx.Logger
}

func NewHub(cfg *config.Config) (*Hub, error) {
	// get database setup
	var db loader.Loader
	if cfg.DB.Type == config.DBConfigTypeMemory {
		db = loader.NewMemoryLoader()
	} else {
		db = loader.NewDefaultLoader(cfg.DB.Path)
	}
	return NewHubWithDB(db, cfg)
}

func NewHubWithDB(db loader.Loader, cfg *config.Config) (*Hub, error) {
	err := db.Setup()
	if err != nil {
		return nil, err
	}
	s := &Hub{
		rooms:    make(map[string]*room.Room),
		messages: []*message.Message{},
		db:       db,
		cfg:      cfg,
		logger:   logx.NewDefaultLogger(),
	}

	// start message loop
	go s.MessageLoop()

	return s, nil
}

func (h *Hub) GetRoom(id string) (*room.Room, error) {
	if room, ok := h.rooms[id]; ok {
		return room, nil
	}
	return nil, errors.New("room not found")
}

func (h *Hub) RoomCount() int {
	return len(h.rooms)
}

func (h *Hub) MessageCount() int {
	return len(h.messages)
}

func (h *Hub) Save() {
	for id, r := range h.rooms {
		h.logger.Info("saving", "room_id", id)

		save := r.Save()

		err := h.db.SaveRoom(id, save)
		if err != nil {
			h.logger.Error("failed to save room", "err", err, "room_id", id)
		}
	}
}

func (h *Hub) Load() {
	rooms, err := h.db.LoadAllRooms()
	if err != nil {
		h.logger.Error("failed to load", "err", err)
		return
	}

	for _, load := range rooms {
		r, err := room.Load(load)
		if err != nil {
			h.logger.Error("failed to load room", "err", err)
			continue
		}

		if h.cfg.Mode == config.ModeTest {
			r.SetFetcher(fetch.NewEmptyFetcher())
		}

		id := r.ID()
		h.logger.Info("loading", "room_id", id)
		h.mu.Lock()
		h.rooms[id] = r
		h.mu.Unlock()
		go h.Heartbeat(id)
	}
}

func (h *Hub) Heartbeat(roomID string) {
	h.mu.Lock()
	r, ok := h.rooms[roomID]
	h.mu.Unlock()
	if !ok {
		return
	}
	for {
		now := time.Now()
		diff := now.Sub(*r.LastActive())
		h.logger.Info("inactive", "room_id", roomID, "duration", diff.Seconds())
		if diff.Seconds() > r.GetTimeout() {
			break
		}
		time.Sleep(3600 * time.Second)
	}
	h.logger.Info("clearing board", "room_id", roomID)

	// close the room down
	err := r.Close()
	if err != nil {
		h.logger.Error("failed to close room", "err", err)
	}

	// delete the room from the server map
	delete(h.rooms, roomID)

	// delete it from the database
	err = h.db.DeleteRoom(roomID)
	if err != nil {
		h.logger.Error("failed to delete room", "err", err)
	}
}

func (h *Hub) ReadMessages() {
	messages, err := h.db.LoadAllMessages()
	if err != nil {
		h.logger.Error("failed to load messages", "err", err)
		return
	}
	defer h.db.DeleteAllMessages() //nolint:errcheck

	for _, msg := range messages {
		m := message.New(msg.Text, msg.TTL)
		h.messages = append(h.messages, m)
	}
}

func (h *Hub) SendMessages() {
	// go through each server message
	keep := []*message.Message{}
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
	h.mu.Lock()
	defer h.mu.Unlock()
	// if the room they want doesn't exist, create it
	if _, ok := h.rooms[roomID]; !ok {
		h.logger.Info("new room", "room_id", roomID)
		r := room.NewRoom(roomID)
		if h.cfg.Mode == config.ModeTest {
			r.SetFetcher(fetch.NewEmptyFetcher())
		}
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
	ec := event.NewDefaultEventChannel(ws)
	h.Handler(ec, roomID)
}

func (h *Hub) Handler(ec event.EventChannel, roomID string) {
	// new connection

	// get or create the room
	r := h.GetOrCreateRoom(roomID)

	// send to the room for handling
	h.logger.Info(
		"new connection",
		"event_channel", ec.ID(),
		"room_id", r.ID(),
	)
	h.logger.Info(
		"disconnection",
		"event_channel", ec.ID(),
		"room_id", r.ID(),
		"reason", r.Handle(ec),
	)
}
