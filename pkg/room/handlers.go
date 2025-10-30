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

type EventHandler func(*core.EventJSON) *core.EventJSON

type Middleware func(EventHandler) EventHandler

// for use by server
func (r *Room) CreateHandlers() map[string]EventHandler {
	return map[string]EventHandler{
		"isprotected":   r.HandleIsProtected,
		"checkpassword": r.HandleCheckPassword,
		"debug":         HandleDebug,
		"ping":          HandlePing,

		"upload_sgf": Chain(
			r.HandleUploadSGF,
			r.OutsideBuffer,
			r.Authorized,
			r.CloseOGS,
			r.BroadcastAfter),
		"request_sgf": Chain(
			r.HandleRequestSGF,
			r.OutsideBuffer,
			r.Authorized,
			r.CloseOGS,
			r.BroadcastAfter),
		"trash": Chain(
			r.HandleTrash,
			r.OutsideBuffer,
			r.Authorized,
			r.CloseOGS,
			r.BroadcastAfter),
		"update_nickname": Chain(
			r.HandleUpdateNickname,
			r.BroadcastAfter),
		"update_settings": Chain(
			r.HandleUpdateSettings,
			r.Authorized,
			r.BroadcastConnectedUsersAfter,
			r.BroadcastAfter,
			r.BroadcastFullFrameAfter),
		"add_stone": Chain(
			r.HandleEvent,
			r.OutsideBuffer,
			r.Authorized,
			r.Slow,
			r.BroadcastAfter,
			r.SetTimeAfter),
		"_": Chain(
			r.HandleEvent,
			r.OutsideBuffer,
			r.Authorized,
			r.BroadcastAfter,
			r.SetTimeAfter),
	}
}

// handlers

func (room *Room) HandleIsProtected(evt *core.EventJSON) *core.EventJSON {
	evt.Value = room.HasPassword()
	room.SendTo(evt.UserID, evt)
	return evt
}

func (room *Room) HandleCheckPassword(evt *core.EventJSON) *core.EventJSON {
	p := evt.Value.(string)

	if !core.CorrectPassword(p, room.password) {
		evt.Value = ""
	} else {
		room.auth[evt.UserID] = true
	}
	room.SendTo(evt.UserID, evt)
	return evt
}

func HandleDebug(evt *core.EventJSON) *core.EventJSON {
	return evt
}

func HandlePing(evt *core.EventJSON) *core.EventJSON {
	return evt
}

func (room *Room) HandleUploadSGF(evt *core.EventJSON) *core.EventJSON {
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

func (room *Room) HandleRequestSGF(evt *core.EventJSON) *core.EventJSON {
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

func (room *Room) HandleTrash(evt *core.EventJSON) *core.EventJSON {

	// reset room
	oldBuffer := room.state.InputBuffer
	room.state = state.NewState(room.state.Size, true)

	// reuse old inputbuffer
	room.state.InputBuffer = oldBuffer

	frame := room.state.GenerateFullFrame(core.Full)
	bcast := core.FrameEvent(frame)
	bcast.UserID = evt.UserID
	return bcast
}

func (room *Room) HandleUpdateNickname(evt *core.EventJSON) *core.EventJSON {
	nickname := evt.Value.(string)
	room.nicks[evt.UserID] = nickname
	userEvt := &core.EventJSON{
		Event:  "connected_users",
		Value:  room.nicks,
		Color:  0,
		UserID: evt.UserID,
	}
	return userEvt
}

type Settings struct {
	Buffer   int64
	Size     int
	Password string
}

func (room *Room) HandleUpdateSettings(evt *core.EventJSON) *core.EventJSON {
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

	room.state.InputBuffer = settings.Buffer
	if settings.Size != room.state.Size {
		// essentially trashing
		room.state = state.NewState(settings.Size, true)
		room.state.InputBuffer = buffer
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

func (room *Room) HandleEvent(evt *core.EventJSON) *core.EventJSON {
	var bcast *core.EventJSON
	defer func() {
		if bcast != nil {
			bcast.UserID = evt.UserID
		}
	}()

	frame, err := room.state.AddEvent(evt)
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

func (room *Room) SetTimeAfter(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		evt = handler(evt)
		// set last user information
		room.lastUser = evt.UserID
		now := time.Now()
		room.lastActive = &now
		return evt
	}
}

func (room *Room) BroadcastAfter(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		evt = handler(evt)
		room.Broadcast(evt)
		return evt
	}
}

func (room *Room) BroadcastFullFrameAfter(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		evt = handler(evt)
		frame := room.state.GenerateFullFrame(core.Full)
		bcast := core.FrameEvent(frame)
		room.Broadcast(bcast)
		return evt
	}
}

func (room *Room) BroadcastConnectedUsersAfter(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		evt = handler(evt)
		userEvt := &core.EventJSON{
			Event:  "connected_users",
			Value:  room.nicks,
			Color:  0,
			UserID: "",
		}

		// broadcast connected_users
		room.Broadcast(userEvt)
		return evt
	}
}

func (room *Room) CloseOGS(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		room.DeregisterPlugin("ogs")
		return handler(evt)
	}
}

func (room *Room) Authorized(handler EventHandler) EventHandler {
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
func (room *Room) Slow(handler EventHandler) EventHandler {
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
func (room *Room) OutsideBuffer(handler EventHandler) EventHandler {
	return func(evt *core.EventJSON) *core.EventJSON {
		if room.lastUser != evt.UserID {
			now := time.Now()
			diff := now.Sub(*room.lastActive)
			if diff.Milliseconds() < room.state.InputBuffer {
				// don't do the next handler if too fast
				return evt
			}
		}
		return handler(evt)
	}
}

func Chain(h EventHandler, middleware ...Middleware) EventHandler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}
