/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room

import (
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/fetch"
	"github.com/jarednogo/board/pkg/room/plugin"
	"github.com/jarednogo/board/pkg/state"
	"github.com/jarednogo/board/pkg/zip"
)

type EventType string

const (
	EventTypeIsProtected    EventType = "isprotected"
	EventTypeCheckPassword  EventType = "checkpassword"
	EventTypeDebug          EventType = "debug"
	EventTypePing           EventType = "ping"
	EventTypeUploadSGF      EventType = "upload_sgf"
	EventTypeRequestSGF     EventType = "request_sgf"
	EventTypeTrash          EventType = "trash"
	EventTypeUpdateNickname EventType = "update_nickname"
	EventTypeUpdateSettings EventType = "update_settings"
	EventTypeAddStone       EventType = "add_stone"
	EventTypeGraft          EventType = "graft"
	EventTypeDefault        EventType = "_"
)

type EventHandler func(*core.EventJSON) *core.EventJSON

type Middleware func(EventHandler) EventHandler

// for use by server
func (r *Room) initHandlers() {
	r.handlers = map[string]EventHandler{
		"isprotected":   r.handleIsProtected,
		"checkpassword": r.handleCheckPassword,
		"debug":         handleDebug,
		"ping":          handlePing,

		"upload_sgf": chain(
			r.handleUploadSGF,
			r.outsideBuffer,
			r.authorized,
			r.closeOGS,
			r.broadcastAfter),
		"request_sgf": chain(
			r.handleRequestSGF,
			r.outsideBuffer,
			r.authorized,
			r.closeOGS,
			r.broadcastAfter),
		"trash": chain(
			r.handleTrash,
			r.outsideBuffer,
			r.authorized,
			r.closeOGS,
			r.broadcastAfter),
		"update_nickname": chain(
			r.handleUpdateNickname,
			r.broadcastAfter),
		"update_settings": chain(
			r.handleUpdateSettings,
			r.authorized,
			r.broadcastConnectedUsersAfter,
			r.broadcastAfter,
			r.broadcastFullFrameAfter),
		"add_stone": chain(
			r.handleEvent,
			r.outsideBuffer,
			r.authorized,
			r.slow,
			r.broadcastAfter,
			r.setTimeAfter),
		"graft": chain(
			r.handleEvent,
			r.broadcastFullFrameAfter),
		"_": chain(
			r.handleEvent,
			r.outsideBuffer,
			r.authorized,
			r.broadcastAfter,
			r.setTimeAfter),
	}
}

// handlers

func (room *Room) handleIsProtected(evt *core.EventJSON) *core.EventJSON {
	evt.Value = room.HasPassword()
	room.SendTo(evt.UserID, evt)
	return evt
}

func (room *Room) handleCheckPassword(evt *core.EventJSON) *core.EventJSON {
	p := evt.Value.(string)

	if !core.CorrectPassword(p, room.password) {
		evt.Value = ""
	} else {
		room.auth[evt.UserID] = true
	}
	room.SendTo(evt.UserID, evt)
	return evt
}

func handleDebug(evt *core.EventJSON) *core.EventJSON {
	return evt
}

func handlePing(evt *core.EventJSON) *core.EventJSON {
	return evt
}

func (room *Room) handleUploadSGF(evt *core.EventJSON) *core.EventJSON {
	var bcast *core.EventJSON
	defer func() {
		if bcast != nil {
			bcast.UserID = evt.UserID
		}
	}()

	// it might be a string
	if str, ok := evt.Value.(string); ok {

		decoded, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			bcast = core.ErrorEvent(err.Error())
			return bcast
		}
		if zip.IsZipFile(decoded) {
			filesBytes, err := zip.Decompress(decoded)
			if err != nil {
				bcast = core.ErrorEvent(err.Error())
			} else {
				files := []string{}
				for _, file := range filesBytes {
					files = append(files, string(file))
				}
				merged := core.Merge(files)
				bcast = room.UploadSGF(merged)
			}
		} else {
			bcast = room.UploadSGF(string(decoded))
		}

	} else if arr, ok := evt.Value.([]interface{}); ok {
		// it might be an array of strings
		sgfs := []string{}
		for _, ifc := range arr {
			str := ifc.(string)
			d, err := base64.StdEncoding.DecodeString(str)
			if err != nil {
				bcast = core.ErrorEvent(err.Error())
				return bcast
			}
			sgfs = append(sgfs, string(d))
		}
		sgf := core.Merge(sgfs)
		bcast = room.UploadSGF(sgf)
	} else {
		bcast = core.ErrorEvent("unreachable")
	}

	bcast.UserID = evt.UserID
	return bcast
}

func (room *Room) handleRequestSGF(evt *core.EventJSON) *core.EventJSON {
	var bcast *core.EventJSON
	defer func() {
		if bcast != nil {
			bcast.UserID = evt.UserID
		}
	}()

	url := evt.Value.(string)

	if fetch.IsOGS(url) {

		connectToOGS := false

		spl := strings.Split(url, "/")
		if len(spl) < 2 {
			bcast = core.ErrorEvent("url parsing error")
			return bcast
		}

		ogsType := spl[len(spl)-2]

		switch ogsType {
		case "game":
			ended, err := room.fetcher.OGSCheckEnded(url)
			if err != nil {
				bcast = core.ErrorEvent(err.Error())
				return bcast
			}
			connectToOGS = !ended
		case "review", "demo":
			ogsType = "review"
			connectToOGS = true
		}

		if connectToOGS {

			idStr := spl[len(spl)-1]
			id64, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				bcast = core.ErrorEvent("int parsing error")
				return bcast
			}
			id := int(id64)

			o, err := plugin.NewOGSConnector(room, room.fetcher)
			if err != nil {
				bcast = core.ErrorEvent("ogs connector error")
				return bcast
			}

			args := map[string]interface{}{
				"key":     "ogs",
				"id":      id,
				"ogsType": ogsType,
			}
			room.RegisterPlugin(o, args)

			if ogsType == "game" {
				// finish here
				return core.NopEvent()
			}
		}
	}

	data, err := room.fetcher.ApprovedFetch(evt.Value.(string))
	if err != nil {
		bcast = core.ErrorEvent(err.Error())
	} else if data == "Permission denied" {
		bcast = core.ErrorEvent("Error fetching SGF. Is it a private OGS game?")
	} else {
		bcast = room.UploadSGF(string(data))
	}

	return bcast
}

