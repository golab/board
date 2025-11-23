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
	"github.com/jarednogo/board/pkg/event"
	"github.com/jarednogo/board/pkg/fetch"
	"github.com/jarednogo/board/pkg/room/plugin"
	"github.com/jarednogo/board/pkg/state"
	"github.com/jarednogo/board/pkg/zip"
)

type EventHandler func(event.Event) event.Event

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
			r.logEventType,
			r.closeOGS,
			r.broadcastAfter,
			r.logAfter),
		"request_sgf": chain(
			r.handleRequestSGF,
			r.outsideBuffer,
			r.logEventType,
			r.logEventValue,
			r.authorized,
			r.closeOGS,
			r.broadcastAfter,
			r.logAfter),
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

func (r *Room) handleIsProtected(evt event.Event) event.Event {
	evt.SetValue(r.HasPassword())
	r.SendTo(evt.User(), evt)
	return evt
}

func (r *Room) handleCheckPassword(evt event.Event) event.Event {
	p := evt.Value().(string)

	password := r.GetPassword()

	if !core.CorrectPassword(p, password) {
		evt.SetValue("")
	} else {
		r.SetAuth(evt.User(), true)
	}
	r.SendTo(evt.User(), evt)
	return evt
}

func handleDebug(evt event.Event) event.Event {
	return evt
}

func handlePing(evt event.Event) event.Event {
	return evt
}

func (r *Room) handleUploadSGF(evt event.Event) event.Event {
	var bcast event.Event
	defer func() {
		if bcast != nil {
			bcast.SetUser(evt.User())
		}
	}()

	// it might be a string
	if str, ok := evt.Value().(string); ok {

		decoded, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			bcast = event.ErrorEvent(err.Error())
			return bcast
		}
		if zip.IsZipFile(decoded) {
			filesBytes, err := zip.Decompress(decoded)
			if err != nil {
				bcast = event.ErrorEvent(err.Error())
			} else {
				files := []string{}
				for _, file := range filesBytes {
					files = append(files, string(file))
				}
				merged := core.Merge(files)
				bcast = r.UploadSGF(merged)
			}
		} else {
			bcast = r.UploadSGF(string(decoded))
		}

	} else if arr, ok := evt.Value().([]any); ok {
		// it might be an array of strings
		sgfs := []string{}
		for _, ifc := range arr {
			str := ifc.(string)
			d, err := base64.StdEncoding.DecodeString(str)
			if err != nil {
				bcast = event.ErrorEvent(err.Error())
				return bcast
			}
			sgfs = append(sgfs, string(d))
		}
		sgf := core.Merge(sgfs)
		bcast = r.UploadSGF(sgf)
	} else {
		bcast = event.ErrorEvent("unreachable")
	}

	bcast.SetUser(evt.User())
	return bcast
}

func (r *Room) handleRequestSGF(evt event.Event) event.Event {
	var bcast event.Event
	defer func() {
		if bcast != nil {
			bcast.SetUser(evt.User())
		}
	}()

	url := evt.Value().(string)

	if fetch.IsOGS(url) {

		connectToOGS := false

		spl := strings.Split(url, "/")
		if len(spl) < 2 {
			bcast = event.ErrorEvent("url parsing error")
			return bcast
		}

		ogsType := spl[len(spl)-2]

		switch ogsType {
		case "game":
			ended, err := r.GetFetcher().OGSCheckEnded(url)
			if err != nil {
				bcast = event.ErrorEvent(err.Error())
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
				bcast = event.ErrorEvent("int parsing error")
				return bcast
			}
			id := int(id64)

			o, err := plugin.NewOGSConnector(r, r.GetFetcher())
			if err != nil {
				bcast = event.ErrorEvent("ogs connector error")
				return bcast
			}

			args := map[string]any{
				"key":     "ogs",
				"id":      id,
				"ogsType": ogsType,
			}
			r.RegisterPlugin(o, args)

			if ogsType == "game" {
				// finish here
				return event.NopEvent()
			}
		}
	}

	data, err := r.GetFetcher().ApprovedFetch(evt.Value().(string))
	if err != nil {
		bcast = event.ErrorEvent(err.Error())
	} else if data == "Permission denied" {
		bcast = event.ErrorEvent("Error fetching SGF. Is it a private OGS game?")
	} else {
		bcast = r.UploadSGF(string(data))
	}

	return bcast
}

func (r *Room) handleTrash(evt event.Event) event.Event {
	// reset room
	oldBuffer := r.GetInputBuffer()
	r.SetState(state.NewState(r.Size(), true))

	// reuse old inputbuffer
	r.SetInputBuffer(oldBuffer)

	frame := r.GenerateFullFrame(core.Full)
	bcast := event.FrameEvent(frame)
	bcast.SetUser(evt.User())
	return bcast
}

func (r *Room) handleUpdateNickname(evt event.Event) event.Event {
	nickname := evt.Value().(string)
	r.SetNick(evt.User(), nickname)
	userEvt := event.NewEvent("connected_users", r.Nicks())
	userEvt.SetUser(evt.User())
	return userEvt
}

type Settings struct {
	Buffer   int64
	Size     int
	Password string
}