func (room *Room) handleTrash(evt *core.EventJSON) *core.EventJSON {

	// reset room
	oldBuffer := room.GetInputBuffer()
	room.engine = state.NewState(room.Size(), true)

	// reuse old inputbuffer
	room.SetInputBuffer(oldBuffer)

	frame := room.GenerateFullFrame(core.Full)
	bcast := core.FrameEvent(frame)
	bcast.UserID = evt.UserID
	return bcast
}

func (room *Room) handleUpdateNickname(evt *core.EventJSON) *core.EventJSON {
	nickname := evt.Value.(string)
	room.nicks[evt.UserID] = nickname
	userEvt := &core.EventJSON{
		Event:  "connected_users",
		Value:  room.nicks,
		UserID: evt.UserID,
	}
	return userEvt
}

type Settings struct {
	Buffer   int64
	Size     int
	Password string
}

func (room *Room) handleUpdateSettings(evt *core.EventJSON) *core.EventJSON {
	sMap := evt.Value.(map[string]interface{})
	buffer := int64(sMap["buffer"].(float64))
	size := int(sMap["size"].(float64))
	nickname := sMap["nickname"].(string)

	room.nicks[evt.UserID] = nickname

	password := sMap["password"].(string)
	hashed := ""
	if password != "" {
		hashed = core.Hash(password)
	}
	settings := &Settings{buffer, size, hashed}

	room.SetInputBuffer(settings.Buffer)
	if settings.Size != room.Size() {
		// essentially trashing
		room.engine = state.NewState(settings.Size, true)
		room.SetInputBuffer(buffer)
	}

	// can be changed
	// anyone already in the room is added
	// person who set password automatically gets added
	for connID := range room.conns {
		room.auth[connID] = true
	}
	room.password = hashed

	return evt
}

func (room *Room) handleEvent(evt *core.EventJSON) *core.EventJSON {
	var bcast *core.EventJSON
	defer func() {
		if bcast != nil {
			bcast.UserID = evt.UserID
		}
	}()

	frame, err := room.AddEvent(evt)
	if err != nil {
		bcast = core.ErrorEvent(err.Error())
		return bcast
	}

	if frame != nil {
		bcast = core.FrameEvent(frame)
	} else {
		bcast = evt
	}
	bcast.UserID = evt.UserID
	return bcast
}

// middleware

func (room *Room) setTimeAfter(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		evt = handler(evt)
		// set last user information
		room.mu.Lock()
		defer room.mu.Unlock()
		room.lastUser = evt.UserID
		now := time.Now()
		room.lastActive = &now
		return evt
	}
}

func (room *Room) broadcastAfter(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		evt = handler(evt)
		room.Broadcast(evt)
		return evt
	}
}

func (room *Room) broadcastFullFrameAfter(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		evt = handler(evt)
		frame := room.GenerateFullFrame(core.Full)
		bcast := core.FrameEvent(frame)
		room.Broadcast(bcast)
		return evt
	}
}

func (room *Room) broadcastConnectedUsersAfter(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		evt = handler(evt)
		userEvt := &core.EventJSON{
			Event:  "connected_users",
			Value:  room.nicks,
			UserID: "",
		}

		// broadcast connected_users
		room.Broadcast(userEvt)
		return evt
	}
}

func (room *Room) closeOGS(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		room.DeregisterPlugin("ogs")
		return handler(evt)
	}
}

func (room *Room) authorized(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		id := evt.UserID
		_, ok := room.auth[id]
		if room.password == "" || ok {
			// only go to the next handler if authorized
			evt = handler(evt)
		}
		return evt
	}
}

// this one is to keep the same user from submitting multiple events too quickly
func (room *Room) slow(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		id := evt.UserID
		// check multiple events from the same user in a narrow window (50 ms)
		now := time.Now()
		if last, ok := room.lastMessages[id]; !ok {
			room.lastMessages[id] = &now
		} else {
			diff := now.Sub(*last)
			room.lastMessages[id] = &now
			if diff.Milliseconds() < 50 {
				// don't do the next handler if too fast
				return evt
			}
		}
		return handler(evt)
	}
}

// this one is to keep people from tripping over each other
func (room *Room) outsideBuffer(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		if room.lastUser != evt.UserID {
			now := time.Now()
			diff := now.Sub(*room.lastActive)
			if diff.Milliseconds() < room.GetInputBuffer() {
				// don't do the next handler if too fast
				return evt
			}
		}
		return handler(evt)
	}
}

func chain(h EventHandler, middleware ...Middleware) EventHandler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// HandleAny is only to be used in special occasions because it recreates
// all the handlers
func (room *Room) HandleAny(evt *core.EventJSON) *core.EventJSON {
	// handle the event
	if handler, ok := room.handlers[evt.Event]; ok {
		return handler(evt)
	} else {
		return room.handlers["_"](evt)
	}
}