func (r *Room) handleUpdateSettings(evt event.Event) event.Event {
	sMap := evt.Value().(map[string]any)
	buffer := int64(sMap["buffer"].(float64))
	size := int(sMap["size"].(float64))
	nickname := sMap["nickname"].(string)

	black := sMap["black"].(string)
	white := sMap["white"].(string)
	komi := sMap["komi"].(string)

	if black != "" {
		r.EditPlayerBlack(black)
	}
	if white != "" {
		r.EditPlayerWhite(white)
	}
	if komi != "" {
		r.EditKomi(komi)
	}

	r.SetNick(evt.User(), nickname)

	password := sMap["password"].(string)
	hashed := ""
	if password != "" {
		hashed = core.Hash(password)
	}
	settings := &Settings{buffer, size, hashed}

	r.SetInputBuffer(settings.Buffer)
	if settings.Size != r.Size() {
		// essentially trashing
		r.SetState(state.NewState(settings.Size, true))
		r.SetInputBuffer(buffer)
	}

	// can be changed
	// anyone already in the room is added
	// person who set password automatically gets added
	r.SetAuthAll()
	r.SetPassword(hashed)

	return evt
}

func (r *Room) handleEvent(evt event.Event) event.Event {
	var bcast event.Event
	defer func() {
		if bcast != nil {
			bcast.SetUser(evt.User())
		}
	}()

	cmd, err := state.DecodeToCommand(evt)
	if err != nil {
		bcast = event.ErrorEvent(err.Error())
		return bcast
	}

	frame, err := r.Execute(cmd)
	if err != nil {
		bcast = event.ErrorEvent(err.Error())
		return bcast
	}

	if frame != nil {
		bcast = event.FrameEvent(frame)
	} else {
		bcast = evt
	}
	bcast.SetUser(evt.User())
	return bcast
}

// middleware

func (r *Room) setTimeAfter(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		evt = handler(evt)
		// set last user information
		r.SetLastUser(evt.User())
		now := time.Now()
		r.SetLastActive(&now)
		return evt
	}
}

func (r *Room) broadcastAfter(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		evt = handler(evt)
		r.Broadcast(evt)
		return evt
	}
}

func (r *Room) broadcastFullFrameAfter(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		evt = handler(evt)
		frame := r.GenerateFullFrame(core.Full)
		bcast := event.FrameEvent(frame)
		r.Broadcast(bcast)
		return evt
	}
}

func (r *Room) broadcastConnectedUsersAfter(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		evt = handler(evt)
		userEvt := event.NewEvent("connected_users", r.Nicks())

		// broadcast connected_users
		r.Broadcast(userEvt)
		return evt
	}
}

func (r *Room) closeOGS(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		r.DeregisterPlugin("ogs")
		return handler(evt)
	}
}

func (r *Room) authorized(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		id := evt.User()
		ok := r.GetAuth(id)
		if r.GetPassword() == "" || ok {
			// only go to the next handler if authorized
			evt = handler(evt)
		}
		return evt
	}
}

// this one is to keep the same user from submitting multiple events too quickly
func (r *Room) slow(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		id := evt.User()
		// check multiple events from the same user in a narrow window (50 ms)
		now := time.Now()
		if last, ok := r.GetLastMessages(id); !ok {
			r.SetLastMessages(id, &now)
		} else {
			diff := now.Sub(last)
			r.SetLastMessages(id, &now)
			if diff.Milliseconds() < r.GetUserBuffer() {
				// don't do the next handler if too fast
				return evt
			}
		}
		return handler(evt)
	}
}

// this one is to keep people from tripping over each other
func (r *Room) outsideBuffer(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		lastUser := r.GetLastUser()

		if lastUser != evt.User() {
			now := time.Now()
			diff := now.Sub(r.GetLastActive())
			if diff.Milliseconds() < r.GetInputBuffer() {
				// don't do the next handler if too fast
				return evt
			}
		}
		return handler(evt)
	}
}

func (r *Room) logEventType(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		if r.logger != nil {
			r.logger.Info("handling event with type", "event_type", evt.Type())
		}
		return handler(evt)
	}
}

func (r *Room) logEventValue(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		if r.logger != nil {
			r.logger.Info("handling event with value", "event_value", evt.Value().(string))
		}
		return handler(evt)
	}
}

func (r *Room) logAfter(handler EventHandler) EventHandler {
	return func(evt event.Event) event.Event {
		evt = handler(evt)
		// only log info if we're not connected to ogs
		// (the ogs link doesn't upload the root node data in time)
		if r.logger != nil && r.GetPlugin("ogs") == nil {
			if evt.Type() == "error" {
				r.logger.Error("error while handling", "error", evt.Value().(string))
			} else {
				current := r.Current()
				args := []any{}
				for _, field := range current.AllFields() {
					args = append(args, field.Key)
					args = append(args, strings.Join(field.Values, ", "))
				}
				r.logger.Info("current node", args...)
			}
		}
		return evt
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
func (r *Room) HandleAny(evt event.Event) event.Event {
	// handle the event
	return r.GetHandler(evt.Type())(evt)
}
